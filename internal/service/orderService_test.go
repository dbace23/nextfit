package service_test

import (
	"context"
	"testing"

	"nextfit/internal/repository"
	"nextfit/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrderService_GetPendingOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &repository.OrderRepository{DB: db}
	srv := &service.OrderService{OrderRepository: repo}

	rows := sqlmock.NewRows([]string{
		"order_id", "user_id", "address", "paid_amount", "shipping_fee",
		"payment_type", "order_status", "created_at", "updated_at", "deleted_at",
	}).AddRow(10, 3, "Bandung", 250000.00, 15000.00, "TRANSFER", "PENDING",
		"2025-10-09 09:00:00", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM orders`).
		WillReturnRows(rows)

	orders, err := srv.GetPendingOrders()
	if err != nil {
		t.Fatalf("GetPendingOrders error: %v", err)
	}
	if len(orders) != 1 || orders[0].OrderId != 10 {
		t.Fatalf("unexpected result: %+v", orders)
	}
}

func TestOrderService_UpdateOrderStatus_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &repository.OrderRepository{DB: db}
	srv := &service.OrderService{OrderRepository: repo}

	mock.ExpectExec(`UPDATE orders SET order_status = \?, updated_at = NOW\(\) WHERE order_id = \? AND deleted_at IS NULL`).
		WithArgs("COMPLETED", 99).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := srv.UpdateOrderStatus(99, "COMPLETED"); err != nil {
		t.Fatalf("UpdateOrderStatus error: %v", err)
	}
}

func TestOrderService_UpdateOrderStatus_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &repository.OrderRepository{DB: db}
	srv := &service.OrderService{OrderRepository: repo}

	// Simulate DB failure
	mock.ExpectExec(`UPDATE orders SET order_status = \?, updated_at = NOW\(\) WHERE order_id = \? AND deleted_at IS NULL`).
		WithArgs("COMPLETED", 77).
		WillReturnError(context.DeadlineExceeded)

	if err := srv.UpdateOrderStatus(77, "COMPLETED"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
