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

func TestGetAuthorByIDReturnIvalidID(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/%s", "invalid id"), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: "invalid id"},
	}

	controller := controllers.NewAuthorController(mockRepository)
	controller.GetAuthorByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message": "invalid id"}`, w.Body.String())
}

func TestGetAuthorByIDReturnUnableFetchEntity(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	expectedID := uuid.New()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	mockRepository.On("GetByID", expectedID).Return(models.Author{}, &errors.AuthorGenericError)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/%s", expectedID), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "id", Value: fmt.Sprintf("%s", expectedID)},
	}

	controller := controllers.NewAuthorController(mockRepository)
	controller.GetAuthorByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
}

func TestGetAuthorByNameSucess(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	authors := []models.Author{
		{
			ID:   uuid.New(),
			Name: "Oswald de Andrade",
		},
		{
			ID:   uuid.New(),
			Name: "Mário de Andrade",
		},
	}

	mockRepository.On("GetByName", "de Andrade").Return(authors, nil)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/?name=de%20Andrade", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.Query().Add("name", "de Andrade")

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAuthorByName(c)

	expectedJSON, _ := json.Marshal(authors)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())
}

func TestGetAuthorByNameEmptyParamReturnInvalidQuery(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/?name=", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.Query().Add("name", "")

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAuthorByName(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message": "invalid query param"}`, w.Body.String())
}

func TestGetAuthorByNameLTOneParamReturnInvalidQuery(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/?name=a", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.Query().Add("name", "a")

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAuthorByName(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message": "invalid query param"}`, w.Body.String())
}

func TestGetAuthorByNameReturnUnableFetchEntity(t *testing.T) {
	mockRepository := new(mocks.AuthorRepository)

	mockRepository.On("GetByName", "de Andrade").Return(nil, &errors.AuthorGenericError)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/authors/?name=de%20Andrade", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.Query().Add("name", "de Andrade")

	controller := controllers.NewAuthorController(mockRepository)

	controller.GetAuthorByName(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"message": "unable to fetch entity"}`, w.Body.String())
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
