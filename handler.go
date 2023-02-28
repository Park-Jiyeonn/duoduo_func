package main

import (
	"context"
	base "simple_tiktok/kitex_gen/base"
	interact "simple_tiktok/kitex_gen/interact"
	social "simple_tiktok/kitex_gen/social"
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

// LikeAction implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) LikeAction(ctx context.Context, request *interact.LikeRequest) (resp *interact.LikeResponse, err error) {
	// TODO: Your code here...
	return
}

// GetLikeList implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) GetLikeList(ctx context.Context, request *interact.LikeListRequest) (resp *interact.LikeListResponse, err error) {
	// TODO: Your code here...
	return
}

// CommentAction implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) CommentAction(ctx context.Context, request *interact.CommentRequest) (resp *interact.CommentResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCommentList implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) GetCommentList(ctx context.Context, request *interact.CommentListRequest) (resp *interact.CommentListResponse, err error) {
	// TODO: Your code here...
	return
}

// FollowAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowAction(ctx context.Context, request *social.FollowRequest) (resp *social.FollowResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, request *social.FollowingListRequest) (resp *social.FollowingListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, request *social.FollowerListRequest) (resp *social.FollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, request *social.FriendListRequest) (resp *social.FriendListResponse, err error) {
	// TODO: Your code here...
	return
}

// MessageChat implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MessageChat(ctx context.Context, req *social.MessageChatReq) (resp *social.MessageChatResp, err error) {
	// TODO: Your code here...
	return
}

// MessageAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MessageAction(ctx context.Context, req *social.MessageActionReq) (resp *social.MessageActionResp, err error) {
	// TODO: Your code here...
	return
}

// FollowAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowAction(ctx context.Context, request *social.FollowRequest) (resp *social.FollowResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowList(ctx context.Context, request *social.FollowingListRequest) (resp *social.FollowingListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowerList(ctx context.Context, request *social.FollowerListRequest) (resp *social.FollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFriendList(ctx context.Context, request *social.FriendListRequest) (resp *social.FriendListResponse, err error) {
	// TODO: Your code here...
	return
}

// MessageChat implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) MessageChat(ctx context.Context, req *social.MessageChatReq) (resp *social.MessageChatResp, err error) {
	// TODO: Your code here...
	return
}

// MessageAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) MessageAction(ctx context.Context, req *social.MessageActionReq) (resp *social.MessageActionResp, err error) {
	// TODO: Your code here...
	return
}
