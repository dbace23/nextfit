package model

import (
	"database/sql"
	"time"
)

type UserModel struct {
	UserId         int
	FullName       string
	Password       string
	HashedPassword string
	Email          string
	IsAdmin        bool
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
	DeletedAt      sql.NullTime
}
