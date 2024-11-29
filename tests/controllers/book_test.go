package controllers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/controllers"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookCreateSucess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()
	authorsID := []uuid.UUID{uuid.New(), uuid.New()}

	MBookIn := models.BookIn{
		Book: models.Book{
			Title:           "The Rust Programming Language",
			Edition:         1,
			PublicationYear: 2018,
		},
		AuthorsID: authorsID,
	}

	mockBookRepository.On("Create", &MBookIn.Book).Return(bookID, nil)

	for _, author := range MBookIn.AuthorsID {
		mockBookAuthorRepository.On("Create", &models.BookAuthor{BookID: bookID, AuthorID: author}).Return(nil)
	}

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)

	body := fmt.Sprintf(`{
		"title": "The Rust Programming Language",
		"edition": 1,
		"publication_year": 2018,
		"authors": ["%s", "%s"]
}`, authorsID[0], authorsID[1])

	c.Request, _ = http.NewRequest(http.MethodPost, "/books/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`{"id": "%s"}`, bookID), w.Body.String())
}
