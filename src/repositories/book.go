package repositories

import (
	"github.com/google/uuid"
	custom "github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) (uuid.UUID, error)
	Delete(id uuid.UUID) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (repository *bookRepository) Create(book *models.Book) (uuid.UUID, error) {
	result := repository.db.FirstOrCreate(&book, book)

	if err := result.Error; err != nil {
		return uuid.Nil, err
	}

	if result.RowsAffected < 1 {
		return uuid.Nil, &custom.BookAlreadyExists
	}

	return book.ID, nil
}

func (repository *bookRepository) Delete(id uuid.UUID) error {
	result := repository.db.Delete(&models.Book{}, id)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected < 1 {
		return &custom.BookNotFound
	}

	return nil
}
