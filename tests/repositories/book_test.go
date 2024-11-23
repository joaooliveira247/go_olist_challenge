package repositories_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	book := NewMockBook()

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

	id, err := repository.Create(book)

	assert.Nil(t, err)
	assert.Equal(t, bookID, id)
}

func TestCreateBookReturnAlredyExists(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	book := NewMockBook()

	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM "books" WHERE "books"."title" = $1 AND "books"."edition" = $2 AND "books"."publication_year" = $3 ORDER BY "books"."id" LIMIT $4`,
		),
	).WithArgs(book.Title, book.Edition, book.PublicationYear, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "edition", "publication_year"}).AddRow(bookID, book.Title, book.Edition, book.PublicationYear))

	repository := repositories.NewBookRepository(gormDB)
	id, err := repository.Create(book)

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookAlreadyExists)
	assert.Equal(t, uuid.Nil, id)
}

func TestCreateBookReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	book := NewMockBook()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."title" = $1 AND "books"."edition" = $2 AND "books"."publication_year" = $3 ORDER BY "books"."id" LIMIT $4`)).WithArgs(book.Title, book.Edition, book.PublicationYear, 1).WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "books" ("title","edition","publication_year") VALUES ($1,$2,$3) RETURNING "id"`,
		),
	).WithArgs(book.Title, book.Edition, book.PublicationYear).WillReturnError(&errors.BookGenericError)
	mock.ExpectRollback()

	repository := repositories.NewBookRepository(gormDB)
	id, err := repository.Create(book)

	assert.Error(t, err)
	t.Log(err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
	assert.Equal(t, uuid.Nil, id)
}

func TestDeleteBookSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).WithArgs(bookID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := repositories.NewBookRepository(gormDB)
	err := repository.Delete(bookID)

	assert.Nil(t, err)
}
