package models

type LoginHistory struct {
	ID       int64  `json:"id" form:"id"`
	Ip       string `json:"ip" form:"ip"`
	UserId   int64  `json:"user_id" form:"user_id"`
	CreateAt int64  `json:"create_at" form:"create_at"`
	UpdateAt int64  `json:"update_at" form:"update_at"`
	Hostname string `json:"hostname" form:"hostname"`
}
