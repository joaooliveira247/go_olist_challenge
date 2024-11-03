package repositories

import (
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
)

type AuthorRepository interface {
	Create(author *models.Author) (uuid.UUID, error)
}