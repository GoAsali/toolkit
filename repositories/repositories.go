package repositories

import (
	"github.com/goasali/toolkit/database"
	"gorm.io/gorm"
	"time"
)

type Interface[T any] interface {
	Create(model *T) error
	CreateMap(model map[string]string) error
	Get(id uint, model *T) error
	Update(id uint, model *T) error
	UpdateMap(id uint, model map[string]any) error
	Delete(id ...uint) error
	ForceDelete(id ...uint) error
	// List get list of dedicated db with a condition
	List(models *[]T, condition ...ListQueryExecuteFn) (int64, error)
}

type Repository[T any] struct {
	model interface{}
}

func NewRepositoryInstance[T any](model interface{}) *Repository[T] {
	return &Repository[T]{
		model: model,
	}
}

func (r Repository[T]) Database() *gorm.DB {
	db, _ := database.Database()
	return db
}

func (r Repository[T]) Delete(id ...uint) error {
	db := r.Database()
	return db.Model(r.model).Delete(r.model, id).Error
}

func (r Repository[T]) ForceDelete(id ...uint) error {
	db := r.Database()
	return db.Model(r.model).Unscoped().Delete(r.model, id).Error
}

func (r Repository[T]) CreateMap(_ map[string]string) error {
	return nil
}

func (r Repository[T]) UpdateMap(id uint, model map[string]any) error {
	db := r.Database()
	return db.Model(r.model).Where("id=?", id).Updates(model).Error
}

func (r Repository[T]) Deactivate(id ...uint) error {
	db, err := database.Database()
	if err != nil {
		return err
	}
	return db.Model(r.model).Where("id", id).Update("deactivated_at", time.Now()).Error
}

func (r Repository[T]) Create(model *T) error {
	db, err := database.Database()
	if err != nil {
		return err
	}
	return db.Create(model).Error
}

func (r Repository[T]) List(models *[]T, queryExecute ...ListQueryExecuteFn) (int64, error) {
	db, err := database.Database()
	if err != nil {
		return 0, err
	}
	query := DefaultListQuery()
	for _, fn := range queryExecute {
		fn(query)
	}

	var count int64

	tx := db.Limit(query.perPage)
	tx.Offset(query.perPage * (query.page - 1))

	if query.order == "" {
		tx.Order(query.order)
	} else {
		tx.Order("id desc")
	}

	for _, condition := range query.conditions {
		tx.Where(condition)
	}

	tx.Find(models)
	tx.Count(&count)

	return count, tx.Error
}

func (r Repository[T]) Get(id uint, model *T) error {
	db, err := database.Database()
	if err != nil {
		return err
	}
	return db.Where("id=?", id).First(model).Error
}

func (r Repository[T]) Update(id uint, model *T) error {
	db, err := database.Database()
	if err != nil {
		return err
	}
	return db.Model(r.model).Where("id=?", id).Updates(model).Error
}
