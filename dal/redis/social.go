package redis

import (
	"context"
	"fmt"
	"strconv"
)

// 存储用户被关注列表的键：followers:{userID}，其中userID是当前用户的ID，存储所有关注该用户的用户ID。
func FollowUser(ctx context.Context, UserID int64, ToUserID string) error {
	// 将 ToUserID 添加到 UserID 的关注列表中
	key := fmt.Sprintf("user:%d:follow", UserID)
	err := Rs.SAdd(ctx, key, ToUserID).Err()
	if err != nil {
		return err
	}

	// 将followerID添加到followingID的被关注列表中
	key = fmt.Sprintf("user:%s:fans", ToUserID)
	err = Rs.SAdd(ctx, key, UserID).Err()
	if err != nil {
		return err
	}

	// 如果对方也关注了我，那么我们马上就成为了好友
	key = fmt.Sprintf("user:%s:follow", ToUserID)
	if ret, _ := Rs.SIsMember(ctx, key, UserID).Result(); ret {
		key1 := fmt.Sprintf("user:%d:friend", UserID)
		key2 := fmt.Sprintf("user:%s:friend", ToUserID)
		err = Rs.SAdd(ctx, key1, ToUserID).Err()
		if err != nil {
			return err
		}
		err = Rs.SAdd(ctx, key2, UserID).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func Unfollow(ctx context.Context, userID int64, ToUserID string) error {
	// 删除用户的关注列表中关注的人
	key := fmt.Sprintf("user:%d:follow", userID)
	err := Rs.SRem(ctx, key, ToUserID).Err()
	if err != nil {
		return err
	}

	// 删除 ToUserID 的粉丝列表中的用户
	key = fmt.Sprintf("user:%s:fans", ToUserID)
	err = Rs.SRem(ctx, key, userID).Err()
	if err != nil {
		return err
	}

	// 如果对方也关注了我，那么我们互相失去好友关系
	key = fmt.Sprintf("user:%s:follow", ToUserID)
	if ret, _ := Rs.SIsMember(ctx, key, userID).Result(); ret {
		key1 := fmt.Sprintf("user:%d:friend", userID)
		key2 := fmt.Sprintf("user:%s:friend", ToUserID)
		err = Rs.SRem(ctx, key1, ToUserID).Err()
		if err != nil {
			return err
		}
		err = Rs.SRem(ctx, key2, userID).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func GetFollowing(ctx context.Context, userID int64) ([]int64, error) {
	// 获取该用户关注的所有用户ID
	key := fmt.Sprintf("user:%d:follow", userID)
	members, err := Rs.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var following []int64
	for _, v := range members {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		following = append(following, id)
	}
	return following, nil
}

func GetFans(ctx context.Context, userID int64) ([]int64, error) {
	key := fmt.Sprintf("user:%d:fans", userID)
	// 获取所有关注该用户的用户ID
	members, err := Rs.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var fans []int64
	for _, v := range members {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		fans = append(fans, id)
	}
	return fans, nil
}

func GetMyFriends(ctx context.Context, userID int64) ([]int64, error) {
	key := fmt.Sprintf("user:%d:friend", userID)
	ret, err := Rs.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var Friends []int64
	for _, val := range ret {
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		Friends = append(Friends, num)
	}
	return Friends, nil
}

func HasFollowed(ctx context.Context, userID int64, ToUserID uint) (bool, error) {
	key := fmt.Sprintf("user:%d:follow", userID)
	return Rs.SIsMember(ctx, key, ToUserID).Result()
}
