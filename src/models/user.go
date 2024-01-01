package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	JoinedAt time.Time `json:"joinedAt"`
}

func CreateUser(username string) *User {
	return &User{
		Username: username,
	}
}

func (m *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
