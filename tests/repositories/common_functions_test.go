package repositories_test

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
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

var SingleBookMocked = models.Book{
	Title:           "The Rust Programming Language",
	Edition:         1,
	PublicationYear: 2018,
}
