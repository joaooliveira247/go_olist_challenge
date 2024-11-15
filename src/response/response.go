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
	InvalidRequestBody = Response{http.StatusUnprocessableEntity, gin.H{"message": "request body invalid"}}
)
