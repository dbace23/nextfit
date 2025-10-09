package model

import (
	"database/sql"
	"time"
)

type CategoryModel struct {
	CategoryId int
	Name       string
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
	DeletedAt  sql.NullTime
}
