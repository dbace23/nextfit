package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nextfit/internal/model"
)

type OrderDetailRepository struct {
	DB *sql.DB
}

func (orderDetailRepository *OrderDetailRepository) CreateOrderDetail(tx *sql.Tx, newOrderDetail model.OrderDetailModel) (model.OrderDetailModel, error) {
	query := `
    INSERT INTO order_details (
        order_id, product_id, quantity
    ) VALUES (?, ?, ?)
    `

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(context.Background(), query,
			newOrderDetail.OrderId,
			newOrderDetail.ProductId,
			newOrderDetail.Quantity,
		)
	} else {
		result, err = orderDetailRepository.DB.ExecContext(context.Background(), query,
			newOrderDetail.OrderId,
			newOrderDetail.ProductId,
			newOrderDetail.Quantity,
		)
	}

	if err != nil {
		fmt.Println(err)
		return model.OrderDetailModel{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return model.OrderDetailModel{}, err
	}

	newOrderDetail.OrderDetailId = int(lastInsertId)
	return newOrderDetail, nil
}

func (orderDetailRepository *OrderDetailRepository) Create(newOrderDetail model.OrderDetailModel) (model.OrderDetailModel, error) {
	return orderDetailRepository.CreateOrderDetail(nil, newOrderDetail)
}

func (orderDetailRepository *OrderDetailRepository) CreateMultipleOrderDetails(tx *sql.Tx, orderDetails []model.OrderDetailModel) ([]model.OrderDetailModel, error) {
	var createdOrderDetails []model.OrderDetailModel

	for _, orderDetail := range orderDetails {
		createdOrderDetail, err := orderDetailRepository.CreateOrderDetail(tx, orderDetail)
		if err != nil {
			return nil, err
		}
		createdOrderDetails = append(createdOrderDetails, createdOrderDetail)
	}

	return createdOrderDetails, nil
}

func (orderDetailRepository *OrderDetailRepository) GetOrderDetailsByOrderId(orderId int) ([]model.OrderDetailModel, error) {
	query := `
    SELECT od.order_detail_id, od.order_id, od.product_id, od.quantity, od.created_at, od.updated_at, od.deleted_at
    FROM order_details od
    WHERE od.order_id = ? AND od.deleted_at IS NULL
    ORDER BY od.created_at ASC
    `

	rows, err := orderDetailRepository.DB.Query(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderDetails []model.OrderDetailModel
	for rows.Next() {
		var orderDetail model.OrderDetailModel
		err := rows.Scan(
			&orderDetail.OrderDetailId,
			&orderDetail.OrderId,
			&orderDetail.ProductId,
			&orderDetail.Quantity,
			&orderDetail.CreatedAt,
			&orderDetail.UpdatedAt,
			&orderDetail.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	return orderDetails, nil
}
