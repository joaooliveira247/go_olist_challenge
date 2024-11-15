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
