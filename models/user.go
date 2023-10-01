package models

type User struct {
	Username string
}

func CreateUser(username string) *User {
	return &User{
		Username: username,
	}
}
