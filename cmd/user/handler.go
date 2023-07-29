package main

import (
	"context"
	"strconv"
	"strings"

	"tiktok-demo/cmd/user/pkg/mq"
	"tiktok-demo/cmd/user/pkg/mysql"
	"tiktok-demo/shared/errno"
	"tiktok-demo/shared/kitex_gen/RelationServer"
	"tiktok-demo/shared/kitex_gen/UserServer"
	"tiktok-demo/shared/tools"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	MysqlManager
	RedisManager
	RealtionManager
	EncryptManager
}

type MysqlManager interface {
	CreateUser(user *mysql.User) (*mysql.User, error)
	GetUserByID(id int64) (*mysql.User, error)
	GetUserByName(name string) (*mysql.User, error)
	GetUsersByID(ids []int64) ([]*mysql.User, error)
	FollowUser(userID, followerID int64) error
	UnFollowUser(userID, followerID int64) error
}

type RedisManager interface {
	DeleteRelation(c context.Context, uid int64, tid int64) (bool, error)
	GetUserInfo(c context.Context, uid int64) (*mysql.User, error)
	InsertUserInfo(c context.Context, user *mysql.User) error
}

type RealtionManager interface {
	QueryRelation(ctx context.Context, req *RelationServer.DouyinQueryRelationRequest, callOptions ...callopt.Option) (*RelationServer.DouyinQueryRelationResponse, error)
}

type EncryptManager interface {
	EncryptPassword(code string) string
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *UserServer.DouyinUserRegisterRequest) (resp *UserServer.DouyinUserRegisterResponse, err error) {
	resp = new(UserServer.DouyinUserRegisterResponse)
	resp.UserId = -1 // -1 means register failed
	// 1、check data valid
	if len(req.Username) == 0 || len(req.Password) == 0 || len(req.Username) > 32 || len(req.Password) > 32 {
		resp.BaseResp = tools.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	// 2、check username exist
	_, err = s.MysqlManager.GetUserByName(req.Username)
	if err != nil {
		if err != errno.UserNotExistErr {
			klog.Error("get user by name failed, err: ", err)
			resp.BaseResp = tools.BuildBaseResp(errno.MysqlErr)
			return resp, nil
		}
		// 3、create user
		usr, err := s.MysqlManager.CreateUser(&mysql.User{
			Name:     req.Username,
			Password: req.Password,
		})
		if err != nil {
			klog.Error("create user failed, err: ", err)
			resp.BaseResp = tools.BuildBaseResp(errno.MysqlErr)
			return resp, nil
		}
		// 4、return user id
		resp.BaseResp = tools.BuildBaseResp(errno.Success)
		resp.UserId = usr.Uid
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(errno.UserAlreadyExistErr)
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *UserServer.DouyinUserLoginRequest) (resp *UserServer.DouyinUserLoginResponse, err error) {
	resp = new(UserServer.DouyinUserLoginResponse)
	resp.UserId = -1 // -1 means login failed
	// 1、check data valid
	if len(req.Username) == 0 || len(req.Password) == 0 || len(req.Username) > 32 || len(req.Password) > 32 {
		resp.BaseResp = tools.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}
	// 2、check user exist
	usr, err := s.MysqlManager.GetUserByName(req.Username)
	if err != nil {
		if err != errno.UserNotExistErr {
			klog.Error("get password by name failed, err: ", err)
			resp.BaseResp = tools.BuildBaseResp(errno.MysqlErr)
			return resp, nil
		}
		resp.BaseResp = tools.BuildBaseResp(errno.UserNotExistErr)
		return resp, nil
	}
	// 3、check password
	cryPwd := s.EncryptManager.EncryptPassword(req.Password)
	if usr.Password != cryPwd {
		klog.Error("%s login err: password not match", req.Username)
		resp.BaseResp = tools.BuildBaseResp(errno.AuthorizationFailedErr)
		return resp, nil
	}
	// 4、return user id
	resp.BaseResp = tools.BuildBaseResp(errno.Success)
	resp.UserId = usr.Uid
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *UserServer.DouyinUserRequest) (resp *UserServer.DouyinUserResponse, err error) {
	resp = new(UserServer.DouyinUserResponse)
	// 1、get user info from redis
	usr, err := s.RedisManager.GetUserInfo(ctx, req.UserId)
	resp = tools.BuilduserInfoResp(errno.Success, usr)
	klog.Debug("get user info from redis, resp: %+v", resp)
	if err != nil {
		klog.Error("get user info from redis failed, err: ", err)
		// 2、get user info from mysql
		usr, err := s.MysqlManager.GetUserByID(req.UserId)
		if err == errno.UserNotExistErr {
			resp = tools.BuilduserInfoResp(errno.UserNotExistErr, nil)
			return resp, nil
		}
		if err != nil {
			klog.Error("get user info from mysql failed, err: ", err)
			resp = tools.BuilduserInfoResp(errno.MysqlErr, nil)
			return resp, nil
		}
		// 3、return user info
		resp = tools.BuilduserInfoResp(errno.Success, usr)
		klog.Debug("get user info from mysql, resp: %+v", resp)
		go func() {
			// 4、update user info to redis async
			if err := s.RedisManager.InsertUserInfo(ctx, usr); err != nil {
				klog.Error("insert user info to redis failed, err: ", err)
			}
		}()
	}
	return resp, nil
}

// MGetUserInfo implements the UserServiceImpl interface.
// TODO: refactor to avoid call GetUserInfo
func (s *UserServiceImpl) MGetUserInfo(ctx context.Context, req *UserServer.DouyinMUserRequest) (resp *UserServer.DouyinMUserResponse, err error) {
	resp = new(UserServer.DouyinMUserResponse)
	var u []*UserServer.User
	singReq := new(UserServer.DouyinUserRequest)
	for _, id := range req.UserId {
		singReq.UserId = id
		userInfo, err := s.GetUserInfo(ctx, singReq)
		if err != nil {
			klog.Error("get user info failed, err: ", err)
			return nil, errno.FindUserErr
		}
		u = append(u, userInfo.User)
	}
	resp.BaseResp = tools.BuildBaseResp(errno.Success)
	resp.User = u
	return resp, nil
}

// ChangeUserFollowCount implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangeUserFollowCount(ctx context.Context, req *UserServer.DouyinChangeUserFollowRequest) (resp *UserServer.BaseResp, err error) {
	resp = tools.BuildBaseResp(errno.Success)
	// 1、can not follow self
	if req.UserId == req.ToUserId {
		return resp, nil
	}
	// 2、assemble into a msg, send to rabbitmq
	msg := strings.Builder{}
	msg.WriteString(strconv.Itoa(int(req.UserId)))
	msg.WriteString(":")
	msg.WriteString(strconv.Itoa(int(req.ToUserId)))
	// 3、query relation
	flag, err := s.RealtionManager.QueryRelation(ctx, &RelationServer.DouyinQueryRelationRequest{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		klog.Error("query relation failed, err: ", err)
		resp = tools.BuildBaseResp(err)
		return resp, nil
	}
	if flag.BaseResp.StatusCode != 0 {
		resp = tools.BuildBaseResp(errno.NewErrNo(flag.BaseResp.StatusCode, flag.BaseResp.StatusMsg))
		return resp, nil
	}
	// 3、update db,using delayed double deletion
	if req.IsFollow {
		// follow
		if flag.IsFollow {
			return resp, nil
		}
		// 3.2 put into rabbitmq
		err = mq.AddActor.Publish(ctx, msg.String())
		// 3.3 double deletion
		s.RedisManager.DeleteRelation(ctx, req.UserId, req.ToUserId)

	} else {
		// unfollow
		if !flag.IsFollow {
			return resp, nil
		}
		// 3.2 put into rabbitmq
		err = mq.DelActor.Publish(ctx, msg.String())
		// 3.3 double deletion
		s.RedisManager.DeleteRelation(ctx, req.UserId, req.ToUserId)
	}
	if err != nil {
		return tools.BuildBaseResp(errno.ChangeUserFollowCountErr), nil
	}
	return resp, nil
}
