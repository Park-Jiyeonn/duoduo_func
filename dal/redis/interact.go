package redis

import (
	"context"
	"fmt"
	"strconv"
)

//func HasLiked(ctx context.Context, userID int64, videoID string) (bool, error) {
//	key := fmt.Sprintf("user:%d:likes", userID) // 构造的key是userID，看看这个集合中有没有videoID
//	ret, err := Rs.SIsMember(ctx, key, videoID).Result()
//	if err != nil {
//		return false, err
//	}
//	return ret, nil
//}

// GetLikedVideos 查询用户点赞过的视频
//func GetLikedVideos(ctx context.Context, userID int64) ([]int64, error) {
//	var likedVideos []int64
//	key := fmt.Sprintf("user:%d:likes", userID)
//
//	ret, err := Rs.SMembers(ctx, key).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	for _, val := range ret {
//		num, _ := strconv.ParseInt(val, 10, 64)
//		likedVideos = append(likedVideos, num)
//	}
//	return likedVideos, nil
//}

func GetLikeCount(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("video:%d:like_count", videoID)
	ret, err := Rs.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	sum, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

//func SetLikeInfo(ctx context.Context, userID int64, videoID string) error {
//	// 包含两个操作，一，集合的加入，二，string的修改
//	key := fmt.Sprintf("user:%d:likes", userID)
//	err := Rs.SAdd(ctx, key, videoID).Err()
//	if err != nil {
//		return err
//	}
//
//	key = fmt.Sprintf("video:%s:like_count", videoID)
//	err = Rs.Incr(ctx, key).Err()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func DelLikeInfo(ctx context.Context, userID int64, videoID string) error {
//	key := fmt.Sprintf("user:%d:likes", userID)
//	err := Rs.SRem(ctx, key, videoID).Err()
//	if err != nil {
//		return err
//	}
//
//	key = fmt.Sprintf("video:%s:like_count", videoID)
//	err = Rs.Decr(ctx, key).Err()
//	if err != nil {
//		return err
//	}
//	return nil
//}
