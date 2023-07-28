package main

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/h2non/filetype"
	"gorm.io/gorm"

	"tiktok-demo/cmd/video/config"
	"tiktok-demo/cmd/video/pkg/minio"
	"tiktok-demo/cmd/video/pkg/mysql"
	"tiktok-demo/cmd/video/pkg/pack"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/CommentServer"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
	VideoServer "tiktok-demo/shared/kitex_gen/VideoServer"
	"tiktok-demo/shared/tools"
)

// VideoSrvImpl implements the last service interface defined in the IDL.
type VideoSrvImpl struct {
	MysqlManager *mysql.Manager
	MinioManager *minio.Manager
	RealtionManager
	UserManager
	CommentManager
	FavoriteManager
}

type UserManager interface {
	GetUserInfo(ctx context.Context, req *UserServer.DouyinUserRequest, callOptions ...callopt.Option) (resp *UserServer.DouyinUserResponse, err error)
}

type CommentManager interface {
	CommentList(ctx context.Context, req *CommentServer.DouyinCommentListRequest, callOptions ...callopt.Option) (resp *CommentServer.DouyinCommentListResponse, err error)
}

type FavoriteManager interface {
	GetFavoriteVideo(ctx context.Context, req *FavoriteServer.DouyinVideoBeFavoriteRequest, callOptions ...callopt.Option) (resp *FavoriteServer.DouyinVideoBeFavoriteResponse, err error)
	QueryUserLikeVideo(ctx context.Context, req *FavoriteServer.DouyinQueryFavoriteRequest, callOptions ...callopt.Option) (resp *FavoriteServer.DouyinQueryFavoriteResponse, err error)
}

type RealtionManager interface {
	QueryRelation(ctx context.Context, req *RelationServer.DouyinQueryRelationRequest, callOptions ...callopt.Option) (resp *RelationServer.DouyinQueryRelationResponse, err error)
}

// Feed implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) Feed(ctx context.Context, req *VideoServer.DouyinFeedRequest) (resp *VideoServer.DouyinFeedResponse, err error) {
	if req == nil {
		return pack.BuildFeedResp(nil, nil, 0), nil
	}
	if req.LatestTime == 0 {
		req.LatestTime = int64(time.Now().UnixMilli())
	}
	// 查询数据库
	videoFeed, err := s.MysqlManager.GetVideoByTime(req.LatestTime, consts.VideoNumPerFeed)
	if err != nil {
		klog.Error("mysql get video error", err)
		return pack.BuildFeedResp(err, nil, 0), nil
	}
	nextTime := time.Now().UnixMilli()
	if len(videoFeed) == 0 {
		return pack.BuildFeedResp(nil, nil, nextTime), nil
	} else {
		nextTime = videoFeed[len(videoFeed)-1].CreatedAt.UnixMilli()
	}
	videos := pack.Video2VideoServerList(videoFeed, req.UserId)
	return pack.BuildFeedResp(nil, videos, nextTime), nil

}

// PublishAction implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) PublishAction(ctx context.Context, req *VideoServer.DouyinPublishActionRequest) (resp *VideoServer.DouyinPublishActionResponse, err error) {
	// 处理视频数据，检查视频是否合法
	videoData := []byte(req.Data)
	kind, err := filetype.Match(videoData)
	if err != nil {
		klog.Error("user %d match file type error", req.UserId, err)
		return pack.BuildPublishActionResp(errno.PublishActionErr), nil
	}
	if kind.MIME.Type != "video" {
		klog.Error("user %d upload video type error", req.UserId)
		return pack.BuildPublishActionResp(errno.VideoTypeErr), nil
	}
	contentType := kind.Extension

	// 生成视频id
	sf, err := snowflake.NewNode(consts.VideoSnowflakeNode)
	if err != nil {
		klog.Fatalf("user %d generate snowflake error", req.UserId, err)
		return pack.BuildPublishActionResp(errno.PublishActionErr), nil
	}
	videoId := uint(sf.Generate().Int64())
	videoName := strconv.Itoa(int(videoId)) + "." + contentType
	// 获取视频URL
	videoUrl := "http://" + config.GlobalServerConfig.MinioInfo.MinioURL + ":" + config.GlobalServerConfig.MinioInfo.MinioPort + "/" + s.MinioManager.VideosBucket + "/" + videoName

	// 封面图片
	coverName := strconv.Itoa(int(videoId)) + ".jpg"
	localVideoPath := "http://127.0.0.1:9000/" + s.MinioManager.VideosBucket + "/" + videoName
	coverUrl := "http://" + config.GlobalServerConfig.MinioInfo.MinioURL + ":" + config.GlobalServerConfig.MinioInfo.MinioPort + "/" + s.MinioManager.CoverBucket + "/" + coverName

	// channel用于接收上传结果
	uploadSuccess := make(chan bool)

	// 异步上传视频和封面
	go func() {
		videoReader := bytes.NewReader(videoData)
		err = s.MinioManager.UploadObject(ctx, "video", s.MinioManager.VideosBucket, videoName, videoReader, int64(len(videoData)))
		if err != nil {
			klog.Error("user %d upload video %s error", req.UserId, videoName, err)
			uploadSuccess <- false
			return
		}
		klog.Infof("user %d upload video %s success", req.UserId, videoName)

		// 封面图片
		coverData, err := tools.GetVideoCover(localVideoPath)
		if err != nil {
			klog.Error("user %d get video %s cover error", req.UserId, videoName, err)
			uploadSuccess <- false
			return
		}
		coverReader := bytes.NewReader(coverData)
		err = s.MinioManager.UploadObject(ctx, "cover", s.MinioManager.CoverBucket, coverName, coverReader, int64(len(coverData)))
		if err != nil {
			klog.Error("user %d upload video %s cover error", req.UserId, videoName, err)
			uploadSuccess <- false
			return
		}
		klog.Infof("user %d upload video %s cover success", req.UserId, videoName)
		uploadSuccess <- true
	}()

	// 视频信息入库
	go func() {
		if success := <-uploadSuccess; success {
			videoModel := &mysql.Video{
				Model: gorm.Model{
					ID:        uint(videoId),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				AuthorId:      req.UserId,
				VideoURL:      videoUrl,
				CoverURL:      coverUrl,
				FavoriteCount: 0,
				CommentCount:  0,
				Title:         req.Title,
			}

			if err := s.MysqlManager.PublishVideo(videoModel); err != nil {
				s.MinioManager.RemoveObject(ctx, s.MinioManager.VideosBucket, videoName)
				s.MinioManager.RemoveObject(ctx, s.MinioManager.CoverBucket, coverName)
				klog.Error("publish video %d error", videoModel.ID, err)
				return
			}
		} else {
			s.MinioManager.RemoveObject(ctx, s.MinioManager.VideosBucket, videoName)
			s.MinioManager.RemoveObject(ctx, s.MinioManager.CoverBucket, coverName)
		}
	}()

	return pack.BuildPublishActionResp(nil), nil
}

// PublishList implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) PublishList(ctx context.Context, req *VideoServer.DouyinPublishListRequest) (resp *VideoServer.DouyinPublishListResponse, err error) {
	pubList, err := s.MysqlManager.GetVideoByUserId(req.UserId)
	if err != nil {
		klog.Error("get user %d publish list error", req.UserId, err)
		if err == errno.VideoNotExistErr {
			return pack.BuildPublishListResp(errno.VideoNotExistErr, nil), nil
		}
		return pack.BuildPublishListResp(errno.PublishListErr, nil), nil
	}
	// todo 这里的userId是不是应该是作者的id
	videos := pack.Video2VideoServerList(pubList, req.UserId)
	return pack.BuildPublishListResp(nil, videos), nil
}

// GetVideoListByVideoId implements the VideoSrvImpl interface.
func (s *VideoSrvImpl) GetVideoListByVideoId(ctx context.Context, req *VideoServer.DouyinVideoListByVideoId) (resp *VideoServer.DouyinPublishListResponse, err error) {
	var VideoList []*mysql.Video
	for _, videoId := range req.VideoId {
		video, err := s.MysqlManager.GetVideoById(videoId)
		if err != nil {
			klog.Error("get video %d error", videoId, err)
			return pack.BuildPublishListResp(errno.VideoNotExistErr, nil), nil
		}
		VideoList = append(VideoList, video)
	}
	videos := pack.Video2VideoServerList(VideoList, req.UserId)
	return pack.BuildPublishListResp(nil, videos), nil
}
