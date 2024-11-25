package repositories_test

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return gormDB, mock
}

func NewMockBook() *models.Book {
	return &models.Book{
		Title:           "The Rust Programming Language",
		Edition:         1,
		PublicationYear: 2018,
	}
}

func NewMockBookOut() models.BookOut {
	return models.BookOut{
		Book: models.Book{
			ID: uuid.New(),
			Title: "the Rust Programming Language",
			Edition: 1,
			PublicationYear: 2018,
		},
		AuthorsName: pq.StringArray{"Carol Nichols", "Steve Klabnik"},
	}
}

func NewMockBooks() []models.BookOut {
	return []models.BookOut{
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "The Rust Programming Language",
				Edition:         1,
				PublicationYear: 2018,
			},
			AuthorsName: pq.StringArray{"Carol Nichols", "Steve Klabnik"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "The Go Programming Language",
				Edition:         2,
				PublicationYear: 2015,
			},
			AuthorsName: pq.StringArray{"Alan A. A. Donovan", "Brian W. Kernighan"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "Python Fluente",
				Edition:         1,
				PublicationYear: 2015,
			},
			AuthorsName: pq.StringArray{"Luciano Ramalho"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "Python Fluente",
				Edition:         2,
				PublicationYear: 2022,
			},
			AuthorsName: pq.StringArray{"Luciano Ramalho"},
		},
	}
}
