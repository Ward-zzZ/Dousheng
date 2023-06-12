package pack

import (
	"errors"

	"tiktok-demo/cmd/comment/pkg/mysql"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/CommentServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
)

func CommentActionResp(err errno.ErrNo, comment *CommentServer.Comment) *CommentServer.DouyinCommentActionResponse {
	resp := new(CommentServer.DouyinCommentActionResponse)
	resp.BaseResp = &CommentServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.Comment = comment
	return resp
}

func CommentListResp(err errno.ErrNo, comments []*CommentServer.Comment) *CommentServer.DouyinCommentListResponse {
	resp := new(CommentServer.DouyinCommentListResponse)
	resp.BaseResp = &CommentServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.CommentList = comments
	return resp
}

func BuildCommentActionResp(err error, comment *CommentServer.Comment) *CommentServer.DouyinCommentActionResponse {
	if err == nil {
		return CommentActionResp(errno.Success, comment)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return CommentActionResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return CommentActionResp(s, nil)
}

func BuildCommentListResp(err error, users []*CommentServer.Comment) *CommentServer.DouyinCommentListResponse {
	if err == nil {
		return CommentListResp(errno.Success, users)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return CommentListResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return CommentListResp(s, nil)
}

// insert UserServer.User into CommentServer.Comment
func CommentInfoConvert(u *UserServer.User, v *mysql.Comment) *CommentServer.Comment {
	if u == nil {
		return nil
	}
	return &CommentServer.Comment{
		Id: int64(v.Model.ID),
		User: &CommentServer.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      false, // query relation rpc in api gateway
		},
		Content:    v.Content,
		CreateDate: v.Model.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
