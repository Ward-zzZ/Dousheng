package tools

import (
	"errors"

	"tiktok-demo/cmd/user/pkg/mysql"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/UserServer"
)

func baseResp(err errno.ErrNo) *UserServer.BaseResp {
	return &UserServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
}

func ParseBaseResp(resp *UserServer.BaseResp) error {
	if resp.StatusCode == errno.Success.ErrCode {
		return nil
	}
	return errno.NewErrNo(resp.StatusCode, resp.StatusMsg)
}

func userInfoResp(err errno.ErrNo, user *UserServer.User) *UserServer.DouyinUserResponse {
	resp := new(UserServer.DouyinUserResponse)
	resp.BaseResp = &UserServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.User = user
	return resp
}

func MuserInfoResp(err errno.ErrNo, user []*UserServer.User) *UserServer.DouyinMUserResponse {
	resp := new(UserServer.DouyinMUserResponse)
	resp.BaseResp = &UserServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.User = user
	return resp
}

// UserInfoConvert convert mysql.User to UserServer.User
func UserInfoConvert(u *mysql.User) *UserServer.User {
	if u == nil {
		return nil
	}
	return &UserServer.User{
		Id:            u.Uid,
		Name:          u.Name,
		FollowCount:   int64(u.FollowingCount),
		FollowerCount: int64(u.FollowerCount),
	}
}

func BuildBaseResp(err error) *UserServer.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func BuilduserInfoResp(err error, u *mysql.User) *UserServer.DouyinUserResponse {
	user := UserInfoConvert(u)
	if err == nil {
		return userInfoResp(errno.Success, user)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return userInfoResp(e, user)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return userInfoResp(s, user)
}

func BuildMuserInfoResp(err error, db_users []*mysql.User) *UserServer.DouyinMUserResponse {
	if len(db_users) == 0 {
		return MuserInfoResp(errno.Success, nil)
	}
	users := make([]*UserServer.User, 0)
	for _, u := range db_users {
		users = append(users, UserInfoConvert(u)) //don't query relation ,isFollow always false
	}
	if len(users) != len(db_users) {
		return MuserInfoResp(errno.StructConvertFailedErr, nil)
	}
	if err == nil {
		return MuserInfoResp(errno.Success, users)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return MuserInfoResp(e, users)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return MuserInfoResp(s, users)
}
