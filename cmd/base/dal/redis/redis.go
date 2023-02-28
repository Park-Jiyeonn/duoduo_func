package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"

	"encoding/json"
	"strings"
)

type LikeInfo struct {
	VideoID string `json:"video_id"`
	UserID  string `json:"user_id"`
}

func GetLikeInfo(ctx context.Context, videoID, userID string) (*LikeInfo, error) {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
	result, err := Rs.Get(ctx, key).Bytes()
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
func GetLikedVideos(ctx context.Context, userID int64) ([]int64, error) {
	var likedVideos []int64
	cursor := uint64(0)
	pattern := fmt.Sprintf("like:*:%d", userID)

	for {
		// 扫描 redis 中所有匹配 pattern 的键
		keys, nextCursor, err := Rs.Scan(ctx, cursor, pattern, 10).Result()
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

func GetLikeCount(ctx context.Context, videoID int64) (int64, error) {
	key := fmt.Sprintf("like:count:%d", videoID)
	result, err := Rs.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return result, nil
}

func IncrLikeCount(ctx context.Context, videoID string) error {
	key := fmt.Sprintf("like:count:%s", videoID)
	_, err := Rs.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func DecrLikeCount(ctx context.Context, videoID string) error {
	key := fmt.Sprintf("like:count:%s", videoID)
	_, err := Rs.Decr(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func SetLikeInfo(ctx context.Context, videoID, userID string) error {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
	likeInfo := &LikeInfo{
		VideoID: videoID,
		UserID:  userID,
	}
	value, err := json.Marshal(likeInfo)
	if err != nil {
		return err
	}

	err = Rs.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DelLikeInfo 删除指定视频和用户的点赞信息
func DelLikeInfo(ctx context.Context, videoID, userID string) error {
	key := fmt.Sprintf("like:%s:%s", videoID, userID)
	_, err := Rs.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete like info: %v", err)
	}
	return nil
}

// InitLikeCount DelLikeCount Redis中增加和删除两个新的键
func InitLikeCount(ctx context.Context, videoID int64) error {
	key := fmt.Sprintf("video:%d:like", videoID)
	return Rs.Set(ctx, key, 0, 0).Err()
}

func DelLikeCount(ctx context.Context, videoID int64) error {
	key := fmt.Sprintf("video:%d:like", videoID)
	return Rs.Del(ctx, key).Err()
}
