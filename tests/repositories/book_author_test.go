package repositories_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/joaooliveira247/go_olist_challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookAuthorSuccess(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()
	authorID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(
			`INSERT INTO "book_authors" ("book_id","author_id") VALUES ($1,$2)`,
		),
	).WithArgs(bookID, authorID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := repositories.NewBookAuthorRepository(gormDB)
	err := repository.Create(&models.BookAuthor{BookID: bookID, AuthorID: authorID})

	assert.Nil(t, err)
}

func TestCreateBookAuthorReturnGenericError(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()
	authorID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(
			`INSERT INTO "book_authors" ("book_id","author_id") VALUES ($1,$2)`,
		),
	).WithArgs(bookID, authorID).WillReturnError(&errors.BookAuthorGenericError)
	mock.ExpectRollback()

	repository := repositories.NewBookAuthorRepository(gormDB)
	err := repository.Create(&models.BookAuthor{BookID: bookID, AuthorID: authorID})

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookAuthorGenericError)
}

func TestCreateBookAuthorReturnRelationshipAlreadyExists(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()
	authorID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "book_authors" ("book_id","author_id") VALUES ($1,$2)`),
	).WithArgs(bookID, authorID).WillReturnError(&errors.RelationshipAlreadyExists)
	mock.ExpectRollback()

	repository := repositories.NewBookAuthorRepository(gormDB)
	err := repository.Create(&models.BookAuthor{BookID: bookID, AuthorID: authorID})

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.RelationshipAlreadyExists)
}

func TestDeleteBookAuthorSuccess(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(
			`DELETE FROM "book_authors" WHERE book_id = $1`,
		),
	).WithArgs(bookID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := repositories.NewBookAuthorRepository(gormDB)
	err := repository.Delete(bookID)
	assert.Nil(t, err)
}

func TestDeleteBookAuthorReturnGenericError(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "book_authors" WHERE book_id = $1`),
	).WithArgs(bookID).WillReturnError(&errors.BookAuthorGenericError)
	mock.ExpectRollback()

	repository := repositories.NewBookAuthorRepository(gormDB)
	err := repository.Delete(bookID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookAuthorGenericError)
}
