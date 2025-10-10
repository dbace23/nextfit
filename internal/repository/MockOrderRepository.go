package repository

import (
	"context"
	"database/sql"
	"nextfit/internal/model"
)

type OrderStore interface {
	CreateOrder(tx *sql.Tx, newOrder model.OrderModel) (model.OrderModel, error)
	GetOrdersByUserId(userId int) ([]model.OrderModel, error)
	GetOrderById(orderId int) (model.OrderModel, error)
	GetPendingOrders(ctx context.Context) ([]model.OrderModel, error)
	UpdateOrderStatus(ctx context.Context, orderId int, status string) error
}

type OrderDetailStore interface {
	CreateMultipleOrderDetails(tx *sql.Tx, details []model.OrderDetailModel) ([]model.OrderDetailModel, error)
	GetOrderDetailsByOrderId(orderId int) ([]model.OrderDetailModel, error)
}

type ProductStore interface {
	FindByProductId(id int) (model.ProductModel, error)
}

type UserStore interface {
	FindByUserId(id int) (model.UserModel, error)
}
