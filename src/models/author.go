package models

import (
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
)

type Author struct {
	ID   uuid.UUID `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string    `json:"name,omitempty" binding:"required,min=2" gorm:"type:varchar(255);column:name;unique;not null"`
}

func (author *Author) validate() error {
	if author.Name == "" {
		return &errors.FieldNameInvalid
	}
	return nil
}
