package model

import (
	"database/sql"
	"time"
)

type OrderDetailModel struct {
	OrderDetailId int
	OrderId       int
	ProductId     int
	Quantity      int
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
	DeletedAt     sql.NullTime
}
