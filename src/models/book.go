package models

import "github.com/google/uuid"

type Book struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string    `gorm:"type:varchar(255);not null;column:title"`
	Edition         uint8     `gorm:"type:smallint;column:edition"`
	PublicationYear uint      `gorm:"type:smallint;column:publication_year"`
}
