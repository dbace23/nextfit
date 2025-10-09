package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nextfit/internal/model"
)

type OrderRepository struct {
	DB *sql.DB
}

func (orderRepository *OrderRepository) CreateOrder(tx *sql.Tx, newOrder model.OrderModel) (model.OrderModel, error) {
	query := `
    INSERT INTO orders (
        user_id, address, paid_amount, shipping_fee, payment_type, order_status
    ) VALUES (?, ?, ?, ?, ?, ?)
    `

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(context.Background(), query,
			newOrder.UserId,
			newOrder.Address,
			newOrder.PaidAmount,
			newOrder.ShippingFee,
			newOrder.PaymentType,
			newOrder.OrderStatus,
		)
	} else {
		result, err = orderRepository.DB.ExecContext(context.Background(), query,
			newOrder.UserId,
			newOrder.Address,
			newOrder.PaidAmount,
			newOrder.ShippingFee,
			newOrder.PaymentType,
			newOrder.OrderStatus,
		)
	}

	if err != nil {
		fmt.Println(err)
		return model.OrderModel{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return model.OrderModel{}, err
	}

	newOrder.OrderId = int(lastInsertId)
	return newOrder, nil
}

func (orderRepository *OrderRepository) Create(newOrder model.OrderModel) (model.OrderModel, error) {
	return orderRepository.CreateOrder(nil, newOrder)
}

func (orderRepository *OrderRepository) GetOrdersByUserId(userId int) ([]model.OrderModel, error) {
	query := `
    SELECT order_id, user_id, address, paid_amount, shipping_fee, payment_type, order_status, created_at, updated_at, deleted_at
    FROM orders 
    WHERE user_id = ? AND deleted_at IS NULL
    ORDER BY created_at DESC
    `

	rows, err := orderRepository.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderModel
	for rows.Next() {
		var order model.OrderModel
		err := rows.Scan(
			&order.OrderId,
			&order.UserId,
			&order.Address,
			&order.PaidAmount,
			&order.ShippingFee,
			&order.PaymentType,
			&order.OrderStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (orderRepository *OrderRepository) GetOrderById(orderId int) (model.OrderModel, error) {
	query := `
    SELECT order_id, user_id, address, paid_amount, shipping_fee, payment_type, order_status, created_at, updated_at, deleted_at
    FROM orders 
    WHERE order_id = ? AND deleted_at IS NULL
    `

	var order model.OrderModel
	err := orderRepository.DB.QueryRow(query, orderId).Scan(
		&order.OrderId,
		&order.UserId,
		&order.Address,
		&order.PaidAmount,
		&order.ShippingFee,
		&order.PaymentType,
		&order.OrderStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)

	if err != nil {
		return model.OrderModel{}, err
	}

	return order, nil
}

/// order update get pending -> shows then update
func (orderRepository *OrderRepository) GetPendingOrders(ctx context.Context) ([]model.OrderModel, error) {
	const query = `
		SELECT 
			order_id,
			user_id,
			address,
			paid_amount,
			shipping_fee,
			payment_type,
			order_status,
			created_at
		FROM orders
		WHERE order_status = 'PENDING'
		  AND deleted_at IS NULL
		ORDER BY created_at ASC;
	`

	rows, err := orderRepository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderModel
	for rows.Next() {
		var o model.OrderModel
		if err := rows.Scan(
			&o.OrderId,
			&o.UserId,
			&o.Address,
			&o.PaidAmount,
			&o.ShippingFee,
			&o.PaymentType,
			&o.OrderStatus,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}


func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderId int, status string) error {
    const query = `
        UPDATE orders
        SET order_status = ?, updated_at = NOW()
        WHERE order_id = ?
          AND deleted_at IS NULL
          AND order_status = 'PENDING';
   		 `
    res, err := r.DB.ExecContext(ctx, query, status, orderId)
    if err != nil { return err }

    n, err := res.RowsAffected()
    if err != nil { return err }
    if n == 0 {
        return fmt.Errorf("please choose id from list above")
    }
    return nil
}


