package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Book struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string    `json:"title" binding:"required,gt=2" gorm:"type:varchar(255);not null;column:title"`
	Edition         uint8     `json:"edition" binding:"required,gt=0" gorm:"type:smallint;column:edition"`
	PublicationYear uint      `json:"publication_year" binding:"required,gt=0" gorm:"type:smallint;column:publication_year"`
}

type BookIn struct {
	Book
	AuthorsID []uuid.UUID `json:"authors" binding:"gt=0,required,dive,uuid" gorm:"-"`
}

type BookOut struct {
	Book
	AuthorsName pq.StringArray `json:"authors" gorm:"type:text[];column:authors"`
}
