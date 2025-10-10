package service

import (
	"errors"
	"nextfit/internal/model"
	"nextfit/internal/repository"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func (userService *UserService) Login(email, password string) (model.UserModel, error) {
	user, err := userService.UserRepository.FindByEmail(email)
	if err != nil {
		return model.UserModel{}, err
	}

	if user.UserId == 0 {
		return model.UserModel{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return model.UserModel{}, errors.New("invalid password")
	}

	return user, nil
}

func (userService *UserService) Register(fullName, email, password string, isAdmin bool) (model.UserModel, error) {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return model.UserModel{}, errors.New("invalid email format")
	}

	existingUser, err := userService.UserRepository.FindByEmail(email)
	if err != nil {
		return model.UserModel{}, err
	}
	if existingUser.UserId != 0 {
		return model.UserModel{}, errors.New("email already registered")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserModel{}, err
	}

	newUser := model.UserModel{
		FullName:       fullName,
		HashedPassword: string(hashedPwd),
		Email:          email,
		IsAdmin:        isAdmin,
	}

	createdUser, err := userService.UserRepository.Create(newUser)
	if err != nil {
		return model.UserModel{}, err
	}

	return createdUser, nil
}

func (userService *UserService) GetAll() ([]model.UserModel, error) {
	users, err := userService.UserRepository.GetAll()
	if err != nil {
		return nil, errors.New("failed to get users")
	}

	return users, nil
}

func (userService *UserService) Update(userId int, updatedUser model.UserModel) (model.UserModel, error) {
	if userId <= 0 {
		return model.UserModel{}, errors.New("invalid user ID")
	}

	if updatedUser.FullName == "" {
		return model.UserModel{}, errors.New("full name cannot be empty")
	}

	if updatedUser.Email != "" {
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, updatedUser.Email)
		if !matched {
			return model.UserModel{}, errors.New("invalid email format")
		}

		existingUser, err := userService.UserRepository.FindByEmail(updatedUser.Email)
		if err == nil && existingUser.UserId != 0 && existingUser.UserId != userId {
			return model.UserModel{}, errors.New("email already used by another user")
		}
	}

	_, err := userService.UserRepository.FindByUserId(userId)
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	result, err := userService.UserRepository.Update(userId, updatedUser)
	if err != nil {
		return model.UserModel{}, errors.New("failed to update user")
	}

	return result, nil
}

func (userService *UserService) Delete(userId int) error {
	if userId <= 0 {
		return errors.New("invalid user ID")
	}

	_, err := userService.UserRepository.FindByUserId(userId)
	if err != nil {
		return errors.New("user not found")
	}

	err = userService.UserRepository.Delete(userId)
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (userService *UserService) GetById(userId int) (model.UserModel, error) {
	if userId <= 0 {
		return model.UserModel{}, errors.New("invalid user ID")
	}

	user, err := userService.UserRepository.FindByUserId(userId)
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	if user.UserId == 0 {
		return model.UserModel{}, errors.New("user not found")
	}

	return user, nil
}
