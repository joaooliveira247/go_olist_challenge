package repositories

import (
	"github.com/google/uuid"

	"github.com/joaooliveira247/go_olist_challenge/src/models"
)

type BookAuthorRepository interface {
	Create(relationship *models.BookAuthor) error
	Delete(bookID uuid.UUID) error
}
