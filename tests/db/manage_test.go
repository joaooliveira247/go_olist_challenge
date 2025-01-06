package db_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateAllTablesSuccess(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	// Mock SELECT for "authors" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("authors", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Mock CREATE TABLE for "authors"
	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "authors" ("id" uuid DEFAULT gen_random_uuid(),"name" varchar(255) NOT NULL,PRIMARY KEY ("id"),CONSTRAINT "uni_authors_name" UNIQUE ("name"))`,
	)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock SELECT for "books" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Mock CREATE TABLE for "books"
	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "books" ("id" uuid DEFAULT gen_random_uuid(),"title" varchar(255) NOT NULL,"edition" smallint,"publication_year" smallint,PRIMARY KEY ("id"))`,
	)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock SELECT for "book_authors" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("book_authors", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Mock CREATE TABLE for "book_authors"
	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "book_authors" ("book_id" uuid,"author_id" uuid,PRIMARY KEY ("book_id","author_id"),CONSTRAINT "fk_book_authors_book" FOREIGN KEY ("book_id") REFERENCES "books"("id") ON DELETE CASCADE ON UPDATE CASCADE,CONSTRAINT "fk_book_authors_author" FOREIGN KEY ("author_id") REFERENCES "authors"("id") ON DELETE CASCADE ON UPDATE CASCADE)`,
	)).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.CreateTables(gormDB)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAllTablesReturnErrorWhenCreateAuthorsTableAlreadyExists(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	// Mock SELECT for "authors" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("authors", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := db.CreateTables(gormDB)

	assert.Error(t, err)
}

func TestCreateAllTablesReturnGenericErrorWhenCreateAuthors(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	// Mock SELECT for "authors" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("authors", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Mock CREATE TABLE for "authors"
	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "authors" ("id" uuid DEFAULT gen_random_uuid(),"name" varchar(255) NOT NULL,PRIMARY KEY ("id"),CONSTRAINT "uni_authors_name" UNIQUE ("name"))`,
	)).WillReturnError(&errors.AuthorGenericError)

	err := db.CreateTables(gormDB)

	assert.Error(t, err)
}

func TestCreateAllTablesReturnErrorWhenCreateBooksTableAlreadyExists(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := db.CreateTables(gormDB)

	assert.Error(t, err)
}

func TestCreateAllTablesReturnErrorGenericWhenCreateBooksTable(t *testing.T) {
	gormDB, mock := mocks.SetupMockDB()
	
	// Mock SELECT for "authors" table existence check
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("authors", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Mock CREATE TABLE for "authors"
	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "authors" ("id" uuid DEFAULT gen_random_uuid(),"name" varchar(255) NOT NULL,PRIMARY KEY ("id"),CONSTRAINT "uni_authors_name" UNIQUE ("name"))`,
	)).WillReturnResult(sqlmock.NewResult(1, 1))
	
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`,
	)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(regexp.QuoteMeta(
		`CREATE TABLE "books" ("id" uuid DEFAULT gen_random_uuid(),"title" varchar(255) NOT NULL,"edition" smallint,"publication_year" smallint,PRIMARY KEY ("id"))`,
	)).WillReturnError(&errors.BookGenericError)

	err := db.CreateTables(gormDB)

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}
