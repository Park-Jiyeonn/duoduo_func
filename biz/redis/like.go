package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

type LikeInfo struct {
	VideoID  string `json:"video_id"`
	UserID   string `json:"user_id"`
	LikeTime string `json:"like_time"`
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

// GetLikedVideos 查询用户点赞过的视频
func GetLikedVideos(userID int64) ([]int64, error) {
	var likedVideos []int64
	cursor := uint64(0)
	pattern := fmt.Sprintf("like:*:%d", userID)

	for {
		// 扫描 redis 中所有匹配 pattern 的键
		keys, nextCursor, err := Rs.Scan(context.Background(), cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		// 遍历匹配的键，取出视频 ID
		for _, key := range keys {
			parts := strings.Split(key, ":")
			if len(parts) != 3 {
				continue
			}
			videoIDStr := parts[1]
			videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
			if err != nil {
				continue
			}
			likedVideos = append(likedVideos, videoID)
		}

		// 如果 nextCursor 返回 0，说明遍历完成
		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return likedVideos, nil
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

func SetLikeInfo(videoID, userID, likeTime string) error {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
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

// DelLikeInfo 删除指定视频和用户的点赞信息
func DelLikeInfo(videoID, userID string) error {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
	_, err := Rs.Del(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete like info: %v", err)
	}
	return nil
}

// InitLikeCount DelLikeCount Redis中增加和删除两个新的键
func InitLikeCount(videoID int64) error {
	key := fmt.Sprintf("video:%d:like", videoID)
	return Rs.Set(context.Background(), key, 0, 0).Err()
}

func DelLikeCount(videoID int64) error {
	key := fmt.Sprintf("video:%d:like", videoID)
	return Rs.Del(context.Background(), key).Err()
}
