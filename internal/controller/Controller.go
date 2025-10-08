package controller

import (
	"errors"
	"fmt"
	"nextfit/internal/model"
	"nextfit/internal/service"
	"regexp"

	"github.com/manifoldco/promptui"
)

type Controller struct {
	UserService *service.UserService
}

func (controller *Controller) UnathorizedMenuController() (string, bool) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" {
			return errors.New("please choose a valid menu (1/2/3)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		fmt.Println(err)

		return "", false
	}

	return choice, true
}

func (controller *Controller) LoginController() (model.UserModel, error) {
	validateEmail := func(input string) error {
		if input == "" {
			return errors.New("email cannot be empty")
		}
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, input)
		if !matched {
			return errors.New("invalid email format")
		}
		return nil
	}

	validatePassword := func(input string) error {
		if len(input) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		return nil
	}

	emailPrompt := promptui.Prompt{
		Label:    "Email",
		Validate: validateEmail,
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return model.UserModel{}, err

	}

	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Mask:     '*',
		Validate: validatePassword,
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	fmt.Println(password)

	return controller.UserService.Login(email, password)
}

func (controller *Controller) RegisterController() (model.UserModel, error) {
	validateFullName := func(input string) error {
		if len(input) < 3 {
			return errors.New("full name must be at least 3 characters")
		}
		return nil
	}

	validateEmail := func(input string) error {
		if input == "" {
			return errors.New("email cannot be empty")
		}
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, input)
		if !matched {
			return errors.New("invalid email format")
		}
		return nil
	}

	validatePassword := func(input string) error {
		if len(input) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		return nil
	}

	fullNamePrompt := promptui.Prompt{
		Label:    "Full Name",
		Validate: validateFullName,
	}
	fullName, err := fullNamePrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	emailPrompt := promptui.Prompt{
		Label:    "Email",
		Validate: validateEmail,
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Mask:     '*',
		Validate: validatePassword,
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	isAdmin := false

	return controller.UserService.Register(fullName, email, password, isAdmin)
}
