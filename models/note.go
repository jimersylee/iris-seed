package models

type Note struct {
	Id         int64  `json:"id" form:"id"`
	BookId     int64  `json:"bookId" form:"bookId"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
	UpdateTime int64  `gorm:"not null" json:"updateTime" form:"updateTime"`
	Content    string `gorm:"text;not null" json:"content" form:"content"`
}
