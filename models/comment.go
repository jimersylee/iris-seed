package models

type Comment struct {
	Id       int64  `json:"id" form:"id"`
	Content  string `json:"content" form:"content"`
	CreateAt int64  `json:"create_at" form:"create_at"`
	UpdateAt int64  `json:"update_at" form:"update_at"`
	UserId   int64  `json:"user_id" form:"user_id"`
}
