package controllers

import "github.com/joaooliveira247/go_olist_challenge/src/repositories"

type BookController struct {
	repository repositories.BookRepository
}

func NewBookController(repo repositories.BookRepository) *BookController {
	return &BookController{repo}
}
