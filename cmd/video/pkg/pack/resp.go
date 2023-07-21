package pack

import (
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"sync"

	"tiktok-demo/cmd/video/config"
	"tiktok-demo/cmd/video/pkg/mysql"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/CommentServer"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
	"tiktok-demo/shared/kitex_gen/VideoServer"
)

func feedResp(err errno.ErrNo, videos []*VideoServer.Video, nextTime int64) *VideoServer.DouyinFeedResponse {
	resp := new(VideoServer.DouyinFeedResponse)
	resp.BaseResp = &VideoServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.VideoList = videos
	resp.NextTime = nextTime
	return resp
}

func BuildFeedResp(err error, videos []*VideoServer.Video, nextTime int64) *VideoServer.DouyinFeedResponse {
	if err == nil {
		return feedResp(errno.Success, videos, nextTime)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return feedResp(e, nil, nextTime)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return feedResp(s, nil, nextTime)
}

func publishActionResp(err errno.ErrNo) *VideoServer.DouyinPublishActionResponse {
	resp := new(VideoServer.DouyinPublishActionResponse)
	resp.BaseResp = &VideoServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	return resp
}

func BuildPublishActionResp(err error) *VideoServer.DouyinPublishActionResponse {
	if err == nil {
		return publishActionResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return publishActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return publishActionResp(s)
}

func publishListResp(err errno.ErrNo, videos []*VideoServer.Video) *VideoServer.DouyinPublishListResponse {
	resp := new(VideoServer.DouyinPublishListResponse)
	resp.BaseResp = &VideoServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.VideoList = videos

	return resp
}

func BuildPublishListResp(err error, videos []*VideoServer.Video) *VideoServer.DouyinPublishListResponse {
	if err == nil {
		return publishListResp(errno.Success, videos)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return publishListResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return publishListResp(s, nil)
}

// db.Video -> videoServer.Video

func Video2VideoServerList(vs []*mysql.Video, userId int64) []*VideoServer.Video {
	var videos []*VideoServer.Video
	for _, v := range vs {
		videos = append(videos, Video2VideoServer(v, userId))
	}
	return videos
}

func Video2VideoServer(v *mysql.Video, userId int64) *VideoServer.Video {
	author, isFol, CommentCount, FavoriteCount, isFav := GetVideoInfo(context.Background(), v, userId)
	return &VideoServer.Video{
		Id: int64(v.ID),
		Author: &VideoServer.User{
			Id:            author.Id,
			Name:          author.Name,
			FollowCount:   author.FollowCount,
			FollowerCount: author.FollowerCount,
			IsFollow:      isFol,
		},
		PlayUrl:       v.VideoURL,
		CoverUrl:      v.CoverURL,
		FavoriteCount: FavoriteCount,
		CommentCount:  CommentCount,
		IsFavorite:    isFav,
		Title:         v.Title,
	}
}

func GetVideoInfo(ctx context.Context, v *mysql.Video, userId int64) (author *UserServer.User, isFol bool, comCount int64, favCount int64, isFav bool) {
	klog.Infof("GetVideoInfo start")
	// 开启 goroutine 并发调用RPC获取视频信息
	var wg sync.WaitGroup
	wg.Add(5)
	var err error
	// 获取作者信息
	go func() {
		var resp *UserServer.DouyinUserResponse
		resp, err = config.UserClient.GetUserInfo(ctx, &UserServer.DouyinUserRequest{
			UserId: v.AuthorId,
		})
		if resp == nil {
			klog.Infof("GetUserInfo failed with %s in %v", err.Error(), v.AuthorId)
		}
		if resp.BaseResp.StatusCode != 0 {
			klog.Infof("GetUserInfo failed with %s in %v", resp.BaseResp.StatusMsg, v.AuthorId)
		} else {
			author = resp.User
		}
		wg.Done()
	}()
	// 查询 userId 是否关注 authorId
	go func() {
		// 未登录
		if userId == 0 {
			isFol = false
		} else {
			var resp *RelationServer.DouyinQueryRelationResponse
			resp, err = config.RelationClient.QueryRelation(ctx, &RelationServer.DouyinQueryRelationRequest{
				UserId:   userId,
				ToUserId: v.AuthorId,
			})
			if resp == nil {
				klog.Infof("QueryRelation failed with %s in %v", err.Error(), userId)
			}
			if resp.BaseResp.StatusCode != 0 {
				klog.Infof("QueryRelation failed with %s in %v", err.Error(), userId)
			} else {
				isFol = resp.IsFollow
			}
		}
		wg.Done()
	}()
	// 查询视频评论数
	go func() {
		var resp *CommentServer.DouyinCommentListResponse
		resp, err = config.CommentClient.CommentList(ctx, &CommentServer.DouyinCommentListRequest{
			VideoId: int64(v.ID),
		})
		if resp == nil {
			klog.Infof("GetCommentList failed with %s in %v", err.Error(), v.ID)
		}
		if resp.BaseResp.StatusCode != 0 {
			klog.Infof("GetCommentList failed with %s in %v", resp.BaseResp.StatusMsg, v.ID)
		}
		comCount = int64(len(resp.CommentList))
		wg.Done()
	}()
	// 查询视频点赞数
	go func() {
		var resp *FavoriteServer.DouyinVideoBeFavoriteResponse
		resp, err = config.FavoriteClient.GetFavoriteVideo(ctx, &FavoriteServer.DouyinVideoBeFavoriteRequest{
			VideoId: int64(v.ID),
		})
		if resp == nil {
			klog.Infof("GetFavoriteVideo failed with %s in %v", err.Error(), v.ID)
		}
		if resp.BaseResp.StatusCode != 0 {
			klog.Infof("GetFavoriteVideo failed with %s in %v", resp.BaseResp.StatusMsg, v.ID)
		} else {
			favCount = resp.FavoriteCount
		}
		wg.Done()
	}()
	// 查询用户是否点赞视频
	go func() {
		// var resp *FavoriteServer.DouyinQueryFavoriteResponse
		resp, err := config.FavoriteClient.QueryUserLikeVideo(ctx, &FavoriteServer.DouyinQueryFavoriteRequest{
			UserId:  userId,
			VideoId: int64(v.ID),
		})
		if resp == nil {
			klog.Infof("QueryUserLikeVideo failed with %s in %v", err.Error(), userId)
		}
		if resp.BaseResp.StatusCode != 0 {
			klog.Infof("QueryUserLikeVideo failed with %s in %v", resp.BaseResp.StatusMsg, userId)
		} else {
			isFav = resp.Favorite
		}
		wg.Done()
	}()
	wg.Wait()
	return author, isFol, comCount, favCount, isFav
}
