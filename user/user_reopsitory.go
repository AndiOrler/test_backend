package user

import (
	"test_backend/database/crud"
	"test_backend/models"

	"gorm.io/gorm"
)

type model = models.User

type UserRepository interface {
	crud.CrudRepository[model]
}

type userRepo struct {
	crud.CrudRepo[model]
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepo{
		CrudRepo: crud.CrudRepo[model]{
			DB:    db,
			Model: model{},
		},
	}
}
