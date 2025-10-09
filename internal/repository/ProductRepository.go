package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nextfit/internal/model"
)

type ProductRepository struct {
	DB *sql.DB
}

func (productRepository *ProductRepository) Create(newProduct model.ProductModel) (model.ProductModel, error) {
	query := `
    INSERT INTO products (
        product_name, category_id, selling_price
    ) VALUES (?, ?, ?)
    `

	result, err := productRepository.DB.ExecContext(context.Background(), query,
		newProduct.ProductName,
		newProduct.CategoryId,
		newProduct.SellingPrice,
	)

	if err != nil {
		fmt.Println(err)
		return model.ProductModel{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return model.ProductModel{}, err
	}

	newProduct.ProductId = int(lastInsertId)
	return newProduct, nil
}

func (productRepository *ProductRepository) GetAll() ([]model.ProductModel, error) {
	query := `
	SELECT
	product_id, product_name, category_id, selling_price, created_at, updated_at, deleted_at
	FROM products
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC
	`

	result, err := productRepository.DB.QueryContext(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer result.Close()

	var products []model.ProductModel
	for result.Next() {
		var product model.ProductModel
		err := result.Scan(
			&product.ProductId,
			&product.ProductName,
			&product.CategoryId,
			&product.SellingPrice,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (productRepository *ProductRepository) Update(productId int, updatedProduct model.ProductModel) (model.ProductModel, error) {
	query := `
	UPDATE products SET
	product_name = ?, category_id = ?, selling_price = ?, updated_at = NOW()
	WHERE product_id = ? AND deleted_at IS NULL
	`

	result, err := productRepository.DB.ExecContext(context.Background(), query,
		updatedProduct.ProductName,
		updatedProduct.CategoryId,
		updatedProduct.SellingPrice,
		productId,
	)

	if err != nil {
		fmt.Println(err)
		return model.ProductModel{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.ProductModel{}, err
	}

	if rowsAffected == 0 {
		return model.ProductModel{}, fmt.Errorf("product with id %d not found", productId)
	}

	updatedProduct.ProductId = productId
	return updatedProduct, nil
}

func (productRepository *ProductRepository) Delete(productId int) error {
	query := `
	UPDATE products SET
	deleted_at = NOW()
	WHERE product_id = ? AND deleted_at IS NULL
	`

	result, err := productRepository.DB.ExecContext(context.Background(), query, productId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", productId)
	}

	return nil
}

func (productRepository *ProductRepository) FindByProductId(productId int) (model.ProductModel, error) {
	query := `
	SELECT
	product_id, product_name, category_id, selling_price, created_at, updated_at, deleted_at
	FROM products
	WHERE product_id = ? AND deleted_at IS NULL
	`

	result, err := productRepository.DB.QueryContext(context.Background(), query, productId)
	if err != nil {
		return model.ProductModel{}, err
	}
	defer result.Close()

	var product model.ProductModel
	if result.Next() {
		err := result.Scan(
			&product.ProductId,
			&product.ProductName,
			&product.CategoryId,
			&product.SellingPrice,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
		)
		if err != nil {
			return model.ProductModel{}, err
		}
	} else {
		return model.ProductModel{}, fmt.Errorf("product with id %d not found", productId)
	}

	return product, nil
}

func (productRepository *ProductRepository) FindByCategoryId(categoryId int) ([]model.ProductModel, error) {
	query := `
	SELECT
	product_id, product_name, category_id, selling_price, created_at, updated_at, deleted_at
	FROM products
	WHERE category_id = ? AND deleted_at IS NULL
	ORDER BY created_at DESC
	`

	result, err := productRepository.DB.QueryContext(context.Background(), query, categoryId)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var products []model.ProductModel
	for result.Next() {
		var product model.ProductModel
		err := result.Scan(
			&product.ProductId,
			&product.ProductName,
			&product.CategoryId,
			&product.SellingPrice,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
