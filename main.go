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
	userService := &service.UserService{UserRepository: userRepository}
	controller := &controller.Controller{UserService: userService}
	view := &view.View{Controller: controller}

	view.Start()

}
