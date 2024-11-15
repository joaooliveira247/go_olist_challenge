package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int
	Message    gin.H
}
