package model

type FeedRequest struct {
	LatestTime string `json:"latest_time"`
	Token      string `json:"token"`
	UserId     int    `json:"user_id"`
	UserName   string `json:"user_name"`
}

type VideoInfo struct {
	Id            int      `json:"id"`
	Author        UserInfo `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount int      `json:"favorite_count"`
	CommentCount  int      `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}
type FeedResponse struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	NextTime   int         `json:"next_time"`
	VideoList  []VideoInfo `json:"video_list"`
}

type PublishRequest struct {
	Data     []byte `json:"data"`
	Token    string `json:"token"`
	Title    string `json:"title"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type PublishResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type PublishListRequest struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type PublishListResponse struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	VideoList  []VideoInfo `json:"video_list,omitempty"`
}
