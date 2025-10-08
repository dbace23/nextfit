package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nextfit/internal/helper"
	"nextfit/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func (userRepository *UserRepository) Create(newUser model.UserModel) (model.UserModel, error) {
	query := `
    INSERT INTO users (
        full_name, password_hash, email, isAdmin
    ) VALUES (?, ?, ?, ?)
    `

	result, err := userRepository.DB.ExecContext(context.Background(), query,
		newUser.FullName,
		newUser.HashedPassword,
		newUser.Email,
		helper.BoolToTinyInt(newUser.IsAdmin),
	)

	if err != nil {
		fmt.Println(err)
		return model.UserModel{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return model.UserModel{}, err
	}

	newUser.UserId = int(lastInsertId)
	return newUser, nil
}

func (userRepository *UserRepository) FindByEmail(email string) (model.UserModel, error) {
	query := `	
	SELECT
	user_id, full_name, password_hash, email, isAdmin, created_at, updated_at, deleted_at
	FROM users
	WHERE email = ?
	`

	result, err := userRepository.DB.QueryContext(context.Background(), query, email)
	if err != nil {
		fmt.Println(err)
		return model.UserModel{}, err
	}

	var user model.UserModel
	if result.Next() {
		err := result.Scan(&user.UserId, &user.FullName, &user.HashedPassword, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return model.UserModel{}, err
		}

	}

	return user, nil
}

func (userRepository *UserRepository) FindByUserId(UserId int) (model.UserModel, error) {
	query := `	
	SELECT
	user_id, full_name, email, isAdmin, created_at, updated_at, deleted_at
	FROM users
	WHERE user_id = ? 
	`

	result, err := userRepository.DB.QueryContext(context.Background(), query, UserId)
	if err != nil {
		return model.UserModel{}, err
	}

	var user model.UserModel
	if result.Next() {
		err := result.Scan(&user.UserId, &user.FullName, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return model.UserModel{}, err
		}
	}

	return user, nil
}
