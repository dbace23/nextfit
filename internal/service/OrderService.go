package service

import (
	"database/sql"
	"errors"
	"nextfit/internal/model"
	"nextfit/internal/repository"
)

type OrderService struct {
	OrderRepository       *repository.OrderRepository
	OrderDetailRepository *repository.OrderDetailRepository
	ProductRepository     *repository.ProductRepository
	UserRepository        *repository.UserRepository
	DB                    *sql.DB
}

func (orderService *OrderService) CreateOrderWithDetails(newOrder model.OrderModel, orderDetails []model.OrderDetailModel) (model.OrderModel, []model.OrderDetailModel, error) {
	if len(orderDetails) == 0 {
		return model.OrderModel{}, nil, errors.New("order must have at least 1 item")
	}

	if newOrder.UserId <= 0 {
		return model.OrderModel{}, nil, errors.New("invalid user ID")
	}

	if newOrder.Address == "" {
		return model.OrderModel{}, nil, errors.New("address cannot be empty")
	}

	if newOrder.PaidAmount < 0 {
		return model.OrderModel{}, nil, errors.New("paid amount cannot be negative")
	}

	if newOrder.ShippingFee < 0 {
		return model.OrderModel{}, nil, errors.New("shipping fee cannot be negative")
	}

	validPaymentTypes := map[string]bool{
		"CASH":        true,
		"TRANSFER":    true,
		"EWALLET":     true,
		"CREDIT_CARD": true,
	}
	if !validPaymentTypes[newOrder.PaymentType] {
		return model.OrderModel{}, nil, errors.New("invalid payment type")
	}

	validOrderStatuses := map[string]bool{
		"PENDING":   true,
		"COMPLETED": true,
	}
	if !validOrderStatuses[newOrder.OrderStatus] {
		return model.OrderModel{}, nil, errors.New("invalid order status")
	}

	_, err := orderService.UserRepository.FindByUserId(newOrder.UserId)
	if err != nil {
		return model.OrderModel{}, nil, errors.New("user not found")
	}

	for _, orderDetail := range orderDetails {
		if orderDetail.ProductId <= 0 {
			return model.OrderModel{}, nil, errors.New("invalid product ID in order detail")
		}

		if orderDetail.Quantity <= 0 {
			return model.OrderModel{}, nil, errors.New("quantity must be greater than 0")
		}

		_, err := orderService.ProductRepository.FindByProductId(orderDetail.ProductId)
		if err != nil {
			return model.OrderModel{}, nil, errors.New("product not found in order detail")
		}
	}

	tx, err := orderService.DB.Begin()
	if err != nil {
		return model.OrderModel{}, nil, errors.New("failed to start transaction")
	}

	createdOrder, err := orderService.OrderRepository.CreateOrder(tx, newOrder)
	if err != nil {
		tx.Rollback()
		return model.OrderModel{}, nil, errors.New("failed to create order")
	}

	for i := range orderDetails {
		orderDetails[i].OrderId = createdOrder.OrderId
	}

	createdOrderDetails, err := orderService.OrderDetailRepository.CreateMultipleOrderDetails(tx, orderDetails)
	if err != nil {
		tx.Rollback()
		return model.OrderModel{}, nil, errors.New("failed to create order details")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return model.OrderModel{}, nil, errors.New("failed to commit transaction")
	}

	return createdOrder, createdOrderDetails, nil
}

func (orderService *OrderService) GetOrdersByUserId(userId int) ([]model.OrderModel, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	orders, err := orderService.OrderRepository.GetOrdersByUserId(userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderService *OrderService) GetOrderDetailsWithProducts(orderId int, userId int) (model.OrderModel, []map[string]interface{}, error) {
	if orderId <= 0 {
		return model.OrderModel{}, nil, errors.New("invalid order ID")
	}

	if userId <= 0 {
		return model.OrderModel{}, nil, errors.New("invalid user ID")
	}

	order, err := orderService.OrderRepository.GetOrderById(orderId)
	if err != nil {
		return model.OrderModel{}, nil, err
	}

	if order.UserId != userId {
		return model.OrderModel{}, nil, errors.New("order not found or access denied")
	}

	orderDetails, err := orderService.OrderDetailRepository.GetOrderDetailsByOrderId(orderId)
	if err != nil {
		return model.OrderModel{}, nil, err
	}

	var detailsWithProducts []map[string]interface{}
	for _, detail := range orderDetails {
		product, err := orderService.ProductRepository.FindByProductId(detail.ProductId)
		if err != nil {
			continue
		}

		detailWithProduct := map[string]interface{}{
			"order_detail_id": detail.OrderDetailId,
			"product_id":      detail.ProductId,
			"product_name":    product.ProductName,
			"selling_price":   product.SellingPrice,
			"quantity":        detail.Quantity,
			"subtotal":        float64(detail.Quantity) * product.SellingPrice,
		}
		detailsWithProducts = append(detailsWithProducts, detailWithProduct)
	}

	return order, detailsWithProducts, nil
}
