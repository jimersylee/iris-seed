package entities

type LoginDTO struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
