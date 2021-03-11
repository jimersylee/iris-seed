package entities

type TagDTO struct {
	Name string `json:"name" form:"name"`
	Num  int    `json:"num" form:"num"`
}
