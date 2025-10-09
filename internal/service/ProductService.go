package service

import (
	"fmt"
	"nextfit/internal/model"
	"nextfit/internal/repository"
	"strings"
)

type ProductService struct {
	ProductRepository  *repository.ProductRepository
	CategoryRepository *repository.CategoryRepository
}

func (productService *ProductService) GetAll() ([]model.ProductModel, error) {
	products, err := productService.ProductRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}

	return products, nil
}

func (productService *ProductService) Create(newProduct model.ProductModel) (model.ProductModel, error) {
	if strings.TrimSpace(newProduct.ProductName) == "" {
		return model.ProductModel{}, fmt.Errorf("product name cannot be empty")
	}

	if newProduct.CategoryId <= 0 {
		return model.ProductModel{}, fmt.Errorf("invalid category ID")
	}

	_, err := productService.CategoryRepository.FindByCategoryId(newProduct.CategoryId)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("category not found: %v", err)
	}

	if newProduct.SellingPrice < 0 {
		return model.ProductModel{}, fmt.Errorf("selling price cannot be negative")
	}

	newProduct.ProductName = strings.TrimSpace(newProduct.ProductName)

	createdProduct, err := productService.ProductRepository.Create(newProduct)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("failed to create product: %v", err)
	}

	return createdProduct, nil
}

func (productService *ProductService) Update(productId int, updatedProduct model.ProductModel) (model.ProductModel, error) {
	if productId <= 0 {
		return model.ProductModel{}, fmt.Errorf("invalid product ID")
	}

	if strings.TrimSpace(updatedProduct.ProductName) == "" {
		return model.ProductModel{}, fmt.Errorf("product name cannot be empty")
	}

	if updatedProduct.CategoryId <= 0 {
		return model.ProductModel{}, fmt.Errorf("invalid category ID")
	}

	_, err := productService.ProductRepository.FindByProductId(productId)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("product not found: %v", err)
	}

	_, err = productService.CategoryRepository.FindByCategoryId(updatedProduct.CategoryId)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("category not found: %v", err)
	}

	if updatedProduct.SellingPrice < 0 {
		return model.ProductModel{}, fmt.Errorf("selling price cannot be negative")
	}

	updatedProduct.ProductName = strings.TrimSpace(updatedProduct.ProductName)
	updatedProduct.ProductId = productId

	result, err := productService.ProductRepository.Update(productId, updatedProduct)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("failed to update product: %v", err)
	}

	return result, nil
}

func (productService *ProductService) Delete(productId int) error {
	if productId <= 0 {
		return fmt.Errorf("invalid product ID")
	}

	_, err := productService.ProductRepository.FindByProductId(productId)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	err = productService.ProductRepository.Delete(productId)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	return nil
}

func (productService *ProductService) GetById(productId int) (model.ProductModel, error) {
	if productId <= 0 {
		return model.ProductModel{}, fmt.Errorf("invalid product ID")
	}

	product, err := productService.ProductRepository.FindByProductId(productId)
	if err != nil {
		return model.ProductModel{}, fmt.Errorf("product not found: %v", err)
	}

	return product, nil
}

func (productService *ProductService) GetByCategory(categoryId int) ([]model.ProductModel, error) {
	if categoryId <= 0 {
		return nil, fmt.Errorf("invalid category ID")
	}

	_, err := productService.CategoryRepository.FindByCategoryId(categoryId)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	products, err := productService.ProductRepository.FindByCategoryId(categoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %v", err)
	}

	return products, nil
}
