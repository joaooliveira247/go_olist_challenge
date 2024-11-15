package controllers

import (
	"github.com/joaooliveira247/go_olist_challenge/src/repositories"
)

type AuthorController struct {
	repository repositories.AuthorRepository
}

func NewAuthorController(repo repositories.AuthorRepository) *AuthorController {
	return &AuthorController{repo}
}
