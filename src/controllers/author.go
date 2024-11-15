package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/joaooliveira247/go_olist_challenge/src/response"
)

type AuthorController struct {
	repository repositories.AuthorRepository
}

func NewAuthorController(repo repositories.AuthorRepository) *AuthorController {
	return &AuthorController{repo}
}

func (ctrl *AuthorController) CreateAuthor(ctx *gin.Context) {
	var author models.Author

	if err := ctx.ShouldBindJSON(&author); err != nil {
		ctx.JSON(response.InvalidRequestBody.StatusCode, response.InvalidRequestBody.Message)
		return
	}

	id, err := ctrl.repository.Create(&author)

	if err != nil {
		ctx.JSON(response.UnableCreateEntity.StatusCode, response.UnableCreateEntity.Message)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
	return
}
