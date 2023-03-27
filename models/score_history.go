package models

// ScoreHistory score history table
type ScoreHistory struct {
	ID          int64 `json:"id" form:"id"`
	ChangeScore int64 `json:"change_score" form:"change_score"`
	UserId      int64 `json:"user_id" form:"user_id"`
	CreateAt    int64 `json:"create_at" form:"create_at"`
	UpdateAt    int64 `json:"update_at" form:"update_at"`
}
