package models

import (
	"time"
)

type Like struct {
	ID        int64     `json:"id" form:"id"`
	ArticleId int64     `json:"articleId" form:"article_id"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdateAt  time.Time
}
