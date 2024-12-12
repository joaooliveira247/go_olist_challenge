package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/dto"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
	"github.com/joaooliveira247/go_olist_challenge/src/response"
)

type BookController struct {
	bookRepository       repositories.BookRepository
	bookAuthorRepository repositories.BookAuthorRepository
}

func NewBookController(bookRepo repositories.BookRepository, bookAuthorRepo repositories.BookAuthorRepository) *BookController {
	return &BookController{bookRepo, bookAuthorRepo}
}

func (controller *BookController) Create(ctx *gin.Context) {
	var book models.BookIn

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(response.InvalidRequestBody.StatusCode, response.InvalidRequestBody.Message)
		fmt.Println(err)
		return
	}

	bookID, err := controller.bookRepository.Create(&book.Book)

	if err != nil {
		ctx.JSON(response.UnableCreateEntity.StatusCode, response.UnableCreateEntity.Message)
		return
	}

	for _, author := range book.AuthorsID {
		if err := controller.bookAuthorRepository.Create(&models.BookAuthor{BookID: bookID, AuthorID: author}); err != nil {
			ctx.JSON(response.UnableCreateEntity.StatusCode, response.UnableCreateEntity.Message)
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": bookID})
	return
}

func (controller *BookController) GetBooksByQuery(ctx *gin.Context) {
	var bookQuery dto.BookQueryParam

	if err := ctx.ShouldBindQuery(&bookQuery); err != nil {
		ctx.JSON(response.InvalidParam.StatusCode, response.InvalidParam.Message)
		return
	}

	queries, err := bookQuery.AsQuery()

	if err != nil {
		ctx.JSON(response.InvalidParam.StatusCode, response.InvalidParam.Message)
		return
	}

	if len(queries) > 0 {
		books, err := controller.bookRepository.GetBookByQuery(queries)
		if err != nil {
			ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
			return
		}
		ctx.JSON(http.StatusOK, books)
		return
	}

	books, err := controller.bookRepository.GetAll()

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, books)
	return
}

func (controller *BookController) GetBookByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil || id == uuid.Nil {
		ctx.JSON(response.InvalidID.StatusCode, response.InvalidID.Message)
		return
	}

	book, err := controller.bookRepository.GetBookByID(id)

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, book)
	return
}

func (controller *BookController) GetBookByAuthorID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Query("authorID"))

	if err != nil || id == uuid.Nil {
		ctx.JSON(response.InvalidID.StatusCode, response.InvalidID.Message)
		return
	}

	books, err := controller.bookRepository.GetBooksByAuthorID(id)

	if err != nil {
		ctx.JSON(response.UnableFetchEntity.StatusCode, response.UnableFetchEntity.Message)
		return
	}

	ctx.JSON(http.StatusOK, books)
	return
}
