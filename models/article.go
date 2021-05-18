package models

import (
	"gorm.io/datatypes"
)

type Article struct {
	Id         int64          `json:"id" form:"id"`
	Title      string         `json:"title" form:"title"`
	Content    string         `json:"content" form:"content"`
	CreateAt   int64          `json:"create_at" form:"create_at"`
	UpdateAt   int64          `json:"update_at" form:"update_at"`
	UserId     int64          `json:"user_id" form:"user_id"`
	Category   string         `json:"category" form:"category"`
	TagsString datatypes.JSON `json:"tags_string" form:"tags_string"`
	ViewTimes  int64          `json:"view_times" form:"view_times"`
	Likes      int64          `json:"likes" form:"likes"`
	Dislikes   int64          `json:"dislikes" form:"dislikes"`
}
