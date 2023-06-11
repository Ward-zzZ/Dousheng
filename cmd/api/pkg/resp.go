package pkg

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"reflect"
	"tiktok-demo/cmd/api/biz/model/ApiServer"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/FavoriteServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
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

func SendRelationActionResponse(c *app.RequestContext, resp interface{}) {
	switch value := resp.(type) {
	case error:
		Err := errno.ConvertErr(value)
		c.JSON(consts.StatusOK, RelationListResponse{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		})
	case *RelationServer.DouyinRelationActionResponse:
		c.JSON(consts.StatusOK, RelationListResponse{
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
