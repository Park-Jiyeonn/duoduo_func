package consts

import "fmt"

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
