package consts

import (
	"fmt"
	"strconv"
)

func GetUserInfoKey(userId int) string {
	return fmt.Sprintf("user_info_%d", userId)
}
func GetUserLikeListKey(userID int) string {
	return fmt.Sprintf("user_like_list_%d", userID)
}

func GetVideoMsgKey(videoID int) string {
	return fmt.Sprintf("video_message_%d", videoID)
}

func GetIDFromUserMsgKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[10:], 10, 64)
	return
}

func GetIDFromUserLikeListKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[15:], 10, 64)
	return
}

func GetIDFromVideoMsgKey(key string) (id int64) {
	id, _ = strconv.ParseInt(key[14:], 10, 64)
	return
}

func GetFavoriteLmtKey(ipaddr string) string {
	return "favorite_limit_" + ipaddr
}
