package view

import (
	"fmt"
	"nextfit/internal/controller"
	"nextfit/internal/model"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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

func unauthorizedMenus() {
	fmt.Println("\n\n\n==== Welcome To NextFit E-commerce Shop ====")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Exit")
}

func authorizedAdminMenu(user *model.UserModel) {
	fmt.Printf("\n\n\n==== Welcome %s To Admin Menus ====\n", user.FullName)
	fmt.Println("1. Manage Product")
	fmt.Println("2. Manage Categories")
	fmt.Println("3. Manage Users")
	fmt.Println("4. Manage Orders")
	fmt.Println("5. Reports")
	fmt.Println("6. Logout")
}

func authorizedCustomerMenu(user *model.UserModel) {
	fmt.Printf("\n\n\n==== Welcome %s To Customer Menus ====\n", user.FullName)
	fmt.Println("1. Create Orders")
	fmt.Println("2. View Orders")
	fmt.Println("3. View Order Details")
	fmt.Println("4. Logout")
}

func showManageCategoryMenu() {
	fmt.Println("1. Create Category")
	fmt.Println("2. Update Category")
	fmt.Println("3. Delete Category")
	fmt.Println("4. Exit")
}

func showCategories(categories []model.CategoryModel) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Name", "Created At", "Updated At", "Deleted At"})

	for _, cat := range categories {
		updatedAtStr := "-"
		if cat.UpdatedAt.Valid {
			updatedAtStr = cat.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}

		deletedAtStr := "-"
		if cat.DeletedAt.Valid {
			deletedAtStr = cat.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			strconv.Itoa(cat.CategoryId),
			cat.Name,
			cat.CreatedAt.Format("2006-01-02 15:04:05"),
			updatedAtStr,
			deletedAtStr,
		}
		table.Append(row)
	}

	table.Render()
}

func showManageProductMenu() {
	fmt.Println("1. Create Product")
	fmt.Println("2. Update Product")
	fmt.Println("3. Delete Product")
	fmt.Println("4. Exit")
}

func showProducts(products []model.ProductModel) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Name", "Category ID", "Price", "Created At", "Updated At", "Deleted At"})

	for _, product := range products {
		priceStr := fmt.Sprintf("%.2f", product.SellingPrice)

		updatedAtStr := "-"
		if product.UpdatedAt.Valid {
			updatedAtStr = product.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}

		deletedAtStr := "-"
		if product.DeletedAt.Valid {
			deletedAtStr = product.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			strconv.Itoa(product.ProductId),
			product.ProductName,
			strconv.Itoa(product.CategoryId),
			priceStr,
			product.CreatedAt.Format("2006-01-02 15:04:05"),
			updatedAtStr,
			deletedAtStr,
		}
		table.Append(row)
	}

	table.Render()
}

func (view *View) manageProductsMenu() {
	for {
		products, err := view.Controller.GetAllProductsController()
		if err != nil {
			printError("Cannot get products")
			return
		}

		fmt.Println("\n==== Current Products ====")
		showProducts(products)

		fmt.Println("\n==== Manage Products ====")
		showManageProductMenu()

		choice, err := view.Controller.ManageProductsChoiceController()
		if err != nil {
			printError("Error please try again later.")
			continue
		}

		switch choice {
		case "1":
			product, err := view.Controller.CreateProductController()
			if err != nil {
				printError(fmt.Sprintf("Failed to create product: %v", err))
			} else {
				printSuccess(fmt.Sprintf("Product '%s' created successfully with ID: %d", product.ProductName, product.ProductId))
			}

		case "2":
			product, err := view.Controller.UpdateProductController()
			if err != nil {
				printError(fmt.Sprintf("Failed to update product: %v", err))
			} else {
				printSuccess(fmt.Sprintf("Product updated successfully to '%s'", product.ProductName))
			}

		case "3":
			err := view.Controller.DeleteProductController()
			if err != nil {
				printError(fmt.Sprintf("Failed to delete product: %v", err))
			} else {
				printSuccess("Product deleted successfully")
			}

		case "4":
			return
		}

		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}

func showManageUserMenu() {
	fmt.Println("1. Create User")
	fmt.Println("2. Update User")
	fmt.Println("3. Delete User")
	fmt.Println("4. Exit")
}

func showUsers(users []model.UserModel) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Full Name", "Email", "Admin", "Created At", "Updated At", "Deleted At"})

	for _, user := range users {
		adminStr := "No"
		if user.IsAdmin {
			adminStr = "Yes"
		}

		updatedAtStr := "-"
		if user.UpdatedAt.Valid {
			updatedAtStr = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}

		deletedAtStr := "-"
		if user.DeletedAt.Valid {
			deletedAtStr = user.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}

		row := []string{
			strconv.Itoa(user.UserId),
			user.FullName,
			user.Email,
			adminStr,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			updatedAtStr,
			deletedAtStr,
		}
		table.Append(row)
	}

	table.Render()
}

func (view *View) manageUsersMenu() {
	for {
		users, err := view.Controller.GetAllUsersController()
		if err != nil {
			printError("Cannot get users")
			return
		}

		fmt.Println("\n==== Current Users ====")
		showUsers(users)

		fmt.Println("\n==== Manage Users ====")
		showManageUserMenu()

		choice, err := view.Controller.ManageUsersChoiceController()
		if err != nil {
			printError("Error please try again later.")
			continue
		}

		switch choice {
		case "1":
			user, err := view.Controller.CreateUserController()
			if err != nil {
				printError(fmt.Sprintf("Failed to create user: %v", err))
			} else {
				printSuccess(fmt.Sprintf("User '%s' created successfully with ID: %d", user.FullName, user.UserId))
			}

		case "2":
			user, err := view.Controller.UpdateUserController()
			if err != nil {
				printError(fmt.Sprintf("Failed to update user: %v", err))
			} else {
				printSuccess(fmt.Sprintf("User updated successfully to '%s'", user.FullName))
			}

		case "3":
			err := view.Controller.DeleteUserController()
			if err != nil {
				printError(fmt.Sprintf("Failed to delete user: %v", err))
			} else {
				printSuccess("User deleted successfully")
			}

		case "4":
			return
		}

		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}

func (view *View) manageCategoriesMenu() {
	for {
		categories, err := view.Controller.GetAllCategoriesController()
		if err != nil {
			printError("Cannot get categories")
			return
		}

		fmt.Println("\n==== Current Categories ====")
		showCategories(categories)

		fmt.Println("\n==== Manage Categories ====")
		showManageCategoryMenu()

		choice, err := view.Controller.ManageCategoriesChoiceController()
		if err != nil {
			printError("Error please try again later.")
			continue
		}

		switch choice {
		case "1":
			category, err := view.Controller.CreateCategoryController()
			if err != nil {
				printError(fmt.Sprintf("Failed to create category: %v", err))
			} else {
				printSuccess(fmt.Sprintf("Category '%s' created successfully with ID: %d", category.Name, category.CategoryId))
			}

		case "2":
			category, err := view.Controller.UpdateCategoryController()
			if err != nil {
				printError(fmt.Sprintf("Failed to update category: %v", err))
			} else {
				printSuccess(fmt.Sprintf("Category updated successfully to '%s'", category.Name))
			}

		case "3":
			err := view.Controller.DeleteCategoryController()
			if err != nil {
				printError(fmt.Sprintf("Failed to delete category: %v", err))
			} else {
				printSuccess("Category deleted successfully")
			}

		case "4":
			return
		}

		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}


func (view *View) loggedInAdmin(currentUser *model.UserModel) bool {
	authorizedAdminMenu(currentUser)

	choice, err := view.Controller.AuthorizedAdminMenuController()
	if err != nil {
		fmt.Println("Error please try again later.")
		return true
	}

	switch choice {
	case "1":
		view.manageProductsMenu()
	case "2":
		view.manageCategoriesMenu()
	case "3":
		view.manageUsersMenu()
	case "4":
		// Handle manage orders
	case "5":
		view.showReport()
	case "6":
		return false // Logout
	}

	return true
}

func (view *View) loggedInCustomer(currentUser *model.UserModel) bool {
	authorizedCustomerMenu(currentUser)

	choice, err := view.Controller.AuthorizedCustomerMenuController()
	if err != nil {
		fmt.Println("Error please try again later.")
		return true
	}

	switch choice {
	case "1":
		order, orderDetails, err := view.Controller.CreateOrderController(currentUser.UserId)
		if err != nil {
			fmt.Printf("Error creating order: %v\n", err)
		} else {
			fmt.Printf("Order created successfully! Order ID: %d\n", order.OrderId)
			fmt.Printf("Total items: %d\n", len(orderDetails))
			fmt.Printf("Total amount: $%.2f\n", order.PaidAmount)
		}
	case "2":
		orders, err := view.Controller.ViewOrdersController(currentUser.UserId)
		if err != nil {
			fmt.Printf("Error retrieving orders: %v\n", err)
		} else if len(orders) == 0 {
			fmt.Println("You have no orders yet.")
		} else {
			fmt.Printf("\n==== Your Orders ====\n")
			showOrders(orders)
		}
	case "3":
		orderId, err := view.Controller.GetOrderIdInputController()
		if err != nil {
			fmt.Printf("Error getting order ID: %v\n", err)
		} else {
			order, orderDetails, err := view.Controller.ViewOrderDetailsController(orderId, currentUser.UserId)
			if err != nil {
				fmt.Printf("Error retrieving order details: %v\n", err)
			} else {
				showOrderDetails(order, orderDetails)
			}
		}
	case "4":
		return false
	}

	return true
}

func (view *View) loggedIn(currentUser *model.UserModel) {
	for {
		var shouldContinue bool

		if currentUser.IsAdmin {
			shouldContinue = view.loggedInAdmin(currentUser)
		} else {
			shouldContinue = view.loggedInCustomer(currentUser)
		}

		if !shouldContinue {
			break
		}
	}
}

func (view *View) Start() {
	for {
		var currentUser *model.UserModel
		loggedIn := false

		unauthorizedMenus()

		choice, err := view.Controller.UnathorizedMenuController()
		if err != nil {
			fmt.Println("Error please try again later.")
			continue
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
			fmt.Println("Thank you for using NextFit E-commerce!")
			return
		}

		if loggedIn {
			view.loggedIn(currentUser)
			printSuccess("You have been logged out successfully.")
		}
	}
}

//report view
func (view *View) showReport() {
	for {
		fmt.Println("\n==== Reports Menu ====")
		fmt.Println("1. Sales Summary (last 30 days)")
		fmt.Println("2. Top 10 Products (last 30 days)")
		fmt.Println("3. Orders by Status (last 30 days)")
		fmt.Println("4. Back")

		choice, err := view.Controller.ReportMenuChoiceController()
		if err != nil {
			fmt.Println("Error: invalid choice.")
			continue
		}

		switch choice {
		case "1":
			reports, err := view.Controller.GetSalesSummaryDailyController()
			if err != nil {
				fmt.Println("Error getting sales summary:", err)
				continue
			}
			showSalesSummaryTable(reports)

		case "2":
			products, err := view.Controller.GetTopProductsController()
			if err != nil {
				fmt.Println("Error getting top products:", err)
				continue
			}
			showTopProductsTable(products)

		case "3":
			orders, err := view.Controller.GetOrdersByStatusController()
			if err != nil {
				fmt.Println("Error getting orders by status:", err)
				continue
			}
			showOrdersByStatusTable(orders)

		case "4":
			return
		}

		fmt.Println("\nPress Enter to return to the report menu...")
		fmt.Scanln()
	}
}
//////render view

func showSalesSummaryTable(data []model.DailySalesReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Date", "Revenue", "Orders", "Customers", "AOV"})
	for _, r := range data {
		table.Append([]string{
			r.Day.Format("2006-01-02"),
			fmt.Sprintf("%.2f", r.Revenue),
			fmt.Sprintf("%d", r.Orders),
			fmt.Sprintf("%d", r.UniqueCustomers),
			fmt.Sprintf("%.2f", r.AOV),
		})
	}
	table.Render()
}

func showTopProductsTable(data []model.TopProductReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Product ID", "Name", "Quantity", "GMV"})
	for _, p := range data {
		table.Append([]string{
			fmt.Sprintf("%d", p.ProductId),
			p.ProductName,
			fmt.Sprintf("%d", p.Quantity),
			fmt.Sprintf("%.2f", p.GMV),
		})
	}
	table.Render()
}

func showOrdersByStatusTable(data []model.OrdersByStatusReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Status", "Orders", "Revenue"})
	for _, o := range data {
		table.Append([]string{
			o.Status,
			fmt.Sprintf("%d", o.Count),
			fmt.Sprintf("%.2f", o.Revenue),
		})
	}
	table.Render()
}
///////////////////////////////////////////////////////////////////////////////////


func showOrders(orders []model.OrderModel) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Order ID", "Address", "Paid Amount", "Shipping Fee", "Payment Type", "Order Status", "Created At"})

	for _, order := range orders {
		row := []string{
			strconv.Itoa(order.OrderId),
			order.Address,
			fmt.Sprintf("$%.2f", order.PaidAmount),
			fmt.Sprintf("$%.2f", order.ShippingFee),
			order.PaymentType,
			order.OrderStatus,
			order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		table.Append(row)
	}

	table.Render()
}

func showOrderDetails(order model.OrderModel, orderDetails []map[string]interface{}) {
	fmt.Printf("\n==== Order Details ====\n")
	fmt.Printf("Order ID: %d\n", order.OrderId)
	fmt.Printf("Address: %s\n", order.Address)
	fmt.Printf("Payment Type: %s\n", order.PaymentType)
	fmt.Printf("Order Status: %s\n", order.OrderStatus)
	fmt.Printf("Shipping Fee: $%.2f\n", order.ShippingFee)
	fmt.Printf("Order Date: %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("\n==== Items ====\n")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Product ID", "Product Name", "Price", "Quantity", "Subtotal"})

	var grandTotal float64 = 0
	for _, detail := range orderDetails {
		subtotal := detail["subtotal"].(float64)
		grandTotal += subtotal

		row := []string{
			fmt.Sprintf("%d", detail["product_id"].(int)),
			detail["product_name"].(string),
			fmt.Sprintf("$%.2f", detail["selling_price"].(float64)),
			fmt.Sprintf("%d", detail["quantity"].(int)),
			fmt.Sprintf("$%.2f", subtotal),
		}
		table.Append(row)
	}

	table.Render()
	fmt.Printf("\nSubtotal: $%.2f\n", grandTotal)
	fmt.Printf("Shipping Fee: $%.2f\n", order.ShippingFee)
	fmt.Printf("Total Paid: $%.2f\n", order.PaidAmount)
}
