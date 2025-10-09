package repository_test

import (
	"context"
	"testing"

	"nextfit/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPendingOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &repository.OrderRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"order_id", "user_id", "address", "paid_amount", "shipping_fee",
		"payment_type", "order_status", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, 2, "Jl. Jakarta", 100000.00, 5000.00, "CASH", "PENDING",
		"2025-10-10 10:00:00", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM orders`).
		WillReturnRows(rows)

	orders, err := repo.GetPendingOrders(context.Background())
	if err != nil {
		t.Fatalf("GetPendingOrders error: %v", err)
	}
	if len(orders) != 1 || orders[0].OrderStatus != "PENDING" {
		t.Fatalf("unexpected result: %+v", orders)
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &repository.OrderRepository{DB: db}

	mock.ExpectExec(`UPDATE orders SET order_status = \?, updated_at = NOW\(\) WHERE order_id = \? AND deleted_at IS NULL`).
		WithArgs("COMPLETED", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.UpdateOrderStatus(context.Background(), 1, "COMPLETED"); err != nil {
		t.Fatalf("UpdateOrderStatus error: %v", err)
	}
}
