package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveira247/go_olist_challenge/src/controllers"
	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
)

func AuthorRoutes(eng *gin.Engine) {
	gormDB, err := db.GetDBConnection()

	if err != nil {
		log.Fatal("DATABASE: ", err)
	}

	authorRepository := repositories.NewAuthorRepository(gormDB)

	controller := controllers.NewAuthorController(authorRepository)

	authorRouter := eng.Group("/authors")
	{
		authorRouter.POST("/", controller.CreateAuthor)
		authorRouter.GET("/", controller.GetAuthors)
		authorRouter.DELETE("/:id", controller.DeleteAuthor)
	}
}
