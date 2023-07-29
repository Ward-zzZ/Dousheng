package main

import (
	"context"
	"strconv"
	"strings"

	"tiktok-demo/cmd/favorite/pkg/mq"
	"tiktok-demo/cmd/favorite/pkg/mysql"
	"tiktok-demo/cmd/favorite/pkg/pack"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/VideoServer"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct {
	MysqlManager
	RedisManager
	VideoManager
}

type MysqlManager interface {
	GetFavoriteUserIdList(videoId int64) ([]*mysql.Favorite, error)
	GetFavoriteVideoIdList(userId int64) ([]*mysql.Favorite, error)
	QueryFavorite(userId int64, videoId int64) (bool, error)
	FavoriteAction(userId int64, videoId int64, isFavorite bool) error
	InsertFavorite(userId int64, videoId int64, isFavorite bool) error
	UpdateFavorite(userId int64, videoId int64, isFavorite bool) error
	GetFavoriteCountByVideoId(videoId int64) (int64, error)
}

type RedisManager interface {
	SetUserLikeList(c context.Context, uid int64, vids []int64) error
	AddUserLikeList(c context.Context, uid int64, vid int64) (bool, error)
	DelUserLikeList(c context.Context, uid int64, vid int64) (bool, error)
	QueryUserLike(c context.Context, uid int64, vid int64) (bool, error)
	GetUserLikeList(c context.Context, uid int64) ([]int64, error)
	SetVideoLikeCount(c context.Context, vid int64, count int64) error
	AddVideoLikeCount(c context.Context, vid int64) error
	DelVideoLikeCount(c context.Context, vid int64) error
	GetVideoLikeCount(c context.Context, vid int64) (int64, error)
}

type VideoManager interface {
	GetVideoListByVideoId(ctx context.Context, Req *VideoServer.DouyinVideoListByVideoId, callOptions ...callopt.Option) (r *VideoServer.DouyinPublishListResponse, err error)
}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *FavoriteServer.DouyinFavoriteActionRequest) (resp *FavoriteServer.DouyinFavoriteActionResponse, err error) {
	resp = pack.BuildfavoriteActionResp(nil)

	//先检查当前的点赞状态
	isFavorite, err := s.RedisManager.QueryUserLike(ctx, req.UserId, req.VideoId)
	if err != nil {
		// redis不存在，从mysql中获取用户点赞列表
		favoriteList, _ := s.MysqlManager.GetFavoriteVideoIdList(req.UserId)
		if len(favoriteList) == 0 {
			isFavorite = false
		} else {
			favoriteListId := make([]int64, 0)
			for _, v := range favoriteList {
				favoriteListId = append(favoriteListId, v.VideoId)
				if v.VideoId == req.VideoId {
					isFavorite = true
				}
			}
			// 更新redis
			s.RedisManager.SetUserLikeList(ctx, req.UserId, favoriteListId)
		}
	}
	if isFavorite && (req.ActionType == 1) {
		resp = pack.BuildfavoriteActionResp(errno.FavoriteActionTypeErr.WithMessage("已经点赞"))
		return resp, nil
	}
	if !isFavorite && req.ActionType == 2 {
		resp = pack.BuildfavoriteActionResp(errno.FavoriteActionTypeErr.WithMessage("已经取消点赞"))
		return resp, nil
	}

	if req.ActionType != 1 && req.ActionType != 2 {
		resp = pack.BuildfavoriteActionResp(errno.FavoriteActionTypeErr.WithMessage("点赞类型错误"))
	}
	// 先更新redis，redis不存在时跳过更新
	if req.ActionType == 1 {
		if _, err = s.RedisManager.AddUserLikeList(ctx, req.UserId, req.VideoId); err != nil {
			klog.Errorf("redis add user like list err:%v", err)
		}
		if err = s.RedisManager.AddVideoLikeCount(ctx, req.VideoId); err != nil {
			klog.Errorf("redis add video like count err:%v", err)
		}
	} else {
		if _, err = s.RedisManager.DelUserLikeList(ctx, req.UserId, req.VideoId); err != nil {
			klog.Errorf("redis del user like list err:%v", err)
		}
		if err = s.RedisManager.DelVideoLikeCount(ctx, req.VideoId); err != nil {
			klog.Errorf("redis add video like count err:%v", err)
		}
	}

	// 异步更新mysql
	// todo: sync redis and mysql periodically
	msg := strings.Builder{}
	msg.WriteString(strconv.FormatInt(req.UserId, 10))
	msg.WriteString(":")
	msg.WriteString(strconv.FormatInt(req.VideoId, 10))
	msg.WriteString(":")
	msg.WriteString(strconv.Itoa(int(req.ActionType)))
	err = mq.FavoriteActor.Publish(ctx, msg.String())
	if err != nil {
		klog.Errorf("mq publish err:%v", err)
		resp = pack.BuildfavoriteActionResp(errno.FavoriteActionErr)
		return resp, nil
	}

	return resp, nil
}

// GetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteList(ctx context.Context, req *FavoriteServer.DouyinFavoriteListRequest) (resp *FavoriteServer.DouyinFavoriteListResponse, err error) {
	// 1. get user's favorite video id list
	videoIds := make([]int64, 0)
	// 1.1 check redis
	ids, err := s.RedisManager.GetUserLikeList(ctx, req.UserId)
	if err != nil || len(ids) == 0 {
		// 1.2 check mysql
		favList, err := s.MysqlManager.GetFavoriteVideoIdList(req.UserId)
		if err != nil {
			if err == errno.FavoriteVideoListNotExistErr {
				return pack.BuildgetFavoriteListResp(errno.Success, nil), nil
			}
			return pack.BuildgetFavoriteListResp(errno.QueryUserLikeVideoErr, nil), nil
		}
		// 1.3 update redis
		for _, v := range favList {
			videoIds = append(videoIds, v.VideoId)
		}
		s.RedisManager.SetUserLikeList(ctx, req.UserId, videoIds)
	} else {
		videoIds = ids
	}

	// 2. get video list
	videoList, err := s.VideoManager.GetVideoListByVideoId(ctx, &VideoServer.DouyinVideoListByVideoId{
		VideoId: videoIds,
		UserId:  req.UserId,
	})
	if videoList.BaseResp.StatusCode != errno.SuccessCode {
		klog.Errorf("get video list err:%v", videoList.BaseResp)
		return pack.BuildgetFavoriteListResp(errno.FavoriteVideoListErr, nil), nil
	}
	video := pack.ConvertVideos(videoList.VideoList)
	resp = pack.BuildgetFavoriteListResp(nil, video)

	return
}

// GetFavoriteUser implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteUser(ctx context.Context, req *FavoriteServer.DouyinUserBeFavoriteRequest) (resp *FavoriteServer.DouyinUserBeFavoriteResponse, err error) {
	return
}

// GetFavoriteVideo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteVideo(ctx context.Context, req *FavoriteServer.DouyinVideoBeFavoriteRequest) (resp *FavoriteServer.DouyinVideoBeFavoriteResponse, err error) {
	// 1. check redis
	count, err := s.RedisManager.GetVideoLikeCount(ctx, req.VideoId)
	if err != nil {
		// 2. check mysql
		count, err = s.MysqlManager.GetFavoriteCountByVideoId(req.VideoId)
		if err != nil {
			klog.Errorf("get favorite count err:%v", err)
			resp = pack.BuildfavoriteVideoQueryResp(errno.QueryFavoriteCountErr, 0)
			return resp, nil
		}
		// 3. update redis
		//异步更新redis
		go s.RedisManager.SetVideoLikeCount(ctx, req.VideoId, count)
	}
	resp = pack.BuildfavoriteVideoQueryResp(nil, count)
	return
}

// QueryUserLikeVideo implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) QueryUserLikeVideo(ctx context.Context, req *FavoriteServer.DouyinQueryFavoriteRequest) (resp *FavoriteServer.DouyinQueryFavoriteResponse, err error) {
	// 1. check redis
	isFavorite, err := s.RedisManager.QueryUserLike(ctx, req.UserId, req.VideoId)
	if err != nil {
		// 2. check mysql
		isFavorite, err = s.MysqlManager.QueryFavorite(req.UserId, req.VideoId)
		if err != nil {
			klog.Errorf("query user like video err:%v", err)
			resp = pack.BuildQueryUserFavoriteVideoResp(errno.QueryUserLikeVideoErr, false)
			return resp, nil
		}
	}
	resp = pack.BuildQueryUserFavoriteVideoResp(nil, isFavorite)
	return
}
