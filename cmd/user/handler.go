package main

import (
	"context"
	UserServer "tiktok-demo/shared/kitex_gen/UserServer"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *UserServer.DouyinUserRegisterRequest) (resp *UserServer.DouyinUserRegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *UserServer.DouyinUserLoginRequest) (resp *UserServer.DouyinUserLoginResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *UserServer.DouyinUserRequest) (resp *UserServer.DouyinUserResponse, err error) {
	// TODO: Your code here...
	return
}

// MGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUserInfo(ctx context.Context, req *UserServer.DouyinMUserRequest) (resp *UserServer.DouyinMUserResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangeUserFollowCount implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangeUserFollowCount(ctx context.Context, req *UserServer.DouyinChangeUserFollowRequest) (resp *UserServer.BaseResp, err error) {
	// TODO: Your code here...
	return
}
