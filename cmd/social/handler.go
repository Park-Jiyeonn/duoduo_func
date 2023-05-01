package main

import (
	"context"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/redis"
	"simple_tiktok/kitex_gen/base"
	social "simple_tiktok/kitex_gen/social"
	"simple_tiktok/pkg/errno"
	"sort"
	"time"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct{}

// FollowAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowAction(ctx context.Context, request *social.FollowRequest) (resp *social.FollowResponse, err error) {
	// TODO: Your code here...
	resp = new(social.FollowResponse)
	if request.ActionType != "1" {
		err = redis.Unfollow(ctx, *request.UserId, request.ToUserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("取消关注失败")
		}
	} else {
		err = redis.FollowUser(ctx, *request.UserId, request.ToUserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("关注失败")
		}
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return resp, nil
}

// GetFollowList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowList(ctx context.Context, request *social.FollowingListRequest) (resp *social.FollowingListResponse, err error) {
	// TODO: Your code here...
	resp = new(social.FollowingListResponse)
	Followings, err := redis.GetFollowing(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询关注的人失败")
	}

	var users []*base.UserInfo

	for _, ToUserID := range Followings {
		isFollow, err := redis.HasFollowed(ctx, request.UserId, uint(ToUserID))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("查询是否关注失败")
		}
		to, _ := db.GetUserById(ctx, ToUserID)
		f := &base.UserInfo{
			Id:            ToUserID,
			Name:          to.Name,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      isFollow,
		}
		users = append(users, f)
	}

	//fmt.Println(users)

	resp.UserList = users
	resp.StatusCode = 0
	message := ""
	resp.StatusMsg = &message
	message = "success"
	return resp, nil
}

// GetFollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowerList(ctx context.Context, request *social.FollowerListRequest) (resp *social.FollowerListResponse, err error) {
	// TODO: Your code here...
	resp = new(social.FollowerListResponse)

	message := ""
	resp.StatusMsg = &message
	Fans, err := redis.GetFans(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询粉丝失败")
	}
	var users []*base.UserInfo

	for _, ToUserID := range Fans {
		isFollow, err := redis.HasFollowed(ctx, request.UserId, uint(ToUserID))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("查询是否关注失败")
		}
		to, _ := db.GetUserById(ctx, ToUserID)
		f := &base.UserInfo{
			Id:            ToUserID,
			Name:          to.Name,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      isFollow,
		}
		users = append(users, f)
	}

	resp.UserList = users
	resp.StatusCode = 0
	message = "success"
	return resp, nil
}

// GetFriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFriendList(ctx context.Context, request *social.FriendListRequest) (resp *social.FriendListResponse, err error) {
	// TODO: Your code here...
	resp = new(social.FriendListResponse)
	message := ""
	resp.StatusMsg = &message
	Followings, err := redis.GetMyFriends(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询朋友失败")
	}
	var users []*base.UserInfo

	for _, userID := range Followings {
		to, _ := db.GetUserById(ctx, userID)
		f := &base.UserInfo{
			Id:            userID,
			Name:          to.Name,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      true,
		}
		users = append(users, f)
	}

	resp.UserList = users
	resp.StatusCode = 0
	message = "success"
	return resp, nil
}

// MessageChat implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) MessageChat(ctx context.Context, req *social.MessageChatReq) (resp *social.MessageChatResp, err error) {
	// TODO: Your code here...
	resp = new(social.MessageChatResp)

	mes := ""
	resp.StatusMsg = &mes

	sendMessages, err := db.QuerryMessageByID(ctx, *req.UserId, req.ToUserId, req.PreMsgTime)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询消息失败")
	}
	reciveMessage, err := db.QuerryMessageByID(ctx, req.ToUserId, *req.UserId, req.PreMsgTime)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询消息失败")
	}

	MessageList := make([]*social.Message, len(sendMessages)+len(reciveMessage))
	for i := 0; i < len(sendMessages); i++ {
		var msg social.Message
		msg.Id = int64(sendMessages[i].ID)
		msg.ToUserId = sendMessages[i].ToUserID
		msg.FromUserId = sendMessages[i].UserID
		msg.Content = sendMessages[i].Content
		msg.CreateTime = &sendMessages[i].PublishDate
		MessageList[i] = &msg
	}
	for i := len(sendMessages); i < len(sendMessages)+len(reciveMessage); i++ {
		var msg social.Message
		msg.Id = int64(reciveMessage[i-len(sendMessages)].ID)
		msg.ToUserId = reciveMessage[i-len(sendMessages)].ToUserID
		msg.FromUserId = reciveMessage[i-len(sendMessages)].UserID
		msg.Content = reciveMessage[i-len(sendMessages)].Content
		msg.CreateTime = &reciveMessage[i-len(sendMessages)].PublishDate
		MessageList[i] = &msg
	}

	sort.Slice(MessageList, func(i, j int) bool {
		return *MessageList[i].CreateTime < *MessageList[j].CreateTime
	})

	resp.StatusCode = 0
	mes = "success"
	resp.MessageList = MessageList
	return resp, nil
}

// MessageAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) MessageAction(ctx context.Context, req *social.MessageActionReq) (resp *social.MessageActionResp, err error) {
	// TODO: Your code here...
	resp = new(social.MessageActionResp)
	mes := ""
	resp.StatusMsg = &mes
	newMessage := &db.Message{
		UserID:      *req.UserId,
		ToUserID:    req.ToUserId,
		Content:     req.Content,
		PublishDate: time.Now().Unix(),
	}
	err = db.CreateMessage(ctx, newMessage)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建消息失败")
	}

	resp.StatusCode = 0
	mes = "success"
	return resp, nil
}
