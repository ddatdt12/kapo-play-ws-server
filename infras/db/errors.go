package db

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicatedKey = errors.New("duplicated key not allowed")
)
