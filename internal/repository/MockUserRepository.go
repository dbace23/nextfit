package repository

import (
	"context"
	"database/sql"
	"nextfit/internal/model"
)

type UserStore interface {
	Create(newUser model.UserModel) (model.UserModel, error)
	FindByEmail(email string) (model.UserModel, error) 
	FindByUserId(UserId int) (model.UserModel, error)
	GetAll() ([]model.UserModel, error)
	Update(userId int, updatedUser model.UserModel) (model.UserModel, error)
	Delete(userId int) error
}