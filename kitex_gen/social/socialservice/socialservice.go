// Code generated by Kitex v0.4.4. DO NOT EDIT.

package socialservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	social "simple_tiktok/kitex_gen/social"
)

func serviceInfo() *kitex.ServiceInfo {
	return socialServiceServiceInfo
}

var socialServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "SocialService"
	handlerType := (*social.SocialService)(nil)
	methods := map[string]kitex.MethodInfo{
		"FollowAction":    kitex.NewMethodInfo(followActionHandler, newSocialServiceFollowActionArgs, newSocialServiceFollowActionResult, false),
		"GetFollowList":   kitex.NewMethodInfo(getFollowListHandler, newSocialServiceGetFollowListArgs, newSocialServiceGetFollowListResult, false),
		"GetFollowerList": kitex.NewMethodInfo(getFollowerListHandler, newSocialServiceGetFollowerListArgs, newSocialServiceGetFollowerListResult, false),
		"GetFriendList":   kitex.NewMethodInfo(getFriendListHandler, newSocialServiceGetFriendListArgs, newSocialServiceGetFriendListResult, false),
		"MessageChat":     kitex.NewMethodInfo(messageChatHandler, newSocialServiceMessageChatArgs, newSocialServiceMessageChatResult, false),
		"MessageAction":   kitex.NewMethodInfo(messageActionHandler, newSocialServiceMessageActionArgs, newSocialServiceMessageActionResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "social",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func followActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceFollowActionArgs)
	realResult := result.(*social.SocialServiceFollowActionResult)
	success, err := handler.(social.SocialService).FollowAction(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceFollowActionArgs() interface{} {
	return social.NewSocialServiceFollowActionArgs()
}

func newSocialServiceFollowActionResult() interface{} {
	return social.NewSocialServiceFollowActionResult()
}

func getFollowListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetFollowListArgs)
	realResult := result.(*social.SocialServiceGetFollowListResult)
	success, err := handler.(social.SocialService).GetFollowList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetFollowListArgs() interface{} {
	return social.NewSocialServiceGetFollowListArgs()
}

func newSocialServiceGetFollowListResult() interface{} {
	return social.NewSocialServiceGetFollowListResult()
}

func getFollowerListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetFollowerListArgs)
	realResult := result.(*social.SocialServiceGetFollowerListResult)
	success, err := handler.(social.SocialService).GetFollowerList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetFollowerListArgs() interface{} {
	return social.NewSocialServiceGetFollowerListArgs()
}

func newSocialServiceGetFollowerListResult() interface{} {
	return social.NewSocialServiceGetFollowerListResult()
}

func getFriendListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetFriendListArgs)
	realResult := result.(*social.SocialServiceGetFriendListResult)
	success, err := handler.(social.SocialService).GetFriendList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetFriendListArgs() interface{} {
	return social.NewSocialServiceGetFriendListArgs()
}

func newSocialServiceGetFriendListResult() interface{} {
	return social.NewSocialServiceGetFriendListResult()
}

func messageChatHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceMessageChatArgs)
	realResult := result.(*social.SocialServiceMessageChatResult)
	success, err := handler.(social.SocialService).MessageChat(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceMessageChatArgs() interface{} {
	return social.NewSocialServiceMessageChatArgs()
}

func newSocialServiceMessageChatResult() interface{} {
	return social.NewSocialServiceMessageChatResult()
}

func messageActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceMessageActionArgs)
	realResult := result.(*social.SocialServiceMessageActionResult)
	success, err := handler.(social.SocialService).MessageAction(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceMessageActionArgs() interface{} {
	return social.NewSocialServiceMessageActionArgs()
}

func newSocialServiceMessageActionResult() interface{} {
	return social.NewSocialServiceMessageActionResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) FollowAction(ctx context.Context, request *social.FollowRequest) (r *social.FollowResponse, err error) {
	var _args social.SocialServiceFollowActionArgs
	_args.Request = request
	var _result social.SocialServiceFollowActionResult
	if err = p.c.Call(ctx, "FollowAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowList(ctx context.Context, request *social.FollowingListRequest) (r *social.FollowingListResponse, err error) {
	var _args social.SocialServiceGetFollowListArgs
	_args.Request = request
	var _result social.SocialServiceGetFollowListResult
	if err = p.c.Call(ctx, "GetFollowList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowerList(ctx context.Context, request *social.FollowerListRequest) (r *social.FollowerListResponse, err error) {
	var _args social.SocialServiceGetFollowerListArgs
	_args.Request = request
	var _result social.SocialServiceGetFollowerListResult
	if err = p.c.Call(ctx, "GetFollowerList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFriendList(ctx context.Context, request *social.FriendListRequest) (r *social.FriendListResponse, err error) {
	var _args social.SocialServiceGetFriendListArgs
	_args.Request = request
	var _result social.SocialServiceGetFriendListResult
	if err = p.c.Call(ctx, "GetFriendList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) MessageChat(ctx context.Context, req *social.MessageChatReq) (r *social.MessageChatResp, err error) {
	var _args social.SocialServiceMessageChatArgs
	_args.Req = req
	var _result social.SocialServiceMessageChatResult
	if err = p.c.Call(ctx, "MessageChat", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) MessageAction(ctx context.Context, req *social.MessageActionReq) (r *social.MessageActionResp, err error) {
	var _args social.SocialServiceMessageActionArgs
	_args.Req = req
	var _result social.SocialServiceMessageActionResult
	if err = p.c.Call(ctx, "MessageAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}