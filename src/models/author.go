package models

import (
	"github.com/google/uuid"
)

type Author struct {
	ID   uuid.UUID `json:"id,omitempty" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string    `json:"name,omitempty" binding:"required,min=2" gorm:"type:varchar(255);column:name;unique;not null"`
}
