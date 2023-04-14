package models

import (
	"time"
)

type Admin struct {
	ID              int64  `json:"id" form:"id"`
	Name            string `json:"name" form:"name"`
	Email           string
	Avatar          string `json:"avatar" form:"avatar"`
	EmailVerifiedAt time.Time
	Password        string
	RememberToken   string
	CreatedAt       time.Time `json:"created_at" form:"created_at"`
	UpdateAt        time.Time
}
