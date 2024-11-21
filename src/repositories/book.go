package repositories

import "gorm.io/gorm"

type BookRepository interface {
}

type bookRepository struct {
	db *gorm.DB
}
