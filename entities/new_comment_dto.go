package entities

// NewCommentDto comment dto
type NewCommentDto struct {
	Name string `json:"content" form:"content"`
}
