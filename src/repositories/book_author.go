package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	custom "github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
)

type BookAuthorRepository interface {
	Create(relationship *models.BookAuthor) error
	Delete(bookID uuid.UUID) error
}

type bookAuthorRepository struct {
	db *gorm.DB
}

func NewBookAuthorRepository(db *gorm.DB) BookAuthorRepository {
	return &bookAuthorRepository{db}
}

func (repository *bookAuthorRepository) Create(relationship *models.BookAuthor) error {
	if err := repository.db.Create(&relationship).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &custom.RelationshipAlreadyExists
		}
		return err
	}
	return nil
}

func (repository *bookAuthorRepository) Delete(bookID uuid.UUID) error {
	if err := repository.db.Delete(&models.BookAuthor{}, "book_id = ?", bookID).Error; err != nil {
		return err
	}
	return nil
}
