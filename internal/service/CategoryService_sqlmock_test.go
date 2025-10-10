package service_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"nextfit/internal/model"
	"nextfit/internal/repository"
	"nextfit/internal/service"
)

func rx(s string) string { return regexp.QuoteMeta(s) }

func TestCategory_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_name = \?`)).
		WithArgs("Soccer").
		WillReturnRows(sqlmock.NewRows([]string{"category_id"}))


	mock.ExpectExec(rx(`INSERT INTO categories`)).
		WithArgs("Soccer", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(123, 1))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	got, err := svc.Create("  Soccer  ")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.CategoryId == 0 || got.Name != "Soccer" {
		t.Fatalf("unexpected result: %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Create_EmptyName(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	if _, err := svc.Create("   "); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCategory_Create_Duplicate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_name = \?`)).
		WithArgs("Soccer").
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name"}).
			AddRow(7, "Soccer"))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	if _, err := svc.Create("Soccer"); err == nil {
		t.Fatal("expected duplicate name error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Update_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	id := 10

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name", "created_at", "updated_at", "deleted_at"}).
			AddRow(id, "OldName", time.Now(), nil, nil))

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_name = \?`)).
		WithArgs("NewName").
		WillReturnRows(sqlmock.NewRows([]string{"category_id"}))

	mock.ExpectExec(rx(`UPDATE categories SET category_name = \?, updated_at = \? WHERE category_id = \?`)).
		WithArgs("NewName", sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	got, err := svc.Update(id, model.CategoryModel{Name: "  NewName "})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.Name != "NewName" {
		t.Fatalf("unexpected updated name: %+v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Update_InvalidID(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	_, err := svc.Update(0, model.CategoryModel{Name: "X"})
	if err == nil {
		t.Fatal("expected invalid id error")
	}
}

func TestCategory_Update_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	id := 99

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"category_id"}))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	_, err := svc.Update(id, model.CategoryModel{Name: "Any"})
	if err == nil {
		t.Fatal("expected not found error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Update_NameTakenByAnother(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	id := 10

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name"}).
			AddRow(id, "Old"))


	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_name = \?`)).
		WithArgs("Taken").
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name"}).
			AddRow(777, "Taken"))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	_, err := svc.Update(id, model.CategoryModel{Name: "Taken"})
	if err == nil {
		t.Fatal("expected duplicate name error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Delete_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	id := 5

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "category_name"}).
			AddRow(id, "X"))

	mock.ExpectExec(rx(`DELETE FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	if err := svc.Delete(id); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategory_Delete_InvalidID(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	if err := svc.Delete(0); err == nil {
		t.Fatal("expected invalid id error")
	}
}

func TestCategory_Delete_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	id := 77

	mock.ExpectQuery(rx(`SELECT .* FROM categories WHERE category_id = \?`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"category_id"}))

	repo := &repository.CategoryRepository{DB: db}
	svc := &service.CategoryService{CategoryRepository: repo}

	if err := svc.Delete(id); err == nil {
		t.Fatal("expected not found error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
