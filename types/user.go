package types

import "fmt"

type User struct {
	ID       int    `json:"-" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (u User) ValidateUser() error {
	if len(u.Password) > 5 || len(u.Password) < 13 {
		return fmt.Errorf("your password is weak %d , write strong password ", len(u.Password))
	}
	if len(u.Username) < 3 {
		return fmt.Errorf("username must have 3 characters! %d", len(u.Username))
	}
	return nil
}
