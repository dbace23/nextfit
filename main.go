package main

import (
	"nextfit/internal/controller"
	"nextfit/internal/helper"
	"nextfit/internal/repository"
	"nextfit/internal/service"
	"nextfit/internal/view"
)

func main() {

	db := helper.GetMysqlConnection()

	userRepository := &repository.UserRepository{DB: db}
	categoryRepository := &repository.CategoryRepository{DB: db}
	productRepository := &repository.ProductRepository{DB: db}
	orderRepository := &repository.OrderRepository{DB: db}
	orderDetailRepository := &repository.OrderDetailRepository{DB: db}

	categoryService := &service.CategoryService{CategoryRepository: categoryRepository}
	userService := &service.UserService{UserRepository: userRepository}
	productService := &service.ProductService{
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
	}
	orderService := &service.OrderService{
		OrderRepository:       orderRepository,
		OrderDetailRepository: orderDetailRepository,
		ProductRepository:     productRepository,
		UserRepository:        userRepository,
		DB:                    db,
	}

	controller := &controller.Controller{
		UserService:     userService,
		CategoryService: categoryService,
		ProductService:  productService,
		OrderService:    orderService,
	}

	view := &view.View{Controller: controller}

	view.Start()

}
