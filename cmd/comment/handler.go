package main

import (
	"context"
	"strconv"
	"time"

	"tiktok-demo/cmd/comment/pkg/mysql"
	"tiktok-demo/cmd/comment/pkg/pack"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/CommentServer"
	"tiktok-demo/shared/kitex_gen/UserServer"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct {
	MysqlManager
	RedisManager
	UserManager
}

type MysqlManager interface {
	AddComment(videoId int64, userId int64, content string) (*mysql.Comment, error)
	DelComment(videoId int64, commentId int64) error
	GetCommentList(videoId int64) ([]*mysql.Comment, error)
}

type RedisManager interface {
	AddComment(videoId int64, commentId int64, userId int64, content string) error
	DelComment(videoId int64, commentId int64) (bool, error)
	GetCommentList(videoId int64) (map[string]map[string]string, error)
}

type UserManager interface {
	GetUserInfo(ctx context.Context, Req *UserServer.DouyinUserRequest, callOptions ...callopt.Option) (r *UserServer.DouyinUserResponse, err error)
	MGetUserInfo(ctx context.Context, req *UserServer.DouyinMUserRequest, callOptions ...callopt.Option) (resp *UserServer.DouyinMUserResponse, err error)
}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *CommentServer.DouyinCommentActionRequest) (resp *CommentServer.DouyinCommentActionResponse, err error) {
	resp = pack.BuildCommentActionResp(nil, nil)

	// 1. check params
	if req.ActionType != 1 && req.ActionType != 2 {
		return pack.BuildCommentActionResp(errno.CommentActionTypeErr, nil), nil
	}
	// 2. action
	if req.ActionType == 1 {
		// 2.1 add comment to mysql
		comment, err := s.MysqlManager.AddComment(req.VideoId, req.UserId, req.CommentText)
		if err != nil {
			klog.Errorf("Mysql AddComment err:%v", err)
			return pack.BuildCommentActionResp(errno.CommentAddErr, nil), nil
		}
		// 2.2 add comment to redis
		err = s.RedisManager.AddComment(req.VideoId, int64(comment.ID), req.UserId, req.CommentText)
		if err != nil {
			klog.Errorf("Redis AddComment err:%v", err)
			return pack.BuildCommentActionResp(errno.CommentAddErr, nil), nil
		}
		// 2.3 get user info from user service
		// todo: using goroutine
		user, err := s.UserManager.GetUserInfo(ctx, &UserServer.DouyinUserRequest{UserId: req.UserId})
		if user.BaseResp.StatusCode != errno.SuccessCode {
			klog.Errorf("GetUserInfo err:%v", err)
			return pack.BuildCommentActionResp(errno.CommentAddErr, nil), nil
		}
		// 2.4 pack comment
		resp := pack.BuildCommentActionResp(nil, pack.CommentInfoConvert(user.User, comment))
		return resp, nil
	} else {
		s.RedisManager.DelComment(req.VideoId, req.CommentId)
		err = s.MysqlManager.DelComment(req.VideoId, req.CommentId)
		if err != nil {
			klog.Errorf("Mysql DelComment err:%v", err)
			return pack.BuildCommentActionResp(errno.CommentDelErr, nil), nil
		}
		time.Sleep(consts.SleepTime)
		s.RedisManager.DelComment(req.VideoId, req.CommentId)
		return resp, nil
	}
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *CommentServer.DouyinCommentListRequest) (resp *CommentServer.DouyinCommentListResponse, err error) {
	userIDs := make([]int64, 0)
	comments := make([]*mysql.Comment, 0)
	// 1. get comment list from redis
	redisComments, _ := s.RedisManager.GetCommentList(req.VideoId)
	if redisComments != nil {
		for commentId := range redisComments {
			id, _ := strconv.ParseInt(commentId, 10, 64)
			commentInfo := redisComments[commentId]
			timestamp, _ := strconv.ParseInt(commentInfo["create_time"], 10, 64)
			userId, _ := strconv.ParseInt(commentInfo["user_id"], 10, 64)
			userIDs = append(userIDs, userId)
			comments = append(comments, &mysql.Comment{
				VideoId: req.VideoId,
				UserId:  userId,
				Content: commentInfo["content"],
				Model: gorm.Model{
					ID:        uint(id),
					CreatedAt: time.Unix(0, int64(timestamp)),
				},
			})
		}
	} else {
		// 2. get comment list from mysql
		mysqlComments, err := s.MysqlManager.GetCommentList(req.VideoId)
		if err != nil {
			klog.Errorf("Mysql GetCommentList err:%v", err)
			return pack.BuildCommentListResp(errno.GetCommentListErr, nil), nil
		}
		// no comment data
		if mysqlComments == nil {
			return pack.BuildCommentListResp(nil, nil), nil
		}
		for _, comment := range mysqlComments {
			userIDs = append(userIDs, comment.UserId)
			comments = append(comments, comment)
		}
	}

	// 3. get user info from user service
	userInfos, err := s.UserManager.MGetUserInfo(ctx, &UserServer.DouyinMUserRequest{UserId: userIDs})
	if err !=nil{
		klog.Errorf("GetUserInfo err:%v", err)
		return pack.BuildCommentListResp(errno.UserRPCErr, nil), nil
	}
	if userInfos.BaseResp.StatusCode != errno.SuccessCode {
		klog.Errorf("GetUserInfo err:%v", err)
		return pack.BuildCommentListResp(errno.UserRPCErr, nil), nil
	}
	// 4. pack comment list
	respCommentList := make([]*CommentServer.Comment, 0)
	for i, comment := range comments {
		respCommentList = append(respCommentList, pack.CommentInfoConvert(userInfos.User[i], comment))
	}
	if len(respCommentList) != len(comments) {
		return pack.BuildCommentListResp(errno.StructConvertFailedErr, nil), nil
	}
	resp = pack.BuildCommentListResp(nil, respCommentList)
	return
}
