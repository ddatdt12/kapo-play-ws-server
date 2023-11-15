package models

type User struct {
	ID       int64
	Username string
}

func CreateUser(username string) *User {
	return &User{
		Username: username,
	}
}
