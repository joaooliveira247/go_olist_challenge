package models

import "github.com/google/uuid"

type BookAuthor struct {
	BookID   uuid.UUID `gorm:"primaryKey;column:book_id"`
	Book     Book      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:BookID;"`
	AuthorID uuid.UUID `gorm:"primaryKey;column:author_id"`
	Author   Author    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AuthorID;"`
}

func (BookAuthor) TableName() string {
	return "book_author"
}
