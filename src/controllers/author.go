package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/dto"
	custom "github.com/joaooliveira247/go_olist_challenge/src/errors"
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
		if errors.Is(err, &custom.AuthorAlreadyExists) {
			ctx.JSON(response.AuthorAlreadyExists.StatusCode, response.AuthorAlreadyExists.Message)
			return
		}
		ctx.JSON(response.UnableCreateEntity.StatusCode, response.UnableCreateEntity.Message)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
	return
}

func (ctrl *AuthorController) GetAllAuthors(ctx *gin.Context) {
	authors, err := ctrl.repository.GetAll()

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, authors)
	return
}

func (ctrl *AuthorController) GetAuthorByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(response.InvalidID.StatusCode, response.InvalidID.Message)
		return
	}

	author, err := ctrl.repository.GetByID(id)

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, author)
	return
}

func (ctrl *AuthorController) GetAuthorByName(ctx *gin.Context) {
	var queryParam dto.AuthorQueryParam

	if ctx.ShouldBindQuery(&queryParam) != nil {
		ctx.JSON(response.InvalidParam.StatusCode, response.InvalidParam.Message)
		return
	}

	authors, err := ctrl.repository.GetByName(queryParam.Name)

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, authors)
	return
}
