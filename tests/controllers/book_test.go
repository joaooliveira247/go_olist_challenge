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

func TestGetBooksReturnInvalidParam(t *testing.T) {
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

		controller.GetBooks(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message": "invalid query param"}`, w.Body.String())
	}
}

func TestGetBooksReturnInvalidID(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	params := []string{"bookID", "authorID"}

	testCases := []struct {
		name string
		url  string
	}{
		{
			"invalid",
			"/books/?%s=123",
		},
		{
			"empty",
			"/books/?%s=00000000-0000-0000-0000-000000000000",
		},
		{
			"bookID length OK but uuid invalid",
			"/books/?%s=@@@@@@@@-@@@@-@@@@-@@@@-@@@@@@@@@@@@",
		},
	}

	for _, param := range params {
		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%s %s", param, testCase.name), func(t *testing.T) {
				controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

				w := httptest.NewRecorder()
				gin.SetMode(gin.TestMode)

				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf(testCase.url, param), nil)
				c.Request.Header.Set("Content-Type", "application/json")

				controller.GetBooks(c)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
			})
		}
	}
}

func TestGetBooksReturnUnableFetchEntity(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	testCases := []struct {
		name             string
		url              string
		methodRepository string
		returnObj        interface{}
	}{
		{
			"bookID",
			fmt.Sprintf("/books/?bookID=%s", uuid.New().String()),
			"GetBookByID",
			models.BookOut{},
		},
		{
			"AuthorID",
			fmt.Sprintf("/books/?authorID=%s", uuid.New().String()),
			"GetBooksByAuthorID",
			nil,
		},
		{
			"Query",
			"/books/?title=The Go Programming Language&edition=2&publicationYear=2015",
			"GetBookByQuery",
			nil,
		},
		{
			"All",
			"/books/",
			"GetAll",
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockBookRepository.Calls = nil
			mockBookRepository.On(testCase.methodRepository, mock.Anything).Return(testCase.returnObj, &errors.BookGenericError)
			controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, testCase.url, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			controller.GetBooks(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
		})
	}
}

func TestGetBooksQueryBookIDSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	mbook := mocks.NewMockBookOut()

	mockBookRepository.On("GetBookByID", mbook.ID).Return(mbook, nil)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/books/?bookID=%s", mbook.ID.String()), nil)
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GetBooks(c)

	bjson, _ := json.Marshal(mbook)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(bjson), w.Body.String())

}

func TestGetBooksQueryAuhthorIDSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	authorID := uuid.New()
	mbook := mocks.NewMockBooks()[:2]

	mockBookRepository.On("GetBooksByAuthorID", authorID).Return(mbook, nil)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/books/?authorID=%s", authorID.String()), nil)
	c.Request.Header.Set("Content-Type", "application/json")

	controller.GetBooks(c)

	bjson, _ := json.Marshal(mbook)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(bjson), w.Body.String())

}

func TestGetBooksManyQueriesSuccess(t *testing.T) {
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)
	mockBookRepository := new(mocks.BookRepository)

	Mbooks := mocks.NewMockBooks()

	testCases := []struct {
		name       string
		url        string
		query      string
		mockResult []models.BookOut
	}{
		{
			"Title Param Return Two Books",
			"/books/?title=Python Fluente",
			"b.title = 'Python Fluente'",
			Mbooks[:2],
		},
		{
			"Edition Param Return Two Books",
			"/books/?edition=1",
			"b.edition = 1",
			[]models.BookOut{Mbooks[1], Mbooks[3]},
		},
		{
			"publicationYear Param Return Two Books",
			"/books/?publicationYear=2015",
			"b.publication_year = 2015",
			[]models.BookOut{Mbooks[1], Mbooks[2]},
		},
		{
			"Edition and PublicatioYear Return One Book",
			"/books/?edition=2&publicationYear=2015",
			"b.edition = 2 AND b.publication_year = 2015",
			[]models.BookOut{Mbooks[1]},
		},
		{
			"Title, Edition and PublicationYear Return One Book",
			"/books/?title=The Go Programming Language&edition=2&publicationYear=2015",
			"b.title = 'The Go Programming Language' AND b.edition = 2 AND b.publication_year = 2015",
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

			controller.GetBooks(c)

			byteMbooks, _ := json.Marshal(testCase.mockResult)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, string(byteMbooks), w.Body.String())
		})
	}
}

func TestGetBooksAllSuccess(t *testing.T) {
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

	controller.GetBooks(c)

	byteMbooks, _ := json.Marshal(Mbooks)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(byteMbooks), w.Body.String())
}

func TestUpdateBookInfoSucess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	testCases := []struct {
		name  string
		model *models.BookUpdate
		body  string
	}{
		{
			"Update Title",
			&models.BookUpdate{
				BookInfo: models.BookInfo{
					Title: "Python Fluente",
				},
			},
			`{"title": "Python Fluente"}`,
		},
		{
			"Update Edition",
			&models.BookUpdate{
				BookInfo: models.BookInfo{
					Edition: 2,
				},
			},
			`{"edition": 2}`,
		},
		{
			"Update Publication Year",
			&models.BookUpdate{
				BookInfo: models.BookInfo{
					PublicationYear: 2023,
				},
			},
			`{"publication_year": 2023}`,
		},
		{
			"Full Update",
			&models.BookUpdate{
				BookInfo: models.BookInfo{
					Title:           "Python Fluente",
					Edition:         2,
					PublicationYear: 2023,
				},
			},
			`{"title": "Python Fluente", "edition": 2, "publication_year": 2023}`,
		},
	}

	for _, testCase := range testCases {
		mockBookRepository.On("Update", bookID, testCase.model).Return(nil)

		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID), bytes.NewBufferString(testCase.body))
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%s", bookID)}}
		c.Request.Header.Set("Content-Type", "application/json")

		controller.UpdateBook(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	}
}

func TestUpdateBookSingleAuthorIDSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	authors := []uuid.UUID{
		uuid.New(),
	}

	mockBookAuthorRepository.On("Delete", bookID).Return(nil)
	mockBookAuthorRepository.On("Create", &models.BookAuthor{
		BookID: bookID, AuthorID: authors[0],
	}).Return(nil)

	body := fmt.Sprintf(`{"authors": ["%s"]}`, authors[0])

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID), bytes.NewBufferString(body))
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.UpdateBook(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestUpdateBookDoubleAuthorIDSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	authors := []uuid.UUID{
		uuid.New(), uuid.New(),
	}

	mockBookAuthorRepository.On("Delete", bookID).Return(nil)
	for _, author := range authors {
		mockBookAuthorRepository.On("Create", &models.BookAuthor{
			BookID: bookID, AuthorID: author,
		}).Return(nil)
	}

	body := fmt.Sprintf(`{"authors": ["%s", "%s"]}`, authors[0], authors[1])

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID), bytes.NewBufferString(body))
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.UpdateBook(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestUpdateBookFullSucess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	MUpdate := mocks.NewMockUpdateBook()

	mockBookRepository.On("Update", bookID, &MUpdate).Return(nil)
	mockBookAuthorRepository.On("Delete", bookID).Return(nil)
	for _, author := range MUpdate.AuthorsID {
		mockBookAuthorRepository.On("Create", &models.BookAuthor{
			BookID: bookID, AuthorID: author,
		}).Return(nil)
	}

	body := fmt.Sprintf(`{"title": "%s","edition": %d,"publication_year": %d,"authors": ["%s", "%s"]}`, MUpdate.BookInfo.Title, MUpdate.BookInfo.Edition, MUpdate.BookInfo.PublicationYear, MUpdate.AuthorsID[0], MUpdate.AuthorsID[1])

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID), bytes.NewBufferString(body))
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.UpdateBook(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestUpdateBookReturnInvalidID(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	testCases := []struct {
		name string
		url  string
		id   string
	}{
		{
			"invalid uuid",
			"/book/123",
			"123",
		},
		{
			"empty uuid",
			fmt.Sprintf(`/book/%s`, uuid.Nil),
			uuid.Nil.String(),
		},
	}

	for _, testCase := range testCases {
		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, testCase.url, nil)
		c.Params = gin.Params{
			{Key: "id", Value: testCase.id},
		}
		c.Header("Content-Type", "application/json")

		controller.UpdateBook(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
	}
}

func TestUpdateBookReturnInvalidParam(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	testCases := []struct {
		name string
		body string
	}{
		{
			"without title",
			`{"title": }`,
		},
		{
			"edition lt 0",
			`{"edition": -1}`,
		},
		{
			"edition gt uint8",
			`{"edition": 256}`,
		},
		{
			"without edition",
			`{"edition": }`,
		},
		{
			"publication year lt 0",
			`{"publication_year": -1}`,
		},
		{
			"without publication_year",
			`{"publication_year": }`,
		},
		{
			"one uuid invalid",
			`{"authors": [123]}`,
		},
		{
			"uuid nil",
			fmt.Sprintf(`{"authors": [%s]}`, uuid.Nil.String()),
		},
		{
			"two uuid but one invalid",
			fmt.Sprintf(`{"authors": [%d, %s]}`, 123, uuid.New().String()),
		},
		{
			"two uuid but one nil",
			fmt.Sprintf(`{"authors": [%s, %s]}`, uuid.Nil.String(), uuid.New().String()),
		},
		{
			"authors empty",
			`{"authors": }`,
		},
	}

	bookID := uuid.New()

	for _, testCase := range testCases {
		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID.String()), bytes.NewBufferString(testCase.body))
		c.Params = gin.Params{
			{Key: "id", Value: bookID.String()},
		}
		c.Request.Header.Set("Content-Type", "application/json")

		controller.UpdateBook(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"message": "request body invalid"}`, w.Body.String())
	}
}

func TestUpdateBookInfoReturnError(t *testing.T) {
	testCases := []struct {
		name            string
		errorReturn     error
		body            string
		ExpectedCode    int
		ExpectedMessage string
	}{
		{
			"update return nothing to update",
			&errors.BookNothingToUpdate,
			`{"edition": 1}`,
			http.StatusNotModified,
			"",
		},
		{
			"update return generic error",
			&errors.BookGenericError,
			`{"edition": 1}`,
			http.StatusInternalServerError,
			`{"message": "unable to fetch entity"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockBookRepository := new(mocks.BookRepository)
			mockBookAuthorRepository := new(mocks.BookAuthorRepository)

			bookID := uuid.New()

			mockBookRepository.On("Update", mock.Anything, mock.Anything).Return(testCase.errorReturn)

			controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", bookID.String()), bytes.NewBufferString(testCase.body))
			c.Params = gin.Params{
				{Key: "id", Value: bookID.String()},
			}
			c.Request.Header.Set("Content-Type", "application/json")

			controller.UpdateBook(c)

			assert.Equal(t, testCase.ExpectedCode, w.Code)
			if testCase.ExpectedMessage != "" {
				assert.JSONEq(t, testCase.ExpectedMessage, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func TestUpdateBookWhenDeleteRefReturnError(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	body := fmt.Sprintf(`{"authors": ["%s"]}`, bookID.String())

	mockBookAuthorRepository.On("Delete", bookID).Return(&errors.BookAuthorGenericError)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf(`/books/%s`, bookID.String()), bytes.NewBufferString(body))
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.UpdateBook(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}

func TestUpdateBookWhenCreateRefReturnError(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	body := fmt.Sprintf(`{"authors": ["%s"]}`, bookID.String())

	mockBookAuthorRepository.On("Delete", mock.Anything).Return(nil)
	mockBookAuthorRepository.On("Create", mock.Anything).Return(&errors.BookAuthorGenericError)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf(`/books/%s`, bookID.String()), bytes.NewBufferString(body))
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.UpdateBook(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}

func TestDeleteBookReturnInvalidID(t *testing.T) {

	testCases := []struct {
		name string
		id   string
	}{
		{
			"invalid id",
			"123",
		},
		{
			"empty id",
			uuid.Nil.String(),
		},
	}

	for _, testCase := range testCases {
		mockBookRepository := new(mocks.BookRepository)
		mockBookAuthorRepository := new(mocks.BookAuthorRepository)

		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf(`/books/%s`, testCase.id), nil)
		c.Params = gin.Params{
			{Key: "id", Value: testCase.id},
		}
		c.Header("Content-Tyoe", "application/json")

		controller.DeleteBook(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
	}
}

func TestDeleteBookReturnError(t *testing.T) {
	testCases := []struct {
		name        string
		message     string
		returnError error
		status      int
	}{
		{
			"return NotModified",
			"",
			&errors.BookNotFound,
			http.StatusNotModified,
		},
		{
			"return InternalServerError",
			`{"message":"unable to fetch entity"}`,
			&errors.BookGenericError,
			http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		bookID := uuid.New()

		mockBookRepository := new(mocks.BookRepository)
		mockBookAuthorRepository := new(mocks.BookAuthorRepository)

		mockBookRepository.On("Delete", bookID).Return(testCase.returnError)

		controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)

		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf(`/books/%s`, bookID.String()), nil)
		c.Params = gin.Params{
			{Key: "id", Value: bookID.String()},
		}
		c.Header("Content-Tyoe", "application/json")

		controller.DeleteBook(c)

		assert.Equal(t, testCase.status, w.Code)
		assert.Equal(t, testCase.message, w.Body.String())
	}
}

func TestDeleteBookSuccess(t *testing.T) {
	mockBookRepository := new(mocks.BookRepository)
	mockBookAuthorRepository := new(mocks.BookAuthorRepository)

	bookID := uuid.New()

	mockBookRepository.On("Delete", bookID).Return(nil)

	controller := controllers.NewBookController(mockBookRepository, mockBookAuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf(`/books/%s`, bookID.String()), nil)
	c.Params = gin.Params{
		{Key: "id", Value: bookID.String()},
	}
	c.Header("Content-Type", "application/json")

	controller.DeleteBook(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
