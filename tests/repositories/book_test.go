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

func TestCreateBookSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	book := mocks.NewMockBook()

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

	book := mocks.NewMockBook()

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

	book := mocks.NewMockBook()

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

func TestUpdateBookSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET "edition"=$1,"publication_year"=$2 WHERE id = $3`)).WithArgs(2, 2023, bookID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := repositories.NewBookRepository(gormDB)
	err := repository.Update(bookID, &models.Book{
		Edition:         2,
		PublicationYear: 2023,
	})

	assert.Nil(t, err)
}

func TestUpdateBookReturnNothingToUpdate(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET "edition"=$1,"publication_year"=$2 WHERE id = $3`)).WithArgs(2, 2023, bookID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repository := repositories.NewBookRepository(gormDB)
	err := repository.Update(bookID, &models.Book{
		Edition:         2,
		PublicationYear: 2023,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookNothingToUpdate)
}

func TestUpdateBookReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET "edition"=$1,"publication_year"=$2 WHERE id = $3`)).WithArgs(2, 2023, bookID).WillReturnError(&errors.BookGenericError)
	mock.ExpectRollback()

	repository := repositories.NewBookRepository(gormDB)
	err := repository.Update(bookID, &models.Book{
		Edition:         2,
		PublicationYear: 2023,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}

func TestGetAllBooksSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "title", "edition", "publication_year", "authors"})

	MBooks := mocks.NewMockBooks()

	for _, book := range MBooks {
		rows.AddRow(book.ID, book.Title, book.Edition, book.PublicationYear, book.AuthorsName)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select b.id, b.title, b.edition, b.publication_year, array_agg(a.name) as authors from book_author ba inner join books b on ba.book_id = b.id inner join authors a on ba.author_id = a.id group by b.id;`)).WillReturnRows(rows)

	repository := repositories.NewBookRepository(gormDB)
	books, err := repository.GetAll()

	assert.Nil(t, err)
	assert.Len(t, books, 4)
}

func TestGetAllReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`select b.id, b.title, b.edition, b.publication_year, array_agg(a.name) as authors from book_author ba inner join books b on ba.book_id = b.id inner join authors a on ba.author_id = a.id group by b.id;`)).WillReturnError(&errors.BookGenericError)

	repository := repositories.NewBookRepository(gormDB)

	books, err := repository.GetAll()

	assert.Nil(t, books)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}

func TestGetBookByOneQuerySuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "title", "edition", "publication_year", "authors"})

	MBooks := mocks.NewMockBooks()[2:]

	for _, book := range MBooks {
		rows.AddRow(book.ID, book.Title, book.Edition, book.PublicationYear, book.AuthorsName)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE b.title = Python Fluente GROUP BY b.id;`)).WillReturnRows(rows)

	repository := repositories.NewBookRepository(gormDB)
	books, err := repository.GetBookByQuery(map[string]interface{}{"title": "Python Fluente"})

	assert.Nil(t, err)
	assert.Equal(t, MBooks, books)
	assert.Len(t, books, 2)
}

func TestGetBookByQuerySuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	MBook := mocks.NewMockBookOut()

	rows := sqlmock.NewRows([]string{"id", "title", "edition", "publication_year", "authors"}).AddRow(MBook.ID, MBook.Title, MBook.Edition, MBook.PublicationYear, MBook.AuthorsName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE b.title = the Rust Programming Language AND b.edition = 1 AND b.publication_year = 2018 GROUP BY b.id;`)).WillReturnRows(rows)

	repository := repositories.NewBookRepository(gormDB)
	book, err := repository.GetBookByQuery(map[string]interface{}{"title": MBook.Title, "edition": MBook.Edition, "publication_year": MBook.PublicationYear})

	assert.Nil(t, err)
	assert.Len(t, book, 1)
	assert.Equal(t, MBook, book[0])
}

func TestGetBookByQueryReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE b.title = the Rust Programming Language AND b.edition = 1 AND b.publication_year = 2018 GROUP BY b.id;`)).WillReturnError(&errors.BookGenericError)

	repository := repositories.NewBookRepository(gormDB)
	book, err := repository.GetBookByQuery(map[string]interface{}{"title": "the Rust Programming Language", "edition": uint8(1), "publication_year": uint(2018)})

	assert.Nil(t, book)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}

func TestGetBookByIDSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	Mbook := mocks.NewMockBookOut()

	rows := sqlmock.NewRows([]string{"id", "title", "edition", "publication_year", "authors"}).AddRow(Mbook.Book.ID, Mbook.Book.Title, Mbook.Book.Edition, Mbook.Book.PublicationYear, Mbook.AuthorsName)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE ba.book_id = $1 GROUP BY b.id ORDER BY b.id LIMIT 1;`,
	)).WithArgs(Mbook.Book.ID).WillReturnRows(rows)

	repository := repositories.NewBookRepository(gormDB)

	book, err := repository.GetBookByID(Mbook.Book.ID)

	assert.Equal(t, Mbook, book)
	assert.Nil(t, err)
}

func TestGetBookByIDReturnBookNotFound(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE ba.book_id = $1 GROUP BY b.id ORDER BY b.id LIMIT 1;`)).WithArgs(bookID).WillReturnRows(sqlmock.NewRows([]string{}))

	repository := repositories.NewBookRepository(gormDB)
	book, err := repository.GetBookByID(bookID)

	assert.Empty(t, book)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookNotFound)
}

func TestGetBookByIDReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id INNER JOIN authors a ON ba.author_id = a.id WHERE ba.book_id = $1 GROUP BY b.id ORDER BY b.id LIMIT 1;`,
	)).WithArgs(bookID).WillReturnError(&errors.BookGenericError)

	repository := repositories.NewBookRepository(gormDB)

	book, err := repository.GetBookByID(bookID)

	assert.Empty(t, book)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}

func TestGetBooksByAuthorIDSuccess(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	authorID := uuid.New()

	MBooks := mocks.NewMockBooks()[2:]

	rows := sqlmock.NewRows([]string{"id", "title", "edition", "publication_year", "authors"})

	for _, book := range MBooks {
		rows.AddRow(book.Book.ID, book.Book.Title, book.Book.Edition, book.Book.PublicationYear, book.AuthorsName)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id
		INNER JOIN authors a ON ba.author_id = a.id WHERE ba.author_id = $1 GROUP BY b.id ORDER BY b.id;`)).WithArgs(authorID).WillReturnRows(rows)

	repository := repositories.NewBookRepository(gormDB)

	books, err := repository.GetBooksByAuthorID(authorID)

	assert.Len(t, books, 2)
	assert.Equal(t, MBooks, books)
	assert.Nil(t, err)
}

func TestGetBooksByAuthorIDReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	authorID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, b.title, b.edition, b.publication_year, array_agg(a.name) AS authors FROM book_author ba INNER JOIN books b ON ba.book_id = b.id
		INNER JOIN authors a ON ba.author_id = a.id WHERE ba.author_id = $1 GROUP BY b.id ORDER BY b.id;`)).WithArgs(authorID).WillReturnError(&errors.BookGenericError)

	repository := repositories.NewBookRepository(gormDB)

	books, err := repository.GetBooksByAuthorID(authorID)

	assert.Empty(t, books)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)

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

func TestDeleteBookReturnNotFound(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).WithArgs(bookID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repository := repositories.NewBookRepository(gormDB)
	err := repository.Delete(bookID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookNotFound)
}

func TestDeleteBookReturnGenericError(t *testing.T) {
	gormDB, mock := SetupMockDB()

	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	bookID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).WithArgs(bookID).WillReturnError(&errors.BookGenericError)
	mock.ExpectRollback()

	repository := repositories.NewBookRepository(gormDB)

	err := repository.Delete(bookID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, &errors.BookGenericError)
}
