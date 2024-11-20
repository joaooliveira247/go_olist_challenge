package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joaooliveira247/go_olist_challenge/src/models"
)

type BookAuthorRepository interface {
	Create(relationship *models.BookAuthor) error
	Delete(bookID uuid.UUID) error
}

type bookAuthorRepository struct {
	db *gorm.DB
}

