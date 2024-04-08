package crud

import (
	"log"

	"gorm.io/gorm"
)

type CrudRepository[T any] interface {
	GetAll(model *T) []T
	Create(model *T) error
}

type CrudRepo[T any] struct {
	DB    *gorm.DB
	Model T
}

func (r *CrudRepo[T]) GetAll(model *T) []T {
	log.Printf("getting all %v from crud repo", model)

	var allElements []T

	r.DB.Find(&allElements)

	return allElements
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
