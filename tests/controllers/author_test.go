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
)

func TestCreateSuccess(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)
	expectedID := uuid.New()
	author := models.Author{
		Name: "Luciano Ramalho",
	}
	mockRepository.On("Create", &author).Return(expectedID, nil)

	controller := controllers.NewAuthorController(mockRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	body := `{"name": "Luciano Ramalho"}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/authors/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateAuthor(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`{"id": "%s"}`, expectedID), w.Body.String())
}

func TestCreateReturnAlreadyExists(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)
	author := models.Author{
		Name: "Luciano Ramalho",
	}
	mockRepository.On("Create", &author).Return(uuid.UUID{}, &errors.AuthorAlreadyExists)

	controller := controllers.NewAuthorController(mockRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	body := `{"name": "Luciano Ramalho"}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/authors/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateAuthor(c)
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.JSONEq(t, `{"message": "author already exists"}`, w.Body.String())
}

func TestCreateReturnInvalidRequestBody(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)
	controller := controllers.NewAuthorController(mockRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	body := `{"name": "L"}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/authors/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateAuthor(c)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `{"message": "request body invalid"}`, w.Body.String())
}

func TestCreateReturnUnableCreateEntity(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	author := models.Author{
		Name: "Luciano Ramalho",
	}

	mockRepository.On("Create", &author).Return(uuid.UUID{}, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)

	body := `{"name": "Luciano Ramalho"}`

	c.Request, _ = http.NewRequest(http.MethodPost, "/authors/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockRepository)

	controller.CreateAuthor(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to create entity"}`, w.Body.String())
}

func TestGetAuthorsByQueryReturnErrorInAuthorID(t *testing.T) {
	mockAuthorRepository := new(mocks.AuthorRepository)

	testCases := []struct {
		name string
		url  string
	}{
		{
			"Invalid numeric id",
			"/authors/?authorID=123",
		},
		{
			"Invalid special character id",
			"/authors/?authorID=@_ç",
		},
		{
			"Invalid letter id",
			"/authors/?authorID=authorabc",
		},
		{
			"Empty uuid",
			fmt.Sprintf("/authors/?authorID=%s", uuid.Nil.String()),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, testCase.url, nil)
			c.Header("Content-Type", "application/json")

			controller := controllers.NewAuthorController(mockAuthorRepository)
			controller.GetAuthors(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
		})
	}
}

func TestGetAuthorsByIDReturnNotFound(t *testing.T) {
	authorID := uuid.New()

	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetByID", authorID).Return(models.Author{}, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/?authorID=%s", authorID), nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message": "author not found"}`, w.Body.String())
}

func TestGetAuthorsByNameReturnNotFound(t *testing.T) {
	authorName := "Edgar Allan Poe"
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetByName", authorName).Return(nil, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/?name=%s", authorName), nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message": "author not found"}`, w.Body.String())
}

func TestGetAuthorsByIDSuccess(t *testing.T) {
	authorID := uuid.New()

	mockAuthor := models.Author{
		ID:   authorID,
		Name: "Edgar Allan Poe",
	}

	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetByID", authorID).Return(mockAuthor, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/?authorID=%s", authorID), nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	bMock, _ := json.Marshal(mockAuthor)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(bMock), w.Body.String())
}

func TestGetAuthorsByNameSuccess(t *testing.T) {
	authorName := "Stephen"

	mockAuthors := []models.Author{
		{
			ID:   uuid.New(),
			Name: "Stephen King",
		},
		{
			ID:   uuid.New(),
			Name: "Stephen Hawking",
		},
	}

	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetByName", authorName).Return(mockAuthors, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/?name=%s", authorName), nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	bMock, _ := json.Marshal(mockAuthors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(bMock), w.Body.String())
}

func TestGetAuthorsReturnUnableFetchEntity(t *testing.T) {
	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetAll").Return(nil, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/", nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}

func TestGetAuthorsSuccess(t *testing.T) {
	mockAuthors := []models.Author{
		{
			ID:   uuid.New(),
			Name: "Edgar Allan Poe",
		},
		{
			ID:   uuid.New(),
			Name: "Stephen King",
		},
		{
			ID:   uuid.New(),
			Name: "Stephen Hawking",
		},
	}

	mockAuthorRepository := new(mocks.AuthorRepository)
	mockAuthorRepository.On("GetAll").Return(mockAuthors, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/", nil)
	c.Header("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockAuthorRepository)
	controller.GetAuthors(c)

	bMock, _ := json.Marshal(mockAuthors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(bMock), w.Body.String())

}

func TestDeleteAuthorSuccess(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	expectedID := uuid.New()

	mockRepository.On("Delete", expectedID).Return(nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/authors/%s", expectedID), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: fmt.Sprintf("%s", expectedID)},
	}

	controller := controllers.NewAuthorController(mockRepository)

	controller.DeleteAuthor(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeleteAuthorReturnInvalidID(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/authors/56", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: "56"},
	}

	controller := controllers.NewAuthorController(mockRepository)

	controller.DeleteAuthor(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
}

func TestDeleteAuthorReturnAuthorNotFound(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	expectedID := uuid.New()

	mockRepository.On("Delete", expectedID).Return(&errors.AuthorNotFound)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/authors/%s", expectedID), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: fmt.Sprintf("%s", expectedID)},
	}

	controller := controllers.NewAuthorController(mockRepository)
	controller.DeleteAuthor(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"message": "author not found"}`, w.Body.String())
}

func TestDeleteAuthorReturnUnableFetchEntity(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	expectedID := uuid.New()

	mockRepository.On("Delete", expectedID).Return(&errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/authors/%s", expectedID), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: fmt.Sprintf("%s", expectedID)},
	}

	controller := controllers.NewAuthorController(mockRepository)
	controller.DeleteAuthor(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}
