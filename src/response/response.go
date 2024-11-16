package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int
	Message    gin.H
}

var (
	InvalidRequestBody    = Response{http.StatusUnprocessableEntity, gin.H{"message": "request body invalid"}}
	InvalidID             = Response{http.StatusUnprocessableEntity, gin.H{"message": "invalid id"}}
	AuthorAlreadyExists   = Response{http.StatusConflict, gin.H{"message": "author already exists"}}
	UnableConnectDatabase = Response{http.StatusInternalServerError, gin.H{"message": "unable to connect to database"}}
	UnableCreateEntity    = Response{http.StatusInternalServerError, gin.H{"message": "unable to create entity"}}
	UnableFetchEntity     = Response{http.StatusInternalServerError, gin.H{"message": "unable to fetch entity"}}
)
