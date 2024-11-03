package repositories

import (
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(author *models.Author) (uuid.UUID, error)
}

type authorRepository struct {
	db *gorm.DB
}
