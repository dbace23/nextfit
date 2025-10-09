package service

import (
	"errors"
	"fmt"
	"nextfit/internal/model"
	"nextfit/internal/repository"
	"strings"
)

type CategoryService struct {
	CategoryRepository *repository.CategoryRepository
}

func (categoryService *CategoryService) GetAll() ([]model.CategoryModel, error) {
	categories, err := categoryService.CategoryRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (categoryService *CategoryService) Create(name string) (model.CategoryModel, error) {
	trimName := strings.TrimSpace(name)
	if trimName == "" || len(trimName) < 1 {
		return model.CategoryModel{}, errors.New("category name cannot be empty")
	}

	_, err := categoryService.CategoryRepository.FindByName(trimName)
	if err == nil {
		return model.CategoryModel{}, fmt.Errorf("category with name '%s' already exists", name)
	}

	newCategory := model.CategoryModel{
		Name: trimName,
	}

	createdCategory, err := categoryService.CategoryRepository.Create(newCategory)
	if err != nil {
		return model.CategoryModel{}, fmt.Errorf("failed to create category: %v", err)
	}

	return createdCategory, nil
}

func (categoryService *CategoryService) Update(categoryId int, updatedCategory model.CategoryModel) (model.CategoryModel, error) {
	if categoryId <= 0 {
		return model.CategoryModel{}, fmt.Errorf("invalid category ID")
	}

	trimName := strings.TrimSpace(updatedCategory.Name)
	if trimName == "" {
		return model.CategoryModel{}, fmt.Errorf("category name cannot be empty")
	}

	_, err := categoryService.CategoryRepository.FindByCategoryId(categoryId)
	if err != nil {
		return model.CategoryModel{}, fmt.Errorf("category not found: %v", err)
	}

	categoryWithSameName, err := categoryService.CategoryRepository.FindByName(trimName)
	if err == nil && categoryWithSameName.CategoryId != categoryId {
		return model.CategoryModel{}, fmt.Errorf("category with name '%s' already exists", trimName)
	}

	updatedCategory.Name = trimName
	updatedCategory.CategoryId = categoryId

	result, err := categoryService.CategoryRepository.Update(categoryId, updatedCategory)
	if err != nil {
		return model.CategoryModel{}, fmt.Errorf("failed to update category: %v", err)
	}

	return result, nil
}

func (categoryService *CategoryService) Delete(categoryId int) error {
	if categoryId <= 0 {
		return fmt.Errorf("invalid category ID")
	}

	_, err := categoryService.CategoryRepository.FindByCategoryId(categoryId)
	if err != nil {
		return fmt.Errorf("category not found: %v", err)
	}

	err = categoryService.CategoryRepository.Delete(categoryId)
	if err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	return nil
}

func (categoryService *CategoryService) GetById(categoryId int) (model.CategoryModel, error) {
	if categoryId <= 0 {
		return model.CategoryModel{}, fmt.Errorf("invalid category ID")
	}

	category, err := categoryService.CategoryRepository.FindByCategoryId(categoryId)
	if err != nil {
		return model.CategoryModel{}, fmt.Errorf("category not found: %v", err)
	}

	return category, nil
}
