package routes

import "github.com/gin-gonic/gin"

func RegistryRoutes(eng *gin.Engine) {
	AuthorRoutes(eng)
	BookRoutes(eng)
}
