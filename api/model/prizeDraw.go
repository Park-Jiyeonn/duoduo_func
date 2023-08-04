package model

type Treasure struct {
	Value  int `json:"value"`
	Weight int `json:"weight"`
}

type TreasureResponse struct {
	StatusCode int        `json:"status_code"`
	StatusMsg  string     `json:"status_msg"`
	Treasure   []Treasure `json:"treasure"`
	Capacity   int        `json:"capacity"`
	Limit      int        `json:"limit"`
}
