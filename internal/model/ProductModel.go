package model

import (
	"database/sql"
	"time"
)

type ProductModel struct {
	ProductId    int
	ProductName  string
	CategoryId   int
	SellingPrice float64
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
}
