package types

type User struct {
	ID       int    `json:"-" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
