package service_test

import (
	"errors"
	"testing"

	"nextfit/internal/model"
	"nextfit/internal/service"
)

type UserStore struct {
	CreateFn func() (model.UserModel, error)
	FindByEmailFn func() (model.UserModel, error) 
	FindByUserIdFn func() (model.UserModel, error)
	GetAllFn func() ([]model.UserModel, error)
	UpdateFn func() (model.UserModel, error)
	DeleteFn func error
}