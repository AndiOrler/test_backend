package crud

import (
	"log"
)

type Crud interface {
	// GetAll() []models.User
	GetAll()
}

type Repository struct {
}

func (r Repository) GetAll() {
	log.Println("getting all users from crud repo")
}

// func (r *Repository) GetAll() ([]models.User, error) {

// 	log.Println("returning all users from crud repo")

// 	return nil, nil
// }
