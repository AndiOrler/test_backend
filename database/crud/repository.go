package crud

import (
	"log"

	"gorm.io/gorm"
)

type CrudRepository[T any] interface {
	GetAll() []T
	Get(id uint) T
	Create(model *T) error
	Delete(model *T, id uint) error
}

type CrudRepo[T any] struct {
	DB    *gorm.DB
	Model T
}

func (r *CrudRepo[T]) GetAll() []T {
	log.Printf("getting all from repo")

	var allElements []T

	r.DB.Find(&allElements)

	return allElements
}

func (r *CrudRepo[T]) Get(id uint) T {

	var result T

	log.Printf("Getting entity with id %v", id)

	tx := r.DB.First(&result, id)

	if tx.Error != nil {
		log.Printf("Getting entity failed %v", tx.Error)
	}

	return result

}

func (r *CrudRepo[T]) Delete(model *T, id uint) error {

	tx := r.DB.Delete(&model, id)
	log.Printf("delted %v", id)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *CrudRepo[T]) Create(model *T) error {
	return r.DB.Create(model).Error
}

func NewRepository[T any](db *gorm.DB, model T) CrudRepository[T] {
	return &CrudRepo[T]{
		DB:    db,
		Model: model,
	}
}
