package state

import (
	"time"
)

type User struct {
	ID        uint
	Username string
	Password string
	Email string
	ApiKey string
	State int
	EmailedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type UserChain struct {
	UserID uint
	Chain string
}

type Record struct {
	ID uint
	Chain string
}

