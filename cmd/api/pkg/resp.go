package pkg

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"reflect"
	"tiktok-demo/cmd/api/biz/model/ApiServer"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/CommentServer"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/MessageServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
	"tiktok-demo/shared/kitex_gen/VideoServer"
)

// user

type RegisterResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	User       *ApiServer.User `json:"user"`
}

func SendRegisterResponse(c *app.RequestContext, err error, id int64, token string) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, RegisterResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		UserId:     id,
		Token:      token,
	})
}

func SendUesrInfoResponse(c *app.RequestContext, err error, u *UserServer.User, isFollow bool) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, UserResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		User: func() *ApiServer.User {
			if u == nil {
				return nil
			}
			return &ApiServer.User{
				Id:            u.Id,
				Name:          u.Name,
				FollowCount:   u.FollowCount,
				FollowerCount: u.FollowerCount,
				IsFollow:      isFollow,
			}
		}(),
	})
}

// relation
type RelationActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type RelationListResponse struct {
	StatusCode int32                  `json:"status_code"`
	StatusMsg  string                 `json:"status_msg"`
	UserList   []*RelationServer.User `json:"user_list"`
}

type FriendListResponse struct {
	StatusCode int32                        `json:"status_code"`
	StatusMsg  string                       `json:"status_msg"`
	UserList   []*RelationServer.FriendUser `json:"user_list"`
}

func SendRelationActionResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, RelationActionResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		})
	case *RelationServer.DouyinRelationActionResponse:
		c.JSON(consts.StatusOK, RelationActionResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

func SendRelationListResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, RelationListResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
			UserList:   nil,
		})
	case *RelationServer.DouyinRelationFollowListResponse:
		c.JSON(consts.StatusOK, RelationListResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
			UserList:   value.UserList,
		})
	case *RelationServer.DouyinRelationFollowerListResponse:
		c.JSON(consts.StatusOK, RelationListResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
			UserList:   value.UserList,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

func SendFriendListResponse(c *app.RequestContext, resp interface{} /*, user *UserServer.DouyinMUserResponse, msg *MessageServer.DouyinMessageLatestMsgResponse*/) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, FriendListResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
			UserList:   nil,
		})
	case *RelationServer.DouyinRelationFriendListResponse:
		c.JSON(consts.StatusOK, FriendListResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
			UserList:   value.UserList,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

// favorite
type FavoriteListResponse struct {
	StatusCode int32                   `json:"status_code"`
	StatusMsg  string                  `json:"status_msg"`
	VideoList  []*FavoriteServer.Video `json:"video_list"`
}

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func SendFavoriteActionResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, FavoriteActionResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		})
	case *FavoriteServer.DouyinFavoriteActionResponse:
		c.JSON(consts.StatusOK, FavoriteActionResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

func SendFavoriteListResponse(c *app.RequestContext, err error, videoList []*FavoriteServer.Video) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, FavoriteListResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
	})
}

// comment

type CommentActionResponse struct {
	StatusCode int32                  `json:"status_code"`
	StatusMsg  string                 `json:"status_msg"`
	Comment    *CommentServer.Comment `json:"comment"`
}

type CommentListResponse struct {
	StatusCode  int32                    `json:"status_code"`
	StatusMsg   string                   `json:"status_msg"`
	CommentList []*CommentServer.Comment `json:"comment_list"`
}

func SendCommentActionResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, CommentActionResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
			Comment:    nil,
		})
	case *CommentServer.DouyinCommentActionResponse:
		c.JSON(consts.StatusOK, CommentActionResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
			Comment:    value.Comment,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

func SendCommentListResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, CommentListResponse{
			StatusCode:  Err.ErrCode,
			StatusMsg:   Err.ErrMsg,
			CommentList: nil,
		})
	case *CommentServer.DouyinCommentListResponse:
		c.JSON(consts.StatusOK, CommentListResponse{
			StatusCode:  value.BaseResp.StatusCode,
			StatusMsg:   value.BaseResp.StatusMsg,
			CommentList: value.CommentList,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

// video
type FeedResponse struct {
	StatusCode int32                `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	VideoList  []*VideoServer.Video `json:"video_list"`
}

type PublishActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type PublishListResponse struct {
	StatusCode int32                `json:"status_code"`
	StatusMsg  string               `json:"status_msg"`
	VideoList  []*VideoServer.Video `json:"video_list"`
}

func SendFeedResponse(c *app.RequestContext, err error, videoList []*VideoServer.Video) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, FeedResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
	})
}

func SendPublishActionResponse(c *app.RequestContext, err error, resp interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, PublishActionResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
	})
}

func SendPublishListResponse(c *app.RequestContext, err error, videoList []*VideoServer.Video) {
	Err := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, PublishListResponse{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
		VideoList:  videoList,
	})
}

// message
type MessageChatResponse struct {
	StatusCode  int32                        `json:"status_code"`
	StatusMsg   string                       `json:"status_msg"`
	MessageList []*MessageServer.ChatMessage `json:"message_list"`
}

type MessageActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func SendMessageChatResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, MessageChatResponse{
			StatusCode:  Err.ErrCode,
			StatusMsg:   Err.ErrMsg,
			MessageList: nil,
		})
	case *MessageServer.DouyinMessageChatResponse:
		c.JSON(consts.StatusOK, MessageChatResponse{
			StatusCode:  value.BaseResp.StatusCode,
			StatusMsg:   value.BaseResp.StatusMsg,
			MessageList: value.ChatList,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}

func SendMessageActionResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, MessageActionResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		})
	case *MessageServer.DouyinMessageActionResponse:
		c.JSON(consts.StatusOK, MessageActionResponse{
			StatusCode: value.BaseResp.StatusCode,
			StatusMsg:  value.BaseResp.StatusMsg,
		})
	default:
		hlog.Error("unknown type of response %v", reflect.TypeOf(resp))
	}
}
