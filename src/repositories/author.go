package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	custom "github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(author *models.Author) (uuid.UUID, error)
	CreateMany(authors *[]models.Author) ([]uuid.UUID, error)
	GetAll() ([]models.Author, error)
	GetByID(id uuid.UUID) (models.Author, error)
	GetByName(name string) ([]models.Author, error)
	Delete(id uuid.UUID) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) authorRepository {
	return authorRepository{db}
}

func (repository *authorRepository) Create(author *models.Author) (uuid.UUID, error) {
	result := repository.db.Create(&author)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return uuid.UUID{}, &custom.AuthorAlreadyExists
		}
		return uuid.UUID{}, err
	}

	return author.ID, nil
}

func (repository *authorRepository) CreateMany(authors *[]models.Author) ([]uuid.UUID, error) {
	result := repository.db.Create(&authors)

	if err := result.Error; err != nil {
		return nil, err
	}

	var authorsIDs []uuid.UUID

	for _, author := range *authors {
		authorsIDs = append(authorsIDs, author.ID)
	}

	return authorsIDs, nil
}

func (repository *authorRepository) GetAll() ([]models.Author, error) {
	var authors []models.Author

	result := repository.db.Find(&authors)

	if err := result.Error; err != nil {
		return nil, err
	}

	return authors, nil
}

func (repository *authorRepository) GetByID(id uuid.UUID) (models.Author, error) {
	var author models.Author

	result := repository.db.First(&author, "id = ?", id)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Author{}, &custom.AuthorNotFound
		}
		return models.Author{}, err
	}

	return author, nil
}

func (repository *authorRepository) GetByName(name string) ([]models.Author, error) {
	var authors []models.Author

	result := repository.db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Find(&authors)

	if err := result.Error; err != nil {
		return nil, err
	}

	return authors, nil
}

func (repository *authorRepository) Delete(id uuid.UUID) error {
	result := repository.db.Delete(models.Author{}, id)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected < 1 {
		return &custom.AuthorNotFound
	}

	return nil
}
