package repository

import (
	"context"
	"database/sql"
	"fmt"
	"nextfit/internal/model"
)

type CategoryRepository struct {
	DB *sql.DB
}

func (categoryRepository *CategoryRepository) Create(newCategory model.CategoryModel) (model.CategoryModel, error) {
	query := `
    INSERT INTO categories (
        category_name
    ) VALUES (?)
    `

	result, err := categoryRepository.DB.ExecContext(context.Background(), query,
		newCategory.Name,
	)

	if err != nil {
		return model.CategoryModel{}, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return model.CategoryModel{}, err
	}

	newCategory.CategoryId = int(lastInsertId)
	return newCategory, nil
}

func (categoryRepository *CategoryRepository) GetAll() ([]model.CategoryModel, error) {
	query := `
    SELECT
    category_id, category_name, created_at, updated_at, deleted_at
    FROM categories
    WHERE deleted_at IS NULL
    ORDER BY created_at DESC
    `
	result, err := categoryRepository.DB.QueryContext(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer result.Close()

	var categories []model.CategoryModel
	for result.Next() {
		var category model.CategoryModel
		err := result.Scan(
			&category.CategoryId,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (categoryRepository *CategoryRepository) Update(categoryId int, updatedCategory model.CategoryModel) (model.CategoryModel, error) {
	query := `
    UPDATE categories SET
    category_name = ?, updated_at = NOW()
    WHERE category_id = ? AND deleted_at IS NULL
    `

	result, err := categoryRepository.DB.ExecContext(context.Background(), query,
		updatedCategory.Name,
		categoryId,
	)

	if err != nil {
		fmt.Println(err)
		return model.CategoryModel{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.CategoryModel{}, err
	}

	if rowsAffected == 0 {
		return model.CategoryModel{}, fmt.Errorf("category with id %d not found", categoryId)
	}

	updatedCategory.CategoryId = categoryId
	return updatedCategory, nil
}

func (categoryRepository *CategoryRepository) Delete(categoryId int) error {
	query := `
    UPDATE categories SET
    deleted_at = NOW()
    WHERE category_id = ? AND deleted_at IS NULL
    `

	result, err := categoryRepository.DB.ExecContext(context.Background(), query, categoryId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", categoryId)
	}

	return nil
}

func (categoryRepository *CategoryRepository) FindByCategoryId(categoryId int) (model.CategoryModel, error) {
	query := `
    SELECT
    category_id, category_name, created_at, updated_at, deleted_at
    FROM categories
    WHERE category_id = ? AND deleted_at IS NULL
    `

	result, err := categoryRepository.DB.QueryContext(context.Background(), query, categoryId)
	if err != nil {
		return model.CategoryModel{}, err
	}
	defer result.Close()

	var category model.CategoryModel
	if result.Next() {
		err := result.Scan(
			&category.CategoryId,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		)
		if err != nil {
			return model.CategoryModel{}, err
		}
	} else {
		return model.CategoryModel{}, fmt.Errorf("category with id %d not found", categoryId)
	}

	return category, nil
}

func (categoryRepository *CategoryRepository) FindByName(name string) (model.CategoryModel, error) {
	query := `
    SELECT
    category_id, category_name, created_at, updated_at, deleted_at
    FROM categories
    WHERE category_name = ? 
    `

	result, err := categoryRepository.DB.QueryContext(context.Background(), query, name)
	if err != nil {
		return model.CategoryModel{}, err
	}
	defer result.Close()

	var category model.CategoryModel
	if result.Next() {
		err := result.Scan(
			&category.CategoryId,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		)
		if err != nil {
			return model.CategoryModel{}, err
		}
	} else {
		return model.CategoryModel{}, fmt.Errorf("category with name '%s' not found", name)
	}

	return category, nil
}
