package view

import (
	"fmt"
	"nextfit/internal/controller"
	"nextfit/internal/model"

	"github.com/fatih/color"
)

type View struct {
	Controller *controller.Controller
	User       *model.UserModel
}

func printError(error string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Println(red(error))
}
func printSuccess(success string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green(success))
}

func UnauthorizedMenus() {
	fmt.Println("\n\n\n==== Welcome To NextFit E-commerce Shop ====")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Exit")
}

func AuthorizedAdminMenu(user *model.UserModel) {
	fmt.Printf("\n\n\n==== Welcome %s To Admin Menus ====\n", user.FullName)
	fmt.Println("1. Manage Product")
	fmt.Println("2. Manage Categories")
	fmt.Println("3. Manage Users")
	fmt.Println("4. Manage Orders")
	fmt.Println("5. Reports")
	fmt.Println("6. Logout")
}

func AuthorizedCustomerMenu(user *model.UserModel) {
	fmt.Printf("\n\n\n==== Welcome %s To Customer Menus ====\n", user.FullName)
	fmt.Println("1. Create Orders")
	fmt.Println("2. View Orders")
	fmt.Println("3. Logout")
}

func (view *View) Start() {
	var currentUser *model.UserModel
	loggedIn := false
	for {
		UnauthorizedMenus()

		choice, err := view.Controller.UnathorizedMenuController()
		if err {
			fmt.Println("Error please try again later.")
		}

		switch choice {
		case "1":
			user, err := view.Controller.LoginController()
			if err != nil {
				printError("Email or Password is invalid")

				continue
			}

			currentUser = &user
			loggedIn = true
		case "2":
			_, err := view.Controller.RegisterController()
			if err != nil {
				printError("Error Registering New User")
			} else {
				printSuccess("Success Creating New User, Please Logged In.")
			}
		case "3":
			fmt.Println("Exit")
		}

		if loggedIn {
			if currentUser.IsAdmin {

			} else {

			}
		}
	}

}
