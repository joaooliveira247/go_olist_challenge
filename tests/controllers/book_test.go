package controllers_test

import (
	"bytes"
	"encoding/json"
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

func TestBookGetByQueryReturnAllSucess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	Mbooks := mocks.NewMockBooks()

	mockBookRepository.On("GetAll").Return(Mbooks, nil)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/books/", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GetBooksByQuery(c)

	byteMbooks, _ := json.Marshal(Mbooks)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(byteMbooks), w.Body.String())
}

func TestBookGetByManyQueriesReturnSucess(t *testing.T) {
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)
	mockBookRepository := new(mocks.BookRepository)

	Mbooks := mocks.NewMockBooks()

	testCases := []struct {
		name       string
		url        string
		query      map[string]interface{}
		mockResult []models.BookOut
	}{
		{
			"Title Param Return Two Books",
			"/books/?title=Python Fluente",
			map[string]interface{}{"title": "Python Fluente"},
			Mbooks[:2],
		},
		{
			"Edition Param Return Two Books",
			"/books/?edition=1",
			map[string]interface{}{"edition": float64(1)},
			[]models.BookOut{Mbooks[1], Mbooks[3]},
		},
		{
			"publicationYear Param Return Two Books",
			"/books/?publicationYear=2015",
			map[string]interface{}{
				"publication_year": float64(2015),
			},
			[]models.BookOut{Mbooks[1], Mbooks[2]},
		},
		{
			"Edition and PublicatioYear Return One Book",
			"/books/?edition=2&publicationYear=2015",
			map[string]interface{}{"edition": float64(2), "publication_year": float64(2015)},
			[]models.BookOut{Mbooks[1]},
		},
		{
			"Title, Edition and PublicationYear Return One Book",
			"/books/?title=The Go Programming Language&edition=2&publicationYear=2015",
			map[string]interface{}{"title": "The Go Programming Language", "edition": float64(2), "publication_year": float64(2015)},
			[]models.BookOut{Mbooks[1]},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockBookRepository.ExpectedCalls = nil
			mockBookRepository.On("GetBookByQuery", testCase.query).Return(testCase.mockResult, nil)

			controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, testCase.url, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			controller.GetBooksByQuery(c)

			byteMbooks, _ := json.Marshal(testCase.mockResult)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, string(byteMbooks), w.Body.String())
		})
	}
}

func TestBookGetByQueryReturnInvalidParam(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	testCases := []struct {
		name string
		url  string
	}{
		{
			"Edition greather than uint8",
			"/books/?edition=512",
		},
		{
			"Edition Lower than one",
			"/books/?edition=-10",
		},
		{
			"Edition with letter",
			"/books/?edition=a",
		},
		{
			"Edition with special character",
			"/books/?edition=@",
		},
		{
			"PublicationYear Lower than one",
			"/books/?publicationYear=-1",
		},
		{
			"PublicationYear with special character",
			"/books/?publicationYear=@",
		},
	}

	for _, testCase := range testCases {
		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, testCase.url, nil)
		c.Request.Header.Set("Content-Type", "application/json")

		controller.GetBooksByQuery(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message": "invalid query param"}`, w.Body.String())
	}
}

func TestBookGetByQueryReturnUnableFetchEntity(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	testCases := []struct {
		name       string
		methodName string
		args       interface{}
		error      error
		url        string
	}{
		{
			"Empty query but error in GetAll",
			"GetAll",
			nil,
			&errors.BookGenericError,
			"/books/",
		},
		{
			"Method has query but error in GetBookByQuery",
			"GetBookByQuery",
			map[string]interface{}{"title": "Python Fluente"},
			&errors.BookGenericError,
			"/books/?title=Python Fluente",
		},
	}

	for _, testCase := range testCases {
		mockBookRepository.Calls = nil

		if testCase.args != nil {
			mockBookRepository.On(testCase.methodName, testCase.args).Return(nil, testCase.error)
		} else {
			mockBookRepository.On(testCase.methodName).Return(nil, testCase.error)
		}

		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, testCase.url, nil)
		c.Request.Header.Set("Content-Type", "application/json")

		controller.GetBooksByQuery(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
	}
}

func TestGetBookByIDSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	Mbook := mocks.NewMockBookOut()

	mockBookRepository.On("GetBookByID", Mbook.ID).Return(Mbook, nil)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/books/%s", Mbook.ID), nil)
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%s", Mbook.ID)}}
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GetBookByID(c)

	expectedBook, _ := json.Marshal(Mbook)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expectedBook), w.Body.String())
}
