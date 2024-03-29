package pack

import (
	"errors"

	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/MessageServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
)

func UserInfoConvert(u *UserServer.User, isFollow bool) *RelationServer.User {
	if u == nil {
		return nil
	}
	return &RelationServer.User{
		Id:            u.Id,
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      isFollow,
	}
}

func FriendUserInfoConvert(u *UserServer.DouyinMUserResponse, f *MessageServer.DouyinMessageLatestMsgResponse) ([]*RelationServer.FriendUser, error) {
	if u == nil || f == nil {
		return nil, errors.New("user or msg is nil")
	}
	user := u.User
	typeList := f.TypeList
	contentList := f.ContentList

	var friendUsers []*RelationServer.FriendUser
	for i, _ := range user {
		friendUser := &RelationServer.FriendUser{
			Id:            user[i].Id,
			Name:          user[i].Name,
			FollowCount:   user[i].FollowCount,
			FollowerCount: user[i].FollowerCount,
			IsFollow:      true,
			MsgType:       int64(typeList[i]),
			Message:       contentList[i],
		}
		friendUsers = append(friendUsers, friendUser)
	}
	return friendUsers, nil
}

func relationActionResp(err errno.ErrNo) *RelationServer.DouyinRelationActionResponse {
	resp := new(RelationServer.DouyinRelationActionResponse)
	resp.BaseResp = &RelationServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	return resp
}

func relationQueryResp(err errno.ErrNo, isFollow bool) *RelationServer.DouyinQueryRelationResponse {
	resp := new(RelationServer.DouyinQueryRelationResponse)
	resp.BaseResp = &RelationServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.IsFollow = isFollow
	return resp
}

func getFollowListResp(err errno.ErrNo, users []*RelationServer.User) *RelationServer.DouyinRelationFollowListResponse {
	resp := new(RelationServer.DouyinRelationFollowListResponse)
	resp.BaseResp = &RelationServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.UserList = users
	return resp
}

func getFollowerListResp(err errno.ErrNo, users []*RelationServer.User) *RelationServer.DouyinRelationFollowerListResponse {
	resp := new(RelationServer.DouyinRelationFollowerListResponse)
	resp.BaseResp = &RelationServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.UserList = users
	return resp
}

func getFriendListResp(err errno.ErrNo, users []*RelationServer.FriendUser) *RelationServer.DouyinRelationFriendListResponse {
	resp := new(RelationServer.DouyinRelationFriendListResponse)
	resp.BaseResp = &RelationServer.BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
	resp.UserList = users
	return resp
}

func BuildrelationActionResp(err error) *RelationServer.DouyinRelationActionResponse {
	if err == nil {
		return relationActionResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return relationActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return relationActionResp(s)
}

func BuildrelationQueryResp(err error, isFollow bool) *RelationServer.DouyinQueryRelationResponse {
	if err == nil {
		return relationQueryResp(errno.Success, isFollow)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return relationQueryResp(e, isFollow)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return relationQueryResp(s, isFollow)
}

func BuildgetFollowListResp(err error, users []*RelationServer.User) *RelationServer.DouyinRelationFollowListResponse {
	if err == nil {
		return getFollowListResp(errno.Success, users)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return getFollowListResp(e, nil)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return getFollowListResp(s, nil)
}

func BuildgetFollowerListResp(err error, users []*RelationServer.User) *RelationServer.DouyinRelationFollowerListResponse {
	if err == nil {
		return getFollowerListResp(errno.Success, users)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return getFollowerListResp(e, nil)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return getFollowerListResp(s, nil)
}

func BuildgetFriendListResp(err error, users []*RelationServer.FriendUser) *RelationServer.DouyinRelationFriendListResponse {
	if err == nil {
		return getFriendListResp(errno.Success, users)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return getFriendListResp(e, nil)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return getFriendListResp(s, nil)
}
