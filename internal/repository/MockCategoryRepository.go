package repository

import "nextfit/internal/model"

type CategoryStore interface {
	GetAll() ([]model.CategoryModel, error)
	FindByName(name string) (model.CategoryModel, error)
	FindByCategoryId(id int) (model.CategoryModel, error)
	Create(m model.CategoryModel) (model.CategoryModel, error)
	Update(id int, m model.CategoryModel) (model.CategoryModel, error)
	Delete(id int) error
}
