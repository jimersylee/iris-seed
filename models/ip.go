package models

import (
	"time"
)

type Ip struct {
	ID           int64  `json:"id" form:"id"`
	Ip           string `json:"ip" form:"ip"`
	Port         int    `json:"port" form:"port"`
	RequestTimes int64  `json:"request_times" form:"request_times"`
	//状态 1：可用 0：不可用
	Status   int       `json:"status" form:"status"`
	CreateAt time.Time `json:"create_at" form:"create_at"`
	UpdateAt time.Time
}
