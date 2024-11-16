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

func TestGetAllAuthorsSuccess(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	authors := []models.Author{
		{
			ID:   uuid.New(),
			Name: "Luciano Ramalho",
		},
		{
			ID:   uuid.New(),
			Name: "J. K. Rowling",
		},
		{
			ID:   uuid.New(),
			Name: "Machado de Assis",
		},
	}

	mockRepository.On("GetAll").Return(authors, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAllAuthors(c)

	expectedJson, _ := json.Marshal(authors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expectedJson), w.Body.String())
}

func TestGetAllAuthorsReturnUnableFetchEntity(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	mockRepository.On("GetAll").Return(nil, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/", nil)

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAllAuthors(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}

func TestGetAuthorByIDSucess(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	expectedID := uuid.New()

	author := models.Author{
		ID:   expectedID,
		Name: "Machado de Assis",
	}

	mockRepository.On("GetByID", expectedID).Return(author, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/%s", expectedID), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: fmt.Sprintf("%s", expectedID)},
	}

	controlller := controllers.NewAuthorController(mockRepository)
	controlller.GetAuthorByID(c)

	expectedJSON, _ := json.Marshal(author)

	t.Log(c.Request.URL, expectedID)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())
}
