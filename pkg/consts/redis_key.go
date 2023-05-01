package consts

import (
	"fmt"
	"strconv"
)

func GetUserInfoKey(userId int64) string {
	return fmt.Sprintf("user_info_%d", userId)
}

func GetUserFollowKey(userId int64) string {
	return fmt.Sprintf("user_follow_%d", userId)
}

func GetUserFollowerKey(userID int64) string {
	return fmt.Sprintf("user_follower_%d", userID)
}

func GetUserLikeKey(userID int64) string {
	return fmt.Sprintf("user_like_%d", userID)
}

func GetVideoMsgKey(videoID int64) string {
	return fmt.Sprintf("video_message_%d", videoID)
}

func GetIDFromUserMsgKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[10:], 10, 64)
	return
}

func GetIDFromUserLikeListKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[10:], 10, 64)
	return
}

func GetIDFromUserFollowListKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[12:], 10, 64)
	return
}

func GetIDFromVideoMsgKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[14:], 10, 64)
	return
}
