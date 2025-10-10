package service_test

import (
	"errors"
	"testing"

	"nextfit/internal/model"
	"nextfit/internal/service"
)

type mockCategoryStore struct {
	getAllFn           func() ([]model.CategoryModel, error)
	findByNameFn       func(string) (model.CategoryModel, error)
	findByCategoryIdFn func(int) (model.CategoryModel, error)
	createFn           func(model.CategoryModel) (model.CategoryModel, error)
	updateFn           func(int, model.CategoryModel) (model.CategoryModel, error)
	deleteFn           func(int) error
}

func (m *mockCategoryStore) GetAll() ([]model.CategoryModel, error) {
	if m.getAllFn != nil { return m.getAllFn() }
	return nil, nil
}
func (m *mockCategoryStore) FindByName(name string) (model.CategoryModel, error) {
	if m.findByNameFn != nil { return m.findByNameFn(name) }
	return model.CategoryModel{}, errors.New("not found")
}
func (m *mockCategoryStore) FindByCategoryId(id int) (model.CategoryModel, error) {
	if m.findByCategoryIdFn != nil { return m.findByCategoryIdFn(id) }
	return model.CategoryModel{}, errors.New("not found")
}
func (m *mockCategoryStore) Create(c model.CategoryModel) (model.CategoryModel, error) {
	if m.createFn != nil { return m.createFn(c) }
	return c, nil
}
func (m *mockCategoryStore) Update(id int, c model.CategoryModel) (model.CategoryModel, error) {
	if m.updateFn != nil { return m.updateFn(id, c) }
	return c, nil
}
func (m *mockCategoryStore) Delete(id int) error {
	if m.deleteFn != nil { return m.deleteFn(id) }
	return nil
}

func TestCategoryService_Create_Success(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByNameFn: func(name string) (model.CategoryModel, error) {
			return model.CategoryModel{}, errors.New("not found") // name can use
		},
		createFn: func(c model.CategoryModel) (model.CategoryModel, error) {
			c.CategoryId = 123
			return c, nil
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}

	got, err := svc.Create("  Soccer ")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.CategoryId != 123 || got.Name != "Soccer" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestCategoryService_Create_Empty(t *testing.T) {
	svc := &service.CategoryService{CategoryRepository: &mockCategoryStore{}}
	if _, err := svc.Create("   "); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCategoryService_Create_Duplicate(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByNameFn: func(name string) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: 7, Name: "Soccer"}, nil // already exists
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	if _, err := svc.Create("Soccer"); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestCategoryService_Update_Success(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: id, Name: "Old"}, nil // exists
		},
		findByNameFn: func(name string) (model.CategoryModel, error) {
			return model.CategoryModel{}, errors.New("not found") // new name can use
		},
		updateFn: func(id int, c model.CategoryModel) (model.CategoryModel, error) {
			return c, nil
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	got, err := svc.Update(10, model.CategoryModel{Name: "  NewName "})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.Name != "NewName" || got.CategoryId != 10 {
		t.Fatalf("bad update: %+v", got)
	}
}

func TestCategoryService_Update_InvalidID(t *testing.T) {
	svc := &service.CategoryService{CategoryRepository: &mockCategoryStore{}}
	if _, err := svc.Update(0, model.CategoryModel{Name: "X"}); err == nil {
		t.Fatal("expected invalid id error")
	}
}

func TestCategoryService_Update_NameTakenByAnother(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: id, Name: "Old"}, nil
		},
		findByNameFn: func(name string) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: 777, Name: "Taken"}, nil // different ID
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	if _, err := svc.Update(10, model.CategoryModel{Name: "Taken"}); err == nil {
		t.Fatal("expected name taken error")
	}
}

func TestCategoryService_Delete_Success(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: id, Name: "X"}, nil
		},
		deleteFn: func(id int) error { return nil },
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	if err := svc.Delete(5); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestCategoryService_Delete_InvalidID(t *testing.T) {
	svc := &service.CategoryService{CategoryRepository: &mockCategoryStore{}}
	if err := svc.Delete(0); err == nil {
		t.Fatal("expected invalid id error")
	}
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{}, errors.New("not found")
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	if err := svc.Delete(77); err == nil {
		t.Fatal("expected not found error")
	}
}

func TestCategoryService_GetById_Success(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{CategoryId: id, Name: "Soccer"}, nil
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	got, err := svc.GetById(7)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.CategoryId != 7 || got.Name != "Soccer" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestCategoryService_GetById_InvalidID(t *testing.T) {
	svc := &service.CategoryService{CategoryRepository: &mockCategoryStore{}}
	if _, err := svc.GetById(0); err == nil {
		t.Fatal("expected invalid id error")
	}
}

func TestCategoryService_GetById_NotFound(t *testing.T) {
	mockRepo := &mockCategoryStore{
		findByCategoryIdFn: func(id int) (model.CategoryModel, error) {
			return model.CategoryModel{}, errors.New("not found")
		},
	}
	svc := &service.CategoryService{CategoryRepository: mockRepo}
	if _, err := svc.GetById(99); err == nil {
		t.Fatal("expected not found error")
	}
}
