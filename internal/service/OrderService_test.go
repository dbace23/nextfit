package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"nextfit/internal/model"
	"nextfit/internal/service"
)

type mockOrderRepo struct {
	createOrderFn       func(tx *sql.Tx, m model.OrderModel) (model.OrderModel, error)
	getOrdersByUserIdFn func(int) ([]model.OrderModel, error)
	getOrderByIdFn      func(int) (model.OrderModel, error)
	getPendingFn        func(context.Context) ([]model.OrderModel, error)
	updateStatusFn      func(context.Context, int, string) error
}

func (m *mockOrderRepo) CreateOrder(tx *sql.Tx, ord model.OrderModel) (model.OrderModel, error) {
	if m.createOrderFn != nil { return m.createOrderFn(tx, ord) }
	return model.OrderModel{}, nil
}
func (m *mockOrderRepo) GetOrdersByUserId(uid int) ([]model.OrderModel, error) {
	if m.getOrdersByUserIdFn != nil { return m.getOrdersByUserIdFn(uid) }
	return nil, nil
}
func (m *mockOrderRepo) GetOrderById(id int) (model.OrderModel, error) {
	if m.getOrderByIdFn != nil { return m.getOrderByIdFn(id) }
	return model.OrderModel{}, nil
}
func (m *mockOrderRepo) GetPendingOrders(ctx context.Context) ([]model.OrderModel, error) {
	if m.getPendingFn != nil { return m.getPendingFn(ctx) }
	return nil, nil
}
func (m *mockOrderRepo) UpdateOrderStatus(ctx context.Context, id int, s string) error {
	if m.updateStatusFn != nil { return m.updateStatusFn(ctx, id, s) }
	return nil
}

type mockOrderDetailRepo struct {
	createMultipleFn       func(tx *sql.Tx, ds []model.OrderDetailModel) ([]model.OrderDetailModel, error)
	getDetailsByOrderIdFn  func(int) ([]model.OrderDetailModel, error)
}

func (m *mockOrderDetailRepo) CreateMultipleOrderDetails(tx *sql.Tx, ds []model.OrderDetailModel) ([]model.OrderDetailModel, error) {
	if m.createMultipleFn != nil { return m.createMultipleFn(tx, ds) }
	return ds, nil
}
func (m *mockOrderDetailRepo) GetOrderDetailsByOrderId(orderId int) ([]model.OrderDetailModel, error) {
	if m.getDetailsByOrderIdFn != nil { return m.getDetailsByOrderIdFn(orderId) }
	return nil, nil
}

type mockProductRepo struct {
	findByIdFn func(int) (model.ProductModel, error)
}
func (m *mockProductRepo) FindByProductId(id int) (model.ProductModel, error) {
	if m.findByIdFn != nil { return m.findByIdFn(id) }
	return model.ProductModel{}, errors.New("not found")
}

type mockUserRepo struct {
	findByIdFn func(int) (model.UserModel, error)
}
func (m *mockUserRepo) FindByUserId(id int) (model.UserModel, error) {
	if m.findByIdFn != nil { return m.findByIdFn(id) }
	return model.UserModel{}, errors.New("not found")
}

// ---------- Helpers, fake sql ----------

func newSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	return db, mock
}

// ---------- Tests: CreateOrderWithDetails ----------

func TestCreateOrderWithDetails_Success(t *testing.T) {
	db, mock := newSQLMock(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	orderRepo := &mockOrderRepo{
		createOrderFn: func(tx *sql.Tx, m model.OrderModel) (model.OrderModel, error) {
			m.OrderId = 101
			return m, nil
		},
	}
	orderDetailRepo := &mockOrderDetailRepo{
		createMultipleFn: func(tx *sql.Tx, ds []model.OrderDetailModel) ([]model.OrderDetailModel, error) {
			for i := range ds {
				ds[i].OrderDetailId = 1000 + i
			}
			return ds, nil
		},
	}
	productRepo := &mockProductRepo{
		findByIdFn: func(id int) (model.ProductModel, error) {
			return model.ProductModel{ProductId: id, ProductName: "Item", SellingPrice: 25000}, nil
		},
	}
	userRepo := &mockUserRepo{
		findByIdFn: func(id int) (model.UserModel, error) {
			return model.UserModel{UserId: id, FullName: "Alice"}, nil
		},
	}

	svc := &service.OrderService{
		OrderRepository:       orderRepo,
		OrderDetailRepository: orderDetailRepo,
		ProductRepository:     productRepo,
		UserRepository:        userRepo,
		DB:                    db,
	}

	order := model.OrderModel{
		UserId:      7,
		Address:     "Jl. Mawar",
		PaidAmount:  100000,
		ShippingFee: 15000,
		PaymentType: "EWALLET",
		OrderStatus: "PENDING",
		CreatedAt:   time.Now(),
	}
	details := []model.OrderDetailModel{
		{ProductId: 1, Quantity: 2},
		{ProductId: 2, Quantity: 1},
	}

	createdOrder, createdDetails, err := svc.CreateOrderWithDetails(order, details)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if createdOrder.OrderId != 101 {
		t.Fatalf("expected order id 101, got %#v", createdOrder)
	}
	if len(createdDetails) != 2 || createdDetails[0].OrderDetailId == 0 {
		t.Fatalf("expected 2 details with ids, got %#v", createdDetails)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("tx expectations not met: %v", err)
	}
}

func TestCreateOrderWithDetails_ValidationErrors(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	svc := &service.OrderService{
		OrderRepository:       &mockOrderRepo{},
		OrderDetailRepository: &mockOrderDetailRepo{},
		ProductRepository:     &mockProductRepo{},
		UserRepository:        &mockUserRepo{},
		DB:                    db,
	}

	// No items
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{}, nil); err == nil {
		t.Fatal("expected error for empty details")
	}
	// Invalid user
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 0}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for invalid user")
	}
	// Empty address
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 1, Address: ""}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for empty address")
	}
	// Negative amounts
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 1, Address: "X", PaidAmount: -1}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for negative paid amount")
	}
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, ShippingFee: -1}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for negative shipping fee")
	}
	// Invalid payment type
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, ShippingFee: 0, PaymentType: "BITCOIN"}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for invalid payment type")
	}
	// Invalid order status
	if _, _, err := svc.CreateOrderWithDetails(model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, ShippingFee: 0, PaymentType: "EWALLET", OrderStatus: "SHIPPED"}, []model.OrderDetailModel{{ProductId: 1, Quantity: 1}}); err == nil {
		t.Fatal("expected error for invalid order status")
	}
}

func TestCreateOrderWithDetails_UserNotFound(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	userRepo := &mockUserRepo{
		findByIdFn: func(id int) (model.UserModel, error) {
			return model.UserModel{}, errors.New("not found")
		},
	}
	productRepo := &mockProductRepo{
		findByIdFn: func(id int) (model.ProductModel, error) {
			return model.ProductModel{ProductId: id}, nil
		},
	}

	svc := &service.OrderService{
		OrderRepository:       &mockOrderRepo{},
		OrderDetailRepository: &mockOrderDetailRepo{},
		ProductRepository:     productRepo,
		UserRepository:        userRepo,
		DB:                    db,
	}

	_, _, err := svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 9, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 1, Quantity: 1}},
	)
	if err == nil {
		t.Fatal("expected user not found error")
	}
}

func TestCreateOrderWithDetails_ProductValidationErrors(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	userRepo := &mockUserRepo{
		findByIdFn: func(id int) (model.UserModel, error) { return model.UserModel{UserId: id}, nil },
	}
	productRepo := &mockProductRepo{
		findByIdFn: func(id int) (model.ProductModel, error) { return model.ProductModel{ProductId: id}, nil },
	}

	svc := &service.OrderService{
		OrderRepository:       &mockOrderRepo{},
		OrderDetailRepository: &mockOrderDetailRepo{},
		ProductRepository:     productRepo,
		UserRepository:        userRepo,
		DB:                    db,
	}

	// invalid product id in detail
	_, _, err := svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 0, Quantity: 1}},
	)
	if err == nil {
		t.Fatal("expected invalid product ID error")
	}

	// invalid quantity
	_, _, err = svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 1, Quantity: 0}},
	)
	if err == nil {
		t.Fatal("expected invalid quantity error")
	}

	// product not found
	productRepo.findByIdFn = func(id int) (model.ProductModel, error) { return model.ProductModel{}, errors.New("not found") }
	_, _, err = svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 1, Quantity: 1}},
	)
	if err == nil {
		t.Fatal("expected product not found error")
	}
}

func TestCreateOrderWithDetails_RepoFailuresRollback(t *testing.T) {
	db, mock := newSQLMock(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectRollback()

	userRepo := &mockUserRepo{
		findByIdFn: func(id int) (model.UserModel, error) { return model.UserModel{UserId: id}, nil },
	}
	productRepo := &mockProductRepo{
		findByIdFn: func(id int) (model.ProductModel, error) { return model.ProductModel{ProductId: id}, nil },
	}

	orderRepo := &mockOrderRepo{
		createOrderFn: func(tx *sql.Tx, m model.OrderModel) (model.OrderModel, error) {
			return model.OrderModel{}, errors.New("create failed")
		},
	}
	orderDetailRepo := &mockOrderDetailRepo{
		createMultipleFn: func(tx *sql.Tx, ds []model.OrderDetailModel) ([]model.OrderDetailModel, error) {
			return nil, errors.New("bulk create failed")
		},
	}

	svc := &service.OrderService{
		OrderRepository:       orderRepo,
		OrderDetailRepository: orderDetailRepo,
		ProductRepository:     productRepo,
		UserRepository:        userRepo,
		DB:                    db,
	}

	// fail on CreateOrder
	_, _, err := svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 1, Quantity: 1}},
	)
	if err == nil {
		t.Fatal("expected create order failure")
	}

	// succeed CreateOrder but fail CreateMultiple
	orderRepo.createOrderFn = func(tx *sql.Tx, m model.OrderModel) (model.OrderModel, error) {
		m.OrderId = 42
		return m, nil
	}
	_, _, err = svc.CreateOrderWithDetails(
		model.OrderModel{UserId: 1, Address: "X", PaidAmount: 1, PaymentType: "CASH", OrderStatus: "PENDING"},
		[]model.OrderDetailModel{{ProductId: 1, Quantity: 1}},
	)
	if err == nil {
		t.Fatal("expected bulk create failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("tx expectations not met: %v", err)
	}
}

// ---------- Tests: GetOrdersByUserId ----------

func TestGetOrdersByUserId_Success(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	orderRepo := &mockOrderRepo{
		getOrdersByUserIdFn: func(id int) ([]model.OrderModel, error) {
			return []model.OrderModel{
				{OrderId: 1, UserId: id},
				{OrderId: 2, UserId: id},
			}, nil
		},
	}
	svc := &service.OrderService{
		OrderRepository: orderRepo,
		DB:              db,
	}

	got, err := svc.GetOrdersByUserId(9)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 || got[0].UserId != 9 {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestGetOrdersByUserId_Invalid(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	svc := &service.OrderService{OrderRepository: &mockOrderRepo{}, DB: db}
	if _, err := svc.GetOrdersByUserId(0); err == nil {
		t.Fatal("expected invalid user id error")
	}
}

// ---------- Tests: GetOrderDetailsWithProducts ----------

func TestGetOrderDetailsWithProducts_Success(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	orderRepo := &mockOrderRepo{
		getOrderByIdFn: func(id int) (model.OrderModel, error) {
			return model.OrderModel{OrderId: id, UserId: 5}, nil
		},
	}
	detailRepo := &mockOrderDetailRepo{
		getDetailsByOrderIdFn: func(id int) ([]model.OrderDetailModel, error) {
			return []model.OrderDetailModel{
				{OrderDetailId: 10, OrderId: id, ProductId: 1, Quantity: 2},
				{OrderDetailId: 11, OrderId: id, ProductId: 2, Quantity: 1},
			}, nil
		},
	}
	productRepo := &mockProductRepo{
		findByIdFn: func(id int) (model.ProductModel, error) {
			if id == 2 { // simulate one missing product; service will skip it
				return model.ProductModel{}, errors.New("not found")
			}
			return model.ProductModel{ProductId: id, ProductName: "P", SellingPrice: 1000}, nil
		},
	}

	svc := &service.OrderService{
		OrderRepository:       orderRepo,
		OrderDetailRepository: detailRepo,
		ProductRepository:     productRepo,
		DB:                    db,
	}

	order, details, err := svc.GetOrderDetailsWithProducts(77, 5)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if order.OrderId != 77 {
		t.Fatalf("bad order: %#v", order)
	}
	if len(details) != 1 { // one skipped due to missing product
		t.Fatalf("expected 1 detail, got %#v", details)
	}
}

func TestGetOrderDetailsWithProducts_InvalidInput(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	svc := &service.OrderService{OrderRepository: &mockOrderRepo{}, OrderDetailRepository: &mockOrderDetailRepo{}, ProductRepository: &mockProductRepo{}, DB: db}

	if _, _, err := svc.GetOrderDetailsWithProducts(0, 1); err == nil {
		t.Fatal("expected invalid order id error")
	}
	if _, _, err := svc.GetOrderDetailsWithProducts(1, 0); err == nil {
		t.Fatal("expected invalid user id error")
	}
}

func TestGetOrderDetailsWithProducts_AccessDenied(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	orderRepo := &mockOrderRepo{
		getOrderByIdFn: func(id int) (model.OrderModel, error) {
			return model.OrderModel{OrderId: id, UserId: 999}, nil // does not match
		},
	}
	svc := &service.OrderService{OrderRepository: orderRepo, OrderDetailRepository: &mockOrderDetailRepo{}, ProductRepository: &mockProductRepo{}, DB: db}
	if _, _, err := svc.GetOrderDetailsWithProducts(7, 5); err == nil {
		t.Fatal("expected access denied error")
	}
}

// ---------- Tests: GetPendingOrders / UpdateOrderStatus ----------

func TestGetPendingOrders_Success(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	orderRepo := &mockOrderRepo{
		getPendingFn: func(ctx context.Context) ([]model.OrderModel, error) {
			return []model.OrderModel{{OrderId: 1, OrderStatus: "PENDING"}}, nil
		},
	}
	svc := &service.OrderService{OrderRepository: orderRepo, DB: db}

	got, err := svc.GetPendingOrders()
	if err != nil || len(got) != 1 || got[0].OrderStatus != "PENDING" {
		t.Fatalf("unexpected: got=%#v err=%v", got, err)
	}
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	db, _ := newSQLMock(t)
	defer db.Close()

	orderRepo := &mockOrderRepo{
		updateStatusFn: func(ctx context.Context, id int, status string) error {
			if id != 555 || status != "SHIPPED" {
				return errors.New("bad args")
			}
			return nil
		},
	}
	svc := &service.OrderService{OrderRepository: orderRepo, DB: db}

	if err := svc.UpdateOrderStatus(555, "SHIPPED"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
