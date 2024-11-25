package repositories

import (
	"github.com/google/uuid"
	custom "github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) (uuid.UUID, error)
	GetAll() ([]models.BookOut, error)
	Update(id uuid.UUID, book *models.Book) error
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

func (repository *bookRepository) GetAll() ([]models.BookOut, error) {
	var books []models.BookOut

	result := repository.db.Raw(`select b.id, b.title, b.edition, b.publication_year, array_agg(a.name) as authors from book_author ba inner join books b on ba.book_id = b.id inner join authors a on ba.author_id = a.id group by b.id;`).Scan(&books)

	if err := result.Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (repository *bookRepository) GetBookByID(id uuid.UUID) (models.BookOut, error) {
	var book models.BookOut

	result := repository.db.Raw(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id
INNER JOIN authors a ON ba.author_id = a.id WHERE ba.book_id = ? GROUP BY b.id ORDER BY b.id LIMIT 1;`, id).Scan(&book)

	if err := result.Error; err != nil {
		return models.BookOut{}, err
	}

	if result.RowsAffected < 1 {
		return models.BookOut{}, &custom.BookNotFound
	}

	return book, nil
}

func (repository *bookRepository) Update(id uuid.UUID, book *models.Book) error {
	result := repository.db.Model(&models.Book{}).Where("id = ?", id).Updates(&models.Book{Title: book.Title, Edition: book.Edition, PublicationYear: book.PublicationYear})

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected < 1 {
		return &custom.BookNothingToUpdate
	}

	return nil
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
