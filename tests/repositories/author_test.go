package repositories_test

import (
	"errors"
	"log"
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
	mock.ExpectQuery(`INSERT INTO "authors" \("name"\) VALUES \(\$1\) RETURNING "id"`).
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
	mock.ExpectQuery(`INSERT INTO "authors" \("name"\) VALUES \(\$1\) RETURNING "id"`).WithArgs(author.Name).WillReturnError(errors.New("some error not mapped"))
	mock.ExpectRollback()

	id, err := repository.Create(author)

	assert.Error(t, err, "some error not mapped")
	assert.Equal(t, uuid.Nil, id)
}
