package model

type CommentRequest struct {
	Token       string `json:"token"`
	VideoId     int    `json:"video_id"`
	ActionType  string `json:"action_type"`
	CommentText string `json:"comment_text,omitempty"`
	CommentId   int    `json:"comment_id,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

type CommentResponse struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	Comment    CommentInfo `json:"comment,omitempty"`
}

type CommentInfo struct {
	Id         int      `json:"id"`
	User       UserInfo `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}

type CommentListRequest struct {
	Token   string `json:"token"`
	VideoId int    `json:"video_id"`
}
type CommentListResponse struct {
	StatusCode  int           `json:"status_code"`
	StatusMsg   string        `json:"status_msg,omitempty"`
	CommentList []CommentInfo `json:"comment_list,omitempty"`
}
