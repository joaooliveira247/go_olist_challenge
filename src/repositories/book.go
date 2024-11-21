package repositories

import "gorm.io/gorm"

type BookRepository interface {
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}
