package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type LikeInfo struct {
	VideoID  int64 `json:"video_id"`
	UserID   int64 `json:"user_id"`
	LikeTime int64 `json:"like_time"`
}

func GetLikeInfo(videoID, userID string) (*LikeInfo, error) {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
	result, err := Rs.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	likeInfo := &LikeInfo{}
	if err := json.Unmarshal(result, likeInfo); err != nil {
		return nil, err
	}

	return likeInfo, nil
}

func GetLikeCount(videoID int64) (int64, error) {
	key := fmt.Sprintf("like:count:%d", videoID)
	result, err := Rs.Get(context.Background(), key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return result, nil
}

func IncrLikeCount(videoID string) error {
	key := fmt.Sprintf("like:count:%s", videoID)
	_, err := Rs.Incr(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return nil
}

func DecrLikeCount(videoID string) error {
	key := fmt.Sprintf("like:count:%s", videoID)
	_, err := Rs.Decr(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return nil
}

func SetLikeInfo(videoID, userID, likeTime int64) error {
	key := fmt.Sprintf("like:%d:%d", videoID, userID)
	likeInfo := &LikeInfo{
		VideoID:  videoID,
		UserID:   userID,
		LikeTime: likeTime,
	}
	value, err := json.Marshal(likeInfo)
	if err != nil {
		return err
	}

	err = Rs.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
