package main

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"tiktok-demo/cmd/relation/pkg/mq"
	"tiktok-demo/cmd/relation/pkg/mysql"
	"tiktok-demo/cmd/relation/pkg/pack"
	// "tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/MessageServer"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct {
	MysqlManager
	RedisManager
	UserManager
	MessageManager
}

type MysqlManager interface {
	AddFollow(userId int64, toUserId int64) error
	DelFollow(userId int64, toUserId int64) error
	QueryRelation(userId int64, toUserId int64) (bool, error)
	GetFansList(userId int64) ([]*mysql.Relation, error)
	GetFollowList(userId int64) ([]*mysql.Relation, error)
	GetFollowSet(userId int64) (map[int64]struct{}, error)
	GetFriendList(userId int64) ([]*mysql.Relation, error)
}

type RedisManager interface {
	AddFollow(c context.Context, uid int64, ids int64) (bool, error)
	UnFollow(c context.Context, uid int64, tid int64) (bool, error)
	SetFollow(c context.Context, uid int64, ids []int64) (bool, error)
	SetFans(c context.Context, uid int64, ids []int64) (bool, error)
	QueryRelation(c context.Context, uid int64, tid int64) (bool, error)
	QueryFollow(c context.Context, uid int64) ([]int64, error)
	QueryFans(c context.Context, uid int64) ([]int64, error)
}

type UserManager interface {
	ChangeUserFollowCount(ctx context.Context, req *UserServer.DouyinChangeUserFollowRequest, callOptions ...callopt.Option) (resp *UserServer.BaseResp, err error)
	MGetUserInfo(ctx context.Context, req *UserServer.DouyinMUserRequest, callOptions ...callopt.Option) (resp *UserServer.DouyinMUserResponse, err error)
}

type MessageManager interface {
	MessageLatestMsg(ctx context.Context, req *MessageServer.DouyinMessageLatestMsgRequest, callOptions ...callopt.Option) (resp *MessageServer.DouyinMessageLatestMsgResponse, err error)
}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *RelationServer.DouyinRelationActionRequest) (resp *RelationServer.DouyinRelationActionResponse, err error) {
	resp = pack.BuildrelationActionResp(nil)

	// 1. check params
	if req.ActionType != 1 && req.ActionType != 2 {
		resp = pack.BuildrelationActionResp(errno.ActionTypeErr)
		return resp, nil
	}
	if req.UserId == req.ToUserId {
		return resp, nil
	}
	isFollow, _ := s.QueryRelation(ctx, &RelationServer.DouyinQueryRelationRequest{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	})
	if isFollow.BaseResp.StatusCode == errno.SuccessCode {
		if isFollow.IsFollow && req.ActionType == 1 {
			resp = pack.BuildrelationActionResp(errno.FollowActionErr.WithMessage("请勿重复关注"))
			return resp, nil
		}
		if !isFollow.IsFollow && req.ActionType == 2 {
			resp = pack.BuildrelationActionResp(errno.FollowActionErr.WithMessage("请勿重复取消关注"))
			return resp, nil
		}
	}
	// 2. action/undo action
	// 2.1 First: rpc user service change follow count
	UserRPCResp, _ := s.UserManager.ChangeUserFollowCount(ctx, &UserServer.DouyinChangeUserFollowRequest{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
		IsFollow: req.ActionType == 1,
	})
	if UserRPCResp.StatusCode != errno.SuccessCode {
		klog.Errorf("UserRPC ChangeUserFollowCount err:%v", UserRPCResp)
		resp = pack.BuildrelationActionResp(errno.UserRPCErr)
		return resp, nil
	}

	// 2.2 Second: mysql add/del follow,,using delayed double deletion
	msg := strings.Builder{}
	msg.WriteString(strconv.FormatInt(req.UserId, 10))
	msg.WriteString(":")
	msg.WriteString(strconv.FormatInt(req.ToUserId, 10))
	if req.ActionType == 1 {
		s.RedisManager.UnFollow(ctx, req.UserId, req.ToUserId)
		err = mq.AddActor.Publish(ctx, msg.String())
	} else {
		s.RedisManager.UnFollow(ctx, req.UserId, req.ToUserId)
		err = mq.DelActor.Publish(ctx, msg.String())
	}
	klog.Infof("RelationAction msg:%s", msg.String())
	if err != nil {
		klog.Errorf("RelationAction mq publish err:%v", err)
		resp = pack.BuildrelationActionResp(errno.FollowActionErr)
		return resp, nil
	}
	return resp, nil
}

// MGetRelationFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetRelationFollowList(ctx context.Context, req *RelationServer.DouyinRelationFollowListRequest) (resp *RelationServer.DouyinRelationFollowListResponse, err error) {
	followIDs := make([]int64, 0)
	// 1. check redis
	ids, err := s.RedisManager.QueryFollow(ctx, req.UserId)
	if ids == nil || err != nil {
		// 2. check mysql
		follows, err := s.MysqlManager.GetFollowList(req.UserId)
		if err != nil {
			klog.Errorf("GetRelationFollowList mysql err:%v", err)
			return pack.BuildgetFollowListResp(errno.GetFollowListErr, nil), nil
		}
		for _, follow := range follows {
			followIDs = append(followIDs, follow.ToUserID)
		}
		// update redis
		if success, err := s.RedisManager.SetFollow(ctx, req.UserId, followIDs); !success || err != nil {
			klog.Errorf("Update %d Redis RelationFollowList redis err:%v", req.UserId, err)
		}
	} else {
		followIDs = append(followIDs, ids...)
	}
	// no follow
	if len(followIDs) == 0 {
		return pack.BuildgetFollowListResp(errno.Success, nil), nil
	}
	// 3. rpc user service get user info
	UserRPCResp, _ := s.UserManager.MGetUserInfo(ctx, &UserServer.DouyinMUserRequest{
		UserId: followIDs,
	})
	if UserRPCResp.BaseResp.StatusCode != errno.SuccessCode {
		klog.Errorf("UserRPC MGetUserInfo err:%v", UserRPCResp.BaseResp)
		return pack.BuildgetFollowListResp(errno.UserRPCErr, nil), nil
	}
	// 4. pack response
	followUsers := make([]*RelationServer.User, 0)
	for _, u := range UserRPCResp.User {
		followUsers = append(followUsers, &RelationServer.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      true,
		})
	}
	if len(followUsers) != len(UserRPCResp.User) {
		klog.Errorf("MGetRelationFollowList len(followUsers) != len(UserRPCResp.Users)")
		return pack.BuildgetFollowListResp(errno.StructConvertFailedErr, nil), nil
	}
	return pack.BuildgetFollowListResp(nil, followUsers), nil
}

// MGetUserRelationFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetUserRelationFollowerList(ctx context.Context, req *RelationServer.DouyinRelationFollowerListRequest) (resp *RelationServer.DouyinRelationFollowerListResponse, err error) {
	followerIDs := make([]int64, 0)
	// 1. check redis
	ids, err := s.RedisManager.QueryFans(ctx, req.UserId)
	if ids == nil || err != nil {
		// 2. check mysql
		followers, err := s.MysqlManager.GetFansList(req.UserId)
		if err != nil {
			klog.Errorf("GetUserRelationFollowerList mysql err:%v", err)
			return pack.BuildgetFollowerListResp(errno.GetFollowerListErr, nil), nil
		}
		for _, follower := range followers {
			followerIDs = append(followerIDs, follower.UserID)
		}
		// update redis
		if success, err := s.RedisManager.SetFans(ctx, req.UserId, followerIDs); !success || err != nil {
			klog.Errorf("Update %d Redis RelationFollowerList redis err:%v", req.UserId, err)
		}
	} else {
		followerIDs = append(followerIDs, ids...)
	}
	// no follower
	if len(followerIDs) == 0 {
		return pack.BuildgetFollowerListResp(errno.Success, nil), nil
	}
	var wg sync.WaitGroup
	wg.Add(2)
	var UsersRPCResp *UserServer.DouyinMUserResponse
	followSet := make(map[int64]struct{})
	// 粉丝的信息
	go func() {
		UsersRPCResp, err = s.UserManager.MGetUserInfo(ctx, &UserServer.DouyinMUserRequest{
			UserId: followerIDs,
		})
		if UsersRPCResp.BaseResp.StatusCode != 0 {
			klog.Errorf("UserRPC MGetUserInfo err:%v", UsersRPCResp.BaseResp.StatusMsg)
		}
		wg.Done()
	}()
	// 当前用户关注的人，用于判断当前用户是否关注了这些粉丝
	go func() {
		follows, _ := s.MGetRelationFollowList(ctx, &RelationServer.DouyinRelationFollowListRequest{
			UserId: req.UserId,
		})
		if follows.BaseResp.StatusCode != 0 {
			klog.Errorf("MGetRelationFollowList err:%v", follows.BaseResp.StatusMsg)
		}
		for _, follow := range follows.UserList {
			followSet[follow.Id] = struct{}{}
		}
		wg.Done()
	}()
	wg.Wait()
	followerUsers := make([]*RelationServer.User, 0)
	for _, u := range UsersRPCResp.User {
		_, ok := followSet[u.Id] // Check whether the user also follows this fan
		followerUsers = append(followerUsers, &RelationServer.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      ok,
		})
	}
	if len(followerUsers) != len(UsersRPCResp.User) {
		klog.Errorf("MGetUserRelationFollowerList len(followerUsers) != len(UsersRPCResp.Users)")
		return pack.BuildgetFollowerListResp(errno.StructConvertFailedErr, nil), nil
	}
	return pack.BuildgetFollowerListResp(nil, followerUsers), nil
}

// QueryRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) QueryRelation(ctx context.Context, req *RelationServer.DouyinQueryRelationRequest) (resp *RelationServer.DouyinQueryRelationResponse, err error) {
	if req.UserId == req.ToUserId {
		return pack.BuildrelationQueryResp(nil, true), nil
	}

	isFollow, err := s.RedisManager.QueryRelation(ctx, req.UserId, req.ToUserId)
	if err == nil {
		return pack.BuildrelationQueryResp(nil, isFollow), nil
	}

	// check relation mysql
	follows, err := s.MysqlManager.GetFollowList(req.UserId)
	if err != nil {
		return pack.BuildrelationQueryResp(errno.QueryFollowErr, false), nil
	}
	followsIDs := make([]int64, 0)
	for _, follow := range follows {
		followsIDs = append(followsIDs, follow.ToUserID)
		if follow.ToUserID == req.ToUserId {
			isFollow = true
		}
	}
	if len(followsIDs) != 0 {
		s.RedisManager.SetFollow(ctx, req.UserId, followsIDs)
	}
	return pack.BuildrelationQueryResp(nil, isFollow), nil
}

// MGetRelationFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetRelationFriendList(ctx context.Context, req *RelationServer.DouyinRelationFriendListRequest) (resp *RelationServer.DouyinRelationFriendListResponse, err error) {
	friendList, err := s.MysqlManager.GetFriendList(req.UserId)
	if err != nil {
		klog.Errorf("MGetRelationFriendList mysql err:%v", err)
		return pack.BuildgetFriendListResp(errno.InternalServerErr, nil), nil
	}
	friendIDs := make([]int64, 0)
	for _, friend := range friendList {
		friendIDs = append(friendIDs, friend.ToUserID)
	}
	if len(friendIDs) == 0 {
		return pack.BuildgetFriendListResp(errno.Success.WithMessage("no friend"), nil), nil
	}

	// user rpc and message rpc
	// 开协程并发请求
	var wg sync.WaitGroup
	wg.Add(2)
	var UsersRPCResp *UserServer.DouyinMUserResponse
	var MessageRPCResp *MessageServer.DouyinMessageLatestMsgResponse
	go func() {
		UsersRPCResp, err = s.UserManager.MGetUserInfo(ctx, &UserServer.DouyinMUserRequest{
			UserId: friendIDs,
		})
		if UsersRPCResp.BaseResp.StatusCode != 0 {
			klog.Errorf("UserRPC MGetUserInfo err:%v", UsersRPCResp.BaseResp.StatusMsg)
		}
		wg.Done()
	}()
	go func() {
		MessageRPCResp, err = s.MessageManager.MessageLatestMsg(ctx, &MessageServer.DouyinMessageLatestMsgRequest{
			UserId:       req.UserId,
			ToUserIdList: friendIDs,
		})
		if MessageRPCResp.BaseResp.StatusCode != 0 {
			klog.Errorf("MessageRPC GetLatestMsg err:%v", MessageRPCResp.BaseResp.StatusMsg)
		}
		wg.Done()
	}()
	wg.Wait()
	friends, err := pack.FriendUserInfoConvert(UsersRPCResp, MessageRPCResp)
	if err != nil {
		klog.Errorf("FriendUserInfoConvert err:%v", err)
		return pack.BuildgetFriendListResp(errno.StructConvertFailedErr, nil), nil
	}
	return pack.BuildgetFriendListResp(nil, friends), nil
}
