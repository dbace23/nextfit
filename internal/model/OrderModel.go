package model

import (
	"database/sql"
	"time"
)

type OrderModel struct {
	OrderId     int
	UserId      int
	Address     string
	PaidAmount  float64
	ShippingFee float64
	PaymentType string
	OrderStatus string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
	DeletedAt   sql.NullTime
}
