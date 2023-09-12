package services

import (
	"github.com/goasali/toolkit/repositories"
)

type ServiceI[M any] interface {
	// Delete model from database (Soft delete supported)
	Delete(id ...uint) error
	// Create new model in database
	Create(data *M) error
	// Update set new values for database
	Update(id uint, model *M) error
	// List get list of
	List(models *[]M, optionCallbacks ...ListOptionsFunc) (int64, error)
	FindById(id uint, model *M) error
	Activate(id ...uint) error
	Repository() repositories.Interface[M]
}

type AbstractService[M any] struct {
	repository repositories.Interface[M]
}

func NewAbstractService[M any](repo repositories.Interface[M]) *AbstractService[M] {
	return &AbstractService[M]{
		repository: repo,
	}
}

func (as AbstractService[M]) Repository() repositories.Interface[M] {
	return as.repository
}

func (as AbstractService[M]) Create(data *M) error {
	if err := as.repository.Create(data); err != nil {
		return err
	}
	return nil
}

func (as AbstractService[M]) Update(id uint, model *M) error {
	if err := as.repository.Update(id, model); err != nil {
		return err
	}
	return nil
}

func (as AbstractService[M]) FindById(id uint, model *M) error {
	if err := as.repository.Get(id, model); err != nil {
		return err
	}
	return nil
}

func (as AbstractService[M]) Delete(id ...uint) error {
	if err := as.repository.Delete(id...); err != nil {
		return err
	}
	return nil
}

func (as AbstractService[M]) Activate(...uint) error {
	panic("implement me")
}
