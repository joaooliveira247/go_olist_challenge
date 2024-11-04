package repositories_test

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/stretchr/testify/assert"
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

func TestCreateSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()
	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)
	expectedID := uuid.New()

	author := &models.Author{
		Name: "Luciano Ramalho",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors" ("name") VALUES ($1) RETURNING "id"`)).
		WithArgs(author.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	mock.ExpectCommit()

	id, err := repository.Create(author)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNotExpectedError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	author := &models.Author{
		Name: "Luciano Ramalho",
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors" ("name") VALUES ($1) RETURNING "id"`)).WithArgs(author.Name).WillReturnError(errors.New("some error not mapped"))
	mock.ExpectRollback()

	id, err := repository.Create(author)

	assert.Error(t, err, "some error not mapped")
	assert.Equal(t, uuid.Nil, id)
}

func TestCreateManySuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	authors := []models.Author{
		{Name: "J. K. Rowling"},
		{Name: "Stephen King"},
	}

	author1ID := uuid.New()
	author2ID := uuid.New()

	rows := mock.NewRows([]string{"id", "name"}).
		AddRow(author1ID, authors[0].Name).
		AddRow(author2ID, authors[1].Name)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name") VALUES ($1),($2) RETURNING "id"`),
	).WithArgs(authors[0].Name, authors[1].Name).WillReturnRows(rows)
	mock.ExpectCommit()

	ids, err := repository.CreateMany(&authors)

	assert.NoError(t, err)
	assert.Equal(t, ids[0], authors[0].ID)
	assert.Equal(t, ids[1], authors[1].ID)
	t.Logf("Author IDs: %v, %v", authors[0].ID, authors[1].ID)
}
