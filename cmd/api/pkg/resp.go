package pkg

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok-demo/cmd/api/biz/model/ApiServer"
	"tiktok-demo/shared/errno"
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
		User: &ApiServer.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      isFollow,
		},
	})
}
