package repositories_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
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
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors" ("name") VALUES ($1) RETURNING "id"`)).WithArgs(author.Name).WillReturnError(&errors.AuthorGenericError)
	mock.ExpectRollback()

	id, err := repository.Create(author)

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.AuthorGenericError)
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

func TestCreateManyNotExpectedError(t *testing.T) {
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

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "authors" ("name") VALUES ($1),($2) RETURNING "id"`),
	).WithArgs(authors[0].Name, authors[1].Name).WillReturnError(errors.New("some error not mapped"))
	mock.ExpectCommit()

	ids, err := repository.CreateMany(&authors)

	assert.Nil(t, ids)
	assert.Error(t, err, "some error not mapped")
}

func TestGetAllSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	authors := []models.Author{
		{Name: "J. K. Rowling"},
		{Name: "Stephen King"},
		{Name: "Luciano Ramalho"},
	}

	rows := mock.NewRows([]string{"id", "name"}).
		AddRow(uuid.New(), authors[0].Name).
		AddRow(uuid.New(), authors[1].Name).
		AddRow(uuid.New(), authors[2].Name)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors"`)).WillReturnRows(rows)

	results, err := repository.GetAll()

	assert.Nil(t, err)
	assert.Len(t, results, 3)
}

func TestGetAllNotExpectedError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors"`)).WillReturnError(errors.New("some error not mapped"))

	results, err := repository.GetAll()

	assert.Nil(t, results)
	assert.Error(t, err, "some error not mapped")
}

func TestGetByIDSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	expectedID := uuid.New()

	row := mock.NewRows([]string{"id", "name"}).AddRow(expectedID, "Luciano Ramalho")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE id = $1 ORDER BY "authors"."id" LIMIT $2`)).WithArgs(expectedID, 1).WillReturnRows(row)

	result, err := repository.GetByID(expectedID)

	assert.Equal(t, "Luciano Ramalho", result.Name)
	assert.Equal(t, expectedID, result.ID)
	assert.Nil(t, err)
}

func TestGetByIDNotExpectedError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	repository := repositories.NewAuthorRepository(gormDB)

	expectedID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE id = $1 ORDER BY "authors"."id" LIMIT $2`)).WithArgs(expectedID, 1).WillReturnError(errors.New("some error not mapped"))

	result, err := repository.GetByID(expectedID)

	assert.Error(t, err, "some error not mapped")
	assert.Equal(t, models.Author{}, result)
}

func TestGetByNameSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	rows := mock.NewRows([]string{"id", "name"}).AddRow(uuid.New(), "Luciano Ramalho").AddRow(uuid.New(), "Luciano Peres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE name LIKE $1`)).WithArgs("%Luciano%").WillReturnRows(rows)

	repository := repositories.NewAuthorRepository(gormDB)

	result, err := repository.GetByName("Luciano")

	assert.Nil(t, err)
	assert.Len(t, result, 2)
}

func TestGetByNameNotExpectedError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "authors" WHERE name like $1"`)).WithArgs("%Luciano%").WillReturnError(errors.New("some error not mapped"))

	repository := repositories.NewAuthorRepository(gormDB)

	result, err := repository.GetByName("Luciano")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	expectedID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM "authors" WHERE id = $1`)).WithArgs(expectedID)
	mock.ExpectCommit()

	repository := repositories.NewAuthorRepository(gormDB)

	err := repository.Delete(expectedID)

	assert.Error(t, err)
}

func TestDeleteNotExpectedError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	expectedID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM "authors" WHERE id = $1`)).WithArgs(expectedID).WillReturnError(errors.New("some error not mapped"))
	mock.ExpectRollback()

	repository := repositories.NewAuthorRepository(gormDB)

	err := repository.Delete(expectedID)

	assert.Error(t, err)
}
