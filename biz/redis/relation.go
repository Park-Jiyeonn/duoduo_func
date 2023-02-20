package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// 存储用户关注列表的键：following:{userID}，其中userID是当前用户的ID，存储该用户关注的所有用户ID。
// 存储用户被关注列表的键：followers:{userID}，其中userID是当前用户的ID，存储所有关注该用户的用户ID。
func FollowUser(followerID, followingID string) error {
	// 将followingID添加到followerID的关注列表中
	err := Rs.SAdd(context.Background(), fmt.Sprintf("following:%s", followerID), followingID).Err()
	if err != nil {
		return err
	}

	// 将followerID添加到followingID的被关注列表中
	err = Rs.SAdd(context.Background(), fmt.Sprintf("followers:%s", followingID), followerID).Err()
	if err != nil {
		// 如果添加失败，需要回滚前一个操作，防止出现数据不一致的情况
		Rs.SRem(context.Background(), fmt.Sprintf("following:%s", followerID), followingID)
		return err
	}

	return nil
}
func Unfollow(userID, followingID string) error {
	// 删除用户的关注列表中关注的人
	_, err := Rs.SRem(context.Background(), fmt.Sprintf("following:%s", userID), followingID).Result()
	if err != nil {
		return err
	}

	// 删除被关注者的粉丝列表中的用户
	_, err = Rs.SRem(context.Background(), fmt.Sprintf("followers:%s", followingID), userID).Result()
	if err != nil {
		// 如果删除被关注者的粉丝列表中的用户失败，需要回滚之前删除的关注关系
		_, _ = Rs.SAdd(context.Background(), fmt.Sprintf("following:%s", userID), followingID).Result()
		return err
	}

	return nil
}
func GetFollowing(userID string) ([]int64, error) {
	// 获取该用户关注的所有用户ID
	members, err := Rs.SMembers(context.Background(), fmt.Sprintf("following:%s", userID)).Result()
	if err != nil {
		return nil, err
	}

	following := make([]int64, len(members))
	for i, v := range members {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		following[i] = id
	}
	return following, nil
}

func GetFollowers(userID string) ([]int64, error) {
	// 获取所有关注该用户的用户ID
	members, err := Rs.SMembers(context.Background(), fmt.Sprintf("followers:%s", userID)).Result()
	if err != nil {
		return nil, err
	}

	Followers := make([]int64, len(members))
	for i, v := range members {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		Followers[i] = id
	}
	return Followers, nil
}

func GetMyFriends(userID string) ([]int64, error) {
	var friendIDs []int64
	followingKeys, err := Rs.SMembers(context.Background(), fmt.Sprintf("following:%s", userID)).Result()
	if err != nil {
		return friendIDs, err
	}

	for _, followingKey := range followingKeys {
		followerID, err := strconv.ParseInt(strings.TrimPrefix(followingKey, "following:"), 10, 64)
		if err != nil {
			return friendIDs, err
		}

		isFriend, err := Rs.SIsMember(context.Background(), fmt.Sprintf("following:%d", followerID), fmt.Sprintf("%s", userID)).Result()
		if err != nil {
			return friendIDs, err
		}

		if isFriend {
			friendIDs = append(friendIDs, followerID)
		}
	}

	return friendIDs, nil
}
