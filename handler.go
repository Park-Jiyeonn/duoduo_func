package main

import (
	"context"
	base "simple_tiktok/kitex_gen/base"
)

// BaseServiceImpl implements the last service interface defined in the IDL.
type BaseServiceImpl struct{}

// UserRegister implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) UserRegister(ctx context.Context, request *base.RegisterRequest) (resp *base.RegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// UserLogin implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) UserLogin(ctx context.Context, request *base.LoginRequest) (resp *base.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetUserInfo(ctx context.Context, request *base.UserInfoRequest) (resp *base.UserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideoList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetVideoList(ctx context.Context, request *base.FeedRequest) (resp *base.FeedResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishAction implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) PublishAction(ctx context.Context, request *base.PublishRequest) (resp *base.PublishResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPublishList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetPublishList(ctx context.Context, request *base.PublishListRequest) (resp *base.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}
