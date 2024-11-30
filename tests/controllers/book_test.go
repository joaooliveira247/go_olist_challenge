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
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestBookCreateReturnInvalidRequestBody(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	requestBodyTests := []struct {
		name        string
		requestBody string
	}{
		{
			"Title Empty",
			`{
				"edition": 1,
				"publication_year": 2018,
				"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
}`,
		},
		{
			"Title Less Than 2",
			`{
			"title": "Th",
			"edition": 1,
			"publication_year": 2018,
			"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
}`,
		},
		{
			"Edition Less Than 1",
			`{
			"title": "The Rust Programming Language",
			"edition": -1,
			"publication_year": 2018,
			"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
}`,
		},
		{
			"Publication Year Less Than 1",
			`{
			"title": "The Rust Programming Language",
			"edition": 1,
			"publication_year": -1,
			"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
}`,
		},
		{
			"One Author ID invalid",
			`{
			"title": "The Rust Programming Language",
			"edition": -1,
			"publication_year": 2018,
			"authors": ["abc", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
}`,
		},
		{
			"Both Authors ID invalid",
			`{
			"title": "The Rust Programming Language",
			"edition": -1,
			"publication_year": 2018,
			"authors": ["abc", "def"]
}`,
		},
		{
			"Valid ID But Empty",
			`{
			"title": "The Rust Programming Language",
			"edition": -1,
			"publication_year": 2018,
			"authors": ["00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000"]
}`,
		},
	}

	for _, testCase := range requestBodyTests {
		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/books/", bytes.NewBufferString(testCase.requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Create(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"message": "request body invalid"}`, w.Body.String())
	}
}

func TestBookCreateReturnUnableCreateEntity(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	mockBookRepository.On("Create", mock.Anything).Return(uuid.Nil, &errors.BookGenericError)

	body := `{
		"title": "The Rust Programming Language",
		"edition": 1,
		"publication_year": 2018,
		"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
	}`

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/books/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to create entity"}`, w.Body.String())
}

func TestBookCreateReturnUnableCreateEntityWhenCreateRelationship(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	mockBookRepository.On("Create", mock.Anything).Return(uuid.New(), nil)
	mockBookAuthorRepository.On("Create", mock.Anything).Return(&errors.BookAuthorGenericError)

	body := `{
		"title": "The Rust Programming Language",
		"edition": 1,
		"publication_year": 2018,
		"authors": ["4ed37603-c983-4137-bbe9-bccfc30b53a6", "31455548-62a9-4935-aa89-c1d2ac036e0f"]
	}`

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/books/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to create entity"}`, w.Body.String())
}
