package models

import (
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
)

type Template struct {
	ID          uint
	Title       string
	Description string
	Cover       string
	IsPublic    bool
	CreatorID   uint
	Creator     User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   types.NullTime
}
