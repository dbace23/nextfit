package controller

import (
	"errors"
	"fmt"
	"nextfit/internal/model"
	"nextfit/internal/service"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
)

type Controller struct {
	UserService     *service.UserService
	CategoryService *service.CategoryService
	ProductService  *service.ProductService
	OrderService    *service.OrderService
}

func (controller *Controller) UnathorizedMenuController() (string, error) {
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
		return "", err
	}

	return choice, nil
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

func (controller *Controller) AuthorizedAdminMenuController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" && input != "5" && input != "6" {
			return errors.New("please choose a valid menu (1/2/3/4/5/6)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3/4/5/6): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return choice, nil
}

func (controller *Controller) AuthorizedCustomerMenuController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" {
			return errors.New("please choose a valid menu (1/2/3/4)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3/4): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return choice, nil
}

func (controller *Controller) GetAllCategoriesController() ([]model.CategoryModel, error) {
	categories, err := controller.CategoryService.GetAll()
	if err != nil {
		return []model.CategoryModel{}, err
	}

	return categories, nil
}

func (controller *Controller) ManageCategoriesChoiceController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" {
			return errors.New("please choose a valid menu (1/2/3/4)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3/4): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return choice, nil
}

func (controller *Controller) CreateCategoryController() (model.CategoryModel, error) {
	validateName := func(input string) error {
		if len(input) < 2 {
			return errors.New("category name must be at least 2 characters")
		}
		return nil
	}

	namePrompt := promptui.Prompt{
		Label:    "Category Name",
		Validate: validateName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return model.CategoryModel{}, err
	}

	return controller.CategoryService.Create(name)
}

func (controller *Controller) UpdateCategoryController() (model.CategoryModel, error) {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("category ID cannot be empty")
		}
		return nil
	}

	validateName := func(input string) error {
		if len(input) < 2 {
			return errors.New("category name must be at least 2 characters")
		}
		return nil
	}

	idPrompt := promptui.Prompt{
		Label:    "Category ID to update",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return model.CategoryModel{}, err
	}

	categoryId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &categoryId); err != nil {
		return model.CategoryModel{}, errors.New("invalid category ID format")
	}

	namePrompt := promptui.Prompt{
		Label:    "New Category Name",
		Validate: validateName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return model.CategoryModel{}, err
	}

	updatedCategory := model.CategoryModel{
		Name: name,
	}

	return controller.CategoryService.Update(categoryId, updatedCategory)
}

func (controller *Controller) DeleteCategoryController() error {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("category ID cannot be empty")
		}
		return nil
	}

	idPrompt := promptui.Prompt{
		Label:    "Category ID to delete",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return err
	}

	categoryId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &categoryId); err != nil {
		return errors.New("invalid category ID format")
	}

	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to delete this category? (y/N)",
		IsConfirm: true,
	}
	_, err = confirmPrompt.Run()
	if err != nil {
		return errors.New("category deletion cancelled")
	}

	return controller.CategoryService.Delete(categoryId)
}

func (controller *Controller) GetAllProductsController() ([]model.ProductModel, error) {
	products, err := controller.ProductService.GetAll()
	if err != nil {
		return []model.ProductModel{}, err
	}

	return products, nil
}

func (controller *Controller) ManageProductsChoiceController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" {
			return errors.New("please choose a valid menu (1/2/3/4)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3/4): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return choice, nil
}

func (controller *Controller) CreateProductController() (model.ProductModel, error) {
	validateName := func(input string) error {
		if len(input) < 2 {
			return errors.New("product name must be at least 2 characters")
		}
		return nil
	}

	validatePrice := func(input string) error {
		if input == "" {
			return errors.New("selling price is required")
		}
		price, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid price format")
		}
		if price < 0 {
			return errors.New("price cannot be negative")
		}
		return nil
	}

	categories, err := controller.CategoryService.GetAll()
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("failed to get categories: %v", err)
	}

	if len(categories) == 0 {
		return model.ProductModel{}, errors.New("no categories available. Please create a category first")
	}

	categoryItems := make([]string, len(categories))
	categoryMap := make(map[string]int)
	for i, category := range categories {
		item := fmt.Sprintf("%s (ID: %d)", category.Name, category.CategoryId)
		categoryItems[i] = item
		categoryMap[item] = category.CategoryId
	}

	namePrompt := promptui.Prompt{
		Label:    "Product Name",
		Validate: validateName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	categoryPrompt := promptui.Select{
		Label: "Select Category",
		Items: categoryItems,
	}
	_, selectedCategory, err := categoryPrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	categoryId := categoryMap[selectedCategory]

	pricePrompt := promptui.Prompt{
		Label:    "Selling Price",
		Validate: validatePrice,
	}
	priceStr, err := pricePrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	price, _ := strconv.ParseFloat(priceStr, 64)

	newProduct := model.ProductModel{
		ProductName:  name,
		CategoryId:   categoryId,
		SellingPrice: price,
	}

	return controller.ProductService.Create(newProduct)
}

func (controller *Controller) UpdateProductController() (model.ProductModel, error) {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("product ID cannot be empty")
		}
		return nil
	}

	validateName := func(input string) error {
		if len(input) < 2 {
			return errors.New("product name must be at least 2 characters")
		}
		return nil
	}

	validatePrice := func(input string) error {
		if input == "" {
			return nil
		}
		price, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid price format")
		}
		if price < 0 {
			return errors.New("price cannot be negative")
		}
		return nil
	}

	idPrompt := promptui.Prompt{
		Label:    "Product ID to update",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	productId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &productId); err != nil {
		return model.ProductModel{}, errors.New("invalid product ID format")
	}

	existingProduct, err := controller.ProductService.GetById(productId)
	if err != nil {
		return model.ProductModel{}, err
	}

	categories, err := controller.CategoryService.GetAll()
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("failed to get categories: %v", err)
	}

	if len(categories) == 0 {
		return model.ProductModel{}, errors.New("no categories available")
	}

	categoryItems := make([]string, len(categories))
	categoryMap := make(map[string]int)
	for i, category := range categories {
		item := fmt.Sprintf("%s (ID: %d)", category.Name, category.CategoryId)
		categoryItems[i] = item
		categoryMap[item] = category.CategoryId
	}

	namePrompt := promptui.Prompt{
		Label:    fmt.Sprintf("Product Name (current: %s)", existingProduct.ProductName),
		Validate: validateName,
		Default:  existingProduct.ProductName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	categoryPrompt := promptui.Select{
		Label: "Select Category",
		Items: categoryItems,
	}
	_, selectedCategory, err := categoryPrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	categoryId := categoryMap[selectedCategory]

	currentPrice := fmt.Sprintf("%.2f", existingProduct.SellingPrice)

	pricePrompt := promptui.Prompt{
		Label:    fmt.Sprintf("Selling Price (current: %s)", currentPrice),
		Validate: validatePrice,
		Default:  currentPrice,
	}
	priceStr, err := pricePrompt.Run()
	if err != nil {
		return model.ProductModel{}, err
	}

	price, _ := strconv.ParseFloat(priceStr, 64)

	updatedProduct := model.ProductModel{
		ProductName:  name,
		CategoryId:   categoryId,
		SellingPrice: price,
	}

	return controller.ProductService.Update(productId, updatedProduct)
}

func (controller *Controller) DeleteProductController() error {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("product ID cannot be empty")
		}
		return nil
	}

	idPrompt := promptui.Prompt{
		Label:    "Product ID to delete",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return err
	}

	productId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &productId); err != nil {
		return errors.New("invalid product ID format")
	}

	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to delete this product? (y/N)",
		IsConfirm: true,
	}
	_, err = confirmPrompt.Run()
	if err != nil {
		return errors.New("product deletion cancelled")
	}

	return controller.ProductService.Delete(productId)
}

func (controller *Controller) GetAllUsersController() ([]model.UserModel, error) {
	users, err := controller.UserService.GetAll()
	if err != nil {
		return []model.UserModel{}, err
	}

	return users, nil
}

func (controller *Controller) ManageUsersChoiceController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" {
			return errors.New("please choose a valid menu (1/2/3/4)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please choose (1/2/3/4): ",
		Validate: validate,
	}

	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return choice, nil
}

func (controller *Controller) CreateUserController() (model.UserModel, error) {
	validateName := func(input string) error {
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

	namePrompt := promptui.Prompt{
		Label:    "Full Name",
		Validate: validateName,
	}
	fullName, err := namePrompt.Run()
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

	adminPrompt := promptui.Select{
		Label: "User Type",
		Items: []string{"Customer", "Admin"},
	}
	_, adminChoice, err := adminPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	isAdmin := adminChoice == "Admin"

	return controller.UserService.Register(fullName, email, password, isAdmin)
}

func (controller *Controller) UpdateUserController() (model.UserModel, error) {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("user ID cannot be empty")
		}
		return nil
	}

	validateName := func(input string) error {
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

	idPrompt := promptui.Prompt{
		Label:    "User ID to update",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	userId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &userId); err != nil {
		return model.UserModel{}, errors.New("invalid user ID format")
	}

	existingUser, err := controller.UserService.GetById(userId)
	if err != nil {
		return model.UserModel{}, err
	}

	namePrompt := promptui.Prompt{
		Label:    fmt.Sprintf("Full Name (current: %s)", existingUser.FullName),
		Validate: validateName,
		Default:  existingUser.FullName,
	}
	fullName, err := namePrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	emailPrompt := promptui.Prompt{
		Label:    fmt.Sprintf("Email (current: %s)", existingUser.Email),
		Validate: validateEmail,
		Default:  existingUser.Email,
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	adminItems := []string{"Customer", "Admin"}
	currentAdminStatus := "Customer"
	if existingUser.IsAdmin {
		currentAdminStatus = "Admin"
	}

	adminPrompt := promptui.Select{
		Label: fmt.Sprintf("User Type (current: %s)", currentAdminStatus),
		Items: adminItems,
	}
	_, adminChoice, err := adminPrompt.Run()
	if err != nil {
		return model.UserModel{}, err
	}

	isAdmin := adminChoice == "Admin"

	updatedUser := model.UserModel{
		FullName: fullName,
		Email:    email,
		IsAdmin:  isAdmin,
	}

	return controller.UserService.Update(userId, updatedUser)
}

func (controller *Controller) DeleteUserController() error {
	validateId := func(input string) error {
		if input == "" {
			return errors.New("user ID cannot be empty")
		}
		return nil
	}

	idPrompt := promptui.Prompt{
		Label:    "User ID to delete",
		Validate: validateId,
	}
	idStr, err := idPrompt.Run()
	if err != nil {
		return err
	}

	userId := 0
	if _, err := fmt.Sscanf(idStr, "%d", &userId); err != nil {
		return errors.New("invalid user ID format")
	}

	confirmPrompt := promptui.Prompt{
		Label:     "Are you sure you want to delete this user? (y/N)",
		IsConfirm: true,
	}
	_, err = confirmPrompt.Run()
	if err != nil {
		return errors.New("user deletion cancelled")
	}

	return controller.UserService.Delete(userId)
}

func (controller *Controller) CreateOrderController(currentUserId int) (model.OrderModel, []model.OrderDetailModel, error) {
	validateAddress := func(input string) error {
		if len(input) < 5 {
			return errors.New("address must be at least 5 characters")
		}
		return nil
	}

	validateAmount := func(input string) error {
		if input == "" {
			return errors.New("amount is required")
		}
		amount, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid amount format")
		}
		if amount < 0 {
			return errors.New("amount cannot be negative")
		}
		return nil
	}

	addressPrompt := promptui.Prompt{
		Label:    "Delivery Address",
		Validate: validateAddress,
	}
	address, err := addressPrompt.Run()
	if err != nil {
		return model.OrderModel{}, nil, err
	}

	shippingPrompt := promptui.Prompt{
		Label:    "Shipping Fee",
		Validate: validateAmount,
		Default:  "10.00",
	}
	shippingStr, err := shippingPrompt.Run()
	if err != nil {
		return model.OrderModel{}, nil, err
	}
	shippingFee, _ := strconv.ParseFloat(shippingStr, 64)

	paymentTypes := []string{"CASH", "TRANSFER", "EWALLET", "CREDIT_CARD"}
	paymentPrompt := promptui.Select{
		Label: "Select Payment Type",
		Items: paymentTypes,
	}
	_, paymentType, err := paymentPrompt.Run()
	if err != nil {
		return model.OrderModel{}, nil, err
	}

	orderStatus := "PENDING"

	products, err := controller.ProductService.GetAll()
	if err != nil {
		return model.OrderModel{}, nil, fmt.Errorf("failed to get products: %v", err)
	}

	if len(products) == 0 {
		return model.OrderModel{}, nil, errors.New("no products available")
	}

	productItems := make([]string, len(products))
	productMap := make(map[string]model.ProductModel)
	for i, product := range products {
		priceStr := fmt.Sprintf("%.2f", product.SellingPrice)
		item := fmt.Sprintf("%s - $%s (ID: %d)", product.ProductName, priceStr, product.ProductId)
		productItems[i] = item
		productMap[item] = product
	}

	var orderDetails []model.OrderDetailModel
	var totalAmount float64 = 0

	fmt.Println("\n==== Add Products to Order ====")
	for {
		productPrompt := promptui.Select{
			Label: "Select Product",
			Items: productItems,
		}
		_, selectedProduct, err := productPrompt.Run()
		if err != nil {
			return model.OrderModel{}, nil, err
		}

		product := productMap[selectedProduct]

		validateQuantity := func(input string) error {
			if input == "" {
				return errors.New("quantity is required")
			}
			qty, err := strconv.Atoi(input)
			if err != nil {
				return errors.New("invalid quantity format")
			}
			if qty <= 0 {
				return errors.New("quantity must be greater than 0")
			}
			return nil
		}

		quantityPrompt := promptui.Prompt{
			Label:    fmt.Sprintf("Quantity for %s", product.ProductName),
			Validate: validateQuantity,
			Default:  "1",
		}
		quantityStr, err := quantityPrompt.Run()
		if err != nil {
			return model.OrderModel{}, nil, err
		}
		quantity, _ := strconv.Atoi(quantityStr)

		orderDetail := model.OrderDetailModel{
			ProductId: product.ProductId,
			Quantity:  quantity,
		}
		orderDetails = append(orderDetails, orderDetail)

		itemTotal := product.SellingPrice * float64(quantity)
		totalAmount += itemTotal

		fmt.Printf("Added: %s x%d = $%.2f\n", product.ProductName, quantity, itemTotal)
		fmt.Printf("Current Total: $%.2f\n", totalAmount)

		continuePrompt := promptui.Select{
			Label: "Add another product?",
			Items: []string{"Yes", "No"},
		}
		_, continueChoice, err := continuePrompt.Run()
		if err != nil {
			return model.OrderModel{}, nil, err
		}

		if continueChoice == "No" {
			break
		}
	}

	if len(orderDetails) == 0 {
		return model.OrderModel{}, nil, errors.New("order must have at least 1 item")
	}

	finalTotal := totalAmount + shippingFee
	fmt.Printf("\nOrder Summary:\n")
	fmt.Printf("Subtotal: $%.2f\n", totalAmount)
	fmt.Printf("Shipping: $%.2f\n", shippingFee)
	fmt.Printf("Total: $%.2f\n", finalTotal)

	confirmPrompt := promptui.Select{
		Label: "Confirm order?",
		Items: []string{"Yes", "Cancel"},
	}
	_, confirmChoice, err := confirmPrompt.Run()
	if err != nil {
		return model.OrderModel{}, nil, err
	}

	if confirmChoice == "Cancel" {
		return model.OrderModel{}, nil, errors.New("order cancelled")
	}

	newOrder := model.OrderModel{
		UserId:      currentUserId,
		Address:     address,
		PaidAmount:  finalTotal,
		ShippingFee: shippingFee,
		PaymentType: paymentType,
		OrderStatus: orderStatus,
	}

	return controller.OrderService.CreateOrderWithDetails(newOrder, orderDetails)
}

func (controller *Controller) ViewOrdersController(userId int) ([]model.OrderModel, error) {
	return controller.OrderService.GetOrdersByUserId(userId)
}

func (controller *Controller) ViewOrderDetailsController(orderId int, userId int) (model.OrderModel, []map[string]interface{}, error) {
	return controller.OrderService.GetOrderDetailsWithProducts(orderId, userId)
}

func (controller *Controller) GetOrderIdInputController() (int, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("order ID cannot be empty")
		}

		orderId, err := strconv.Atoi(input)
		if err != nil || orderId <= 0 {
			return errors.New("please enter a valid order ID (positive number)")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter Order ID",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	orderId, _ := strconv.Atoi(result)
	return orderId, nil
}

//-----------report menu choice  
func (controller *Controller) ReportMenuChoiceController() (string, error) {
	validate := func(input string) error {
		if input != "1" && input != "2" && input != "3" && input != "4" {
			return errors.New("please choose a valid menu (1/2/3/4)")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Please choose (1) Sales SummaryL30d  (2) Top 10 Products  (3) Orders by Status  (4) Back: ",
		Validate: validate,
	}
	choice, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return choice, nil
} 