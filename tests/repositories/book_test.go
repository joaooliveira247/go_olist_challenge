package repositories_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	book := models.Book{
		Title:           "The Rust Programming Language",
		Edition:         1,
		PublicationYear: 2018,
	}

	bookID := uuid.New()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "books" WHERE "books"."title" = $1 AND "books"."edition" = $2 AND "books"."publication_year" = $3 ORDER BY "books"."id" LIMIT $4`,
		),
	).WithArgs(book.Title, book.Edition, book.PublicationYear, 1).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "books" ("title","edition","publication_year") VALUES ($1,$2,$3) RETURNING "id"`,
		),
	).WithArgs(book.Title, book.Edition, book.PublicationYear).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(bookID))
	mock.ExpectCommit()

	repository := repositories.NewBookRepository(gormDB)

	id, err := repository.Create(&book)

	assert.Nil(t, err)
	assert.Equal(t, bookID, id)
}
