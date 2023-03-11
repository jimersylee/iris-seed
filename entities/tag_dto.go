package entities

// TagDTO new tag
type TagDTO struct {
	Name string `json:"name" form:"name"`
	Num  int    `json:"num" form:"num"`
}
