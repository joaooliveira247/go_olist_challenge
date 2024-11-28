package controllers

import (
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
)

type BookController struct {
	bookRepository repositories.BookRepository
	bookAuthorRepository repositories.BookAuthorRepository
}
