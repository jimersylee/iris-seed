package models

// 书表结构
type Book struct {
	Id           int64  `json:"id" form:"id"`
	Name         string `gorm:"size:256;unique;not null" json:"token" form:"token"`
	Author       string `gorm:"not null;size:256" json:"author" form:"author"`
	Img          string `gorm:"not null;size:256" json:"image" form:"image"`
	PublishTime  string `json:"PublishTime" form:"PublishTime"`
	Page         int    `json:"page" form:"page"`
	ISBN         string `json:"isbn" form:"isbn"`
	Introduction string `json:"introduction" form:"introduction"`
	CreateTime   int64  `gorm:"not null" json:"createTime" form:"createTime"`
}
