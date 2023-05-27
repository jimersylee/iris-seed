package models

import (
	"time"
)

type Role struct {
	ID        int64     `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdateAt  time.Time
}
