package models

import (
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/errors"
)

type Author struct {
	ID   uuid.UUID `json:"-" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string    `json:"name,omitempty" gorm:"type:varchar(255);column:name;unique;not null"`
}

func (author *Author) validate() error {
	if author.Name == "" {
		return &errors.FieldNameInvalid
	}
	return nil
}

func (author *Author) Prepare() error {
	if err := author.validate(); err != nil {
		return err
	}
	return nil
}
