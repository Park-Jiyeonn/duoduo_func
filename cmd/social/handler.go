package main

import (
	"context"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/dal/redis"
	"simple_tiktok/kitex_gen/base"
	social "simple_tiktok/kitex_gen/social"
	"simple_tiktok/pkg/errno"
	"sort"
	"strconv"
	"time"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct{}

// FollowAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowAction(ctx context.Context, req *social.FollowRequest) (resp *social.FollowResponse, err error) {
	// TODO: Your code here...
	resp = social.NewFollowResponse()
	userId, toUserId1 := req.GetUserId(), req.GetToUserId()
	toUserId, _ := strconv.ParseInt(toUserId1, 10, 64)
	// 1. 判断需要操作的对象在缓存中是否存在
	if redis.FollowIsExists(ctx, userId) == 0 {
		// 缓存中不存在用户粉丝列表
		userList, err := db.GetFollowList(ctx, userId)
		if err != nil {
			return nil, err
		}
		if len(userList) > 0 {
			kv := make([]string, 0)
			for _, user := range userList {
				kv = append(kv, strconv.FormatInt(user, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFollowList(ctx, userId, kv...) {
				return resp, errno.NewErrNo("redis 缓存失败")
			}
		}
	}

	if redis.FollowerIsExists(ctx, toUserId) == 0 {
		// 缓存中不存在用户粉丝列表
		userList, err := db.GetFollowerList(ctx, toUserId)
		if err != nil {
			return nil, err
		}
		if len(userList) > 0 {
			kv := make([]string, 0)
			for _, user := range userList {
				kv = append(kv, strconv.FormatInt(user, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFollowerList(ctx, toUserId, kv...) {
				return resp, errno.NewErrNo("redis 缓存失败")
			}
		}
	}

	if redis.UserIsExists(ctx, userId) == 0 {
		// 当关注的用户不存在
		user, err := db.GetUserById(ctx, userId)
		if err != nil {
			return resp, err
		}
		if !redis.SetUserInfo(ctx, user) {
			return resp, errno.NewErrNo("redis 缓存失败")
		}
	}

	if redis.UserIsExists(ctx, toUserId) == 0 {
		// 当被关注的用户不存在
		user, err := db.GetUserById(ctx, toUserId)
		if err != nil {
			return resp, err
		}
		if !redis.SetUserInfo(ctx, user) {
			return resp, errno.NewErrNo("redis 缓存失败")
		}
	}

	// 2. 判断操作类型
	var action int64
	follow := redis.IsFollow(ctx, userId, toUserId)
	if follow && req.GetActionType() == "2" {
		// 已关注并取消关注
		action = -1
	} else if !follow && req.GetActionType() == "1" {
		// 关注
		action = 1
	} else {
		// 关系和操作不匹配
		return resp, errno.NewErrNo("关系和操作不匹配")
	}

	// 3. 原子性的进行操作
	// 更新用户关注列表
	if !redis.FollowAction(ctx, userId, toUserId, action) {
		return resp, errno.NewErrNo("redis 更新用户关注列表 failed")
	}
	// 更新用户关注数量
	if !redis.IncrUserField(ctx, userId, "follow_count", action) {
		return resp, errno.NewErrNo("redis 更新用户关注数量 failed")
	}
	// 更新用户被关注数量
	if !redis.IncrUserField(ctx, toUserId, "follower_count", action) {
		return resp, errno.NewErrNo("redis 更新用户被关注数量 failed")
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return resp, nil
}

// GetFollowList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowList(ctx context.Context, req *social.FollowingListRequest) (resp *social.FollowingListResponse, err error) {
	// TODO: Your code here...
	resp = social.NewFollowingListResponse()
	userId := req.GetUserId()
	var followidList []int64
	if redis.FollowIsExists(ctx, userId) == 0 {
		// 缓存中不存在查询用户粉丝列表
		followidList, err := db.GetFollowList(ctx, userId)
		if err != nil {
			return nil, err
		}
		if len(followidList) > 0 {
			kv := make([]string, 0)
			for _, user := range followidList {
				kv = append(kv, strconv.FormatInt(user, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFollowList(ctx, userId, kv...) {
				return resp, errno.NewErrNo("redis 缓存 failed")
			}
		}
	} else {
		// 缓存中存在用户粉丝列表
		followidList = redis.GetFollowList(ctx, userId)
	}

	// 遍历查找查询用户所关注的用户
	var retUsers []*base.UserInfo
	if len(followidList) > 0 {
		for _, followid := range followidList {
			var follow = &model.User{}
			if redis.UserIsExists(ctx, followid) == 0 {
				// 如果要查询的用户不在缓存中
				follow, err = db.GetUserById(ctx, followid)
				if err != nil {
					return resp, err
				}
				redis.SetUserInfo(ctx, follow)
			} else {
				// 要查询的用户位于缓存中
				follow, err = redis.GetUserInfo(ctx, followid)
				if err != nil {
					return resp, err
				}
			}
			retUsers = append(retUsers, &base.UserInfo{
				Id:            int64(follow.ID),
				Name:          follow.Name,
				FollowCount:   follow.FollowCount,
				FollowerCount: follow.FollowerCount,
				IsFollow:      false,
			})
		}
	}
	resp.SetUserList(retUsers)
	resp.StatusCode = 0
	message := ""
	resp.StatusMsg = &message
	message = "success"
	return resp, nil
}

// GetFollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFollowerList(ctx context.Context, req *social.FollowerListRequest) (resp *social.FollowerListResponse, err error) {
	// TODO: Your code here...
	resp = social.NewFollowerListResponse()
	message := ""
	resp.StatusMsg = &message
	userId := req.GetUserId()

	var followeridList []int64
	if redis.FollowerIsExists(ctx, userId) == 0 {
		followeridList, err := db.GetFollowerList(ctx, userId)
		if err != nil {
			return nil, err
		}
		if len(followeridList) > 0 {
			kv := make([]string, 0)
			for _, user := range followeridList {
				kv = append(kv, strconv.FormatInt(user, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFollowerList(ctx, userId, kv...) {
				return resp, errno.NewErrNo("redis 缓存失败")
			}
		}
	} else {
		followeridList = redis.GetFollowerList(ctx, userId)
	}

	// 遍历查找查询用户所关注的用户
	var followerList []*base.UserInfo
	if len(followeridList) > 0 {
		for _, followerid := range followeridList {
			var follower = &model.User{}
			if redis.UserIsExists(ctx, followerid) == 0 {
				// 如果要查询的用户不在缓存中
				follower, err = db.GetUserById(ctx, followerid)
				if err != nil {
					return resp, err
				}
				redis.SetUserInfo(ctx, follower)
			} else {
				// 要查询的用户位于缓存中
				follower, err = redis.GetUserInfo(ctx, followerid)
				if err != nil {
					return resp, err
				}
			}
			followerList = append(followerList, &base.UserInfo{
				Id:            int64(follower.ID),
				Name:          follower.Name,
				FollowCount:   follower.FollowCount,
				FollowerCount: follower.FollowerCount,
				IsFollow:      false,
			})
		}
	}
	resp.SetUserList(followerList)
	resp.StatusCode = 0
	message = "success"
	return resp, nil
}

// GetFriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) GetFriendList(ctx context.Context, req *social.FriendListRequest) (resp *social.FriendListResponse, err error) {
	// TODO: Your code here...
	resp = social.NewFriendListResponse()
	message := ""
	resp.StatusMsg = &message
	userId := req.GetUserId()
	var followList []int64
	// Try to get the user follower list from cache.
	// if missing, fill cache from dao.
	if redis.FollowIsExists(ctx, userId) == 0 {
		followList, err := db.GetFollowList(ctx, userId)
		if err != nil {
			return nil, err
		}
		if len(followList) > 0 {
			kv := make([]string, 0)
			for _, user := range followList {
				kv = append(kv, strconv.FormatInt(user, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFollowList(ctx, userId, kv...) {
				return resp, errno.NewErrNo("redis 缓存失败")
			}
		}
	} else {
		followList = redis.GetFollowList(ctx, userId)
	}

	var friends []*base.UserInfo
	for _, followerId := range followList {
		// Add the follower from cache

		if redis.FollowIsExists(ctx, followerId) == 0 {
			userList, err := db.GetFollowList(ctx, followerId)
			if err != nil {
				return nil, err
			}
			if len(userList) > 0 {
				kv := make([]string, 0)
				for _, user := range userList {
					kv = append(kv, strconv.FormatInt(user, 10))
					kv = append(kv, "1")
				}
				if !redis.SetFollowList(ctx, followerId, kv...) {
					return resp, errno.NewErrNo("redis 缓存 failed")
				}
			}
		}
		isFriend := redis.IsFollow(ctx, followerId, userId)
		if isFriend == false {
			continue
		}
		// Try to get follower UserMessage from cache
		user := &model.User{}
		if redis.UserIsExists(ctx, followerId) == 0 {
			user, err = db.GetUserById(ctx, followerId)
			if err != nil {
				return resp, err
			}
			if !redis.SetUserInfo(ctx, user) {
				return resp, errno.NewErrNo("redis 缓存 failed")
			}
		} else {
			user, err = redis.GetUserInfo(ctx, followerId)
			if err != nil {
				return resp, err
			}
		}

		friend := &base.UserInfo{
			Id:            int64(user.ID),
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true,
		}

		friends = append(friends, friend)
	}

	resp.UserList = friends
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
		msg.ToUserId = sendMessages[i].ToUserId
		msg.FromUserId = sendMessages[i].UserId
		msg.Content = sendMessages[i].Content
		msg.CreateTime = &sendMessages[i].PublishDate
		MessageList[i] = &msg
	}
	for i := len(sendMessages); i < len(sendMessages)+len(reciveMessage); i++ {
		var msg social.Message
		msg.Id = int64(reciveMessage[i-len(sendMessages)].ID)
		msg.ToUserId = reciveMessage[i-len(sendMessages)].ToUserId
		msg.FromUserId = reciveMessage[i-len(sendMessages)].UserId
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
	newMessage := &model.Message{
		UserId:      *req.UserId,
		ToUserId:    req.ToUserId,
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
