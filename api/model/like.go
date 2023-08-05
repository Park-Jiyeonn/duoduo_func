package model

type LikeRequest struct {
	Token      string `json:"token"`
	VideoId    int    `json:"video_id"`
	ActionType string `json:"action_type"`
	UserId     int    `json:"user_id,omitempty"`
	UserName   string `json:"user_name,omitempty"`
}

type LikeResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type LikeListRequest struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type LikeListResponse struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	VideoList  []VideoInfo `json:"video_list,omitempty"`
}
