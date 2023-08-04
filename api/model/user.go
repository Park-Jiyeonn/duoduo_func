package model

type User struct {
	Username string `json:"username" query:"username"`
	Password string `json:"password" query:"password"`
}

type RegisterResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int    `json:"user_id"`
	Token      string `json:"token"`
}

type LoginResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserId     int    `json:"user_id,omitempty"`
	Token      string `json:"token,omitempty"`
}
