package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joaooliveira247/go_olist_challenge/src/controllers"
	"github.com/joaooliveira247/go_olist_challenge/src/db"
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
)

func BookRoutes(eng *gin.Engine) {
	gormDB, err := db.GetDBConnection()

	if err != nil {
		log.Fatal(err)
		return
	}

	bookRepository := repositories.NewBookRepository(gormDB)
	bookAuthorRepository := repositories.NewBookAuthorRepository(gormDB)

	controller := controllers.NewBookController(bookRepository, bookAuthorRepository)

	bookGroup := eng.Group("/books")
	{
		bookGroup.POST("/", controller.Create)
		bookGroup.GET("/", controller.GetBooksByQuery)
		bookGroup.GET("/:id", controller.GetBookByID)
		bookGroup.GET("/:authorID", controller.GetBookByAuthorID)
		bookGroup.PUT("/:id", controller.UpdateBook)
		bookGroup.DELETE("/:id", controller.DeleteBook)
	}
}
