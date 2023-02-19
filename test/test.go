package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type LikeInfo struct {
	VideoID  int64 `json:"video_id"`
	UserID   int64 `json:"user_id"`
	LikeTime int64 `json:"like_time"`
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{rdb, context.Background()}
}

func (c *RedisClient) Close() {
	c.client.Close()
}

func (c *RedisClient) GetLikeInfo(videoID, userID int64) (*LikeInfo, error) {
	key := fmt.Sprintf("like:%d:%d", videoID, userID)
	result, err := c.client.Get(c.ctx, key).Bytes()
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

func (c *RedisClient) GetLikeCount(videoID int64) (int64, error) {
	key := fmt.Sprintf("like:count:%d", videoID)
	result, err := c.client.Get(c.ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return result, nil
}

func (c *RedisClient) IncrLikeCount(videoID int64) error {
	key := fmt.Sprintf("like:count:%d", videoID)
	_, err := c.client.Incr(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) DecrLikeCount(videoID int64) error {
	key := fmt.Sprintf("like:count:%d", videoID)
	_, err := c.client.Decr(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) SetLikeInfo(videoID, userID, likeTime int64) error {
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

	err = c.client.Set(c.ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
