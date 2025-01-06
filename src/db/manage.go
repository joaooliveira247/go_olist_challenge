package db

import (
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"gorm.io/gorm"
)

func CreateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Author{}, &models.BookAuthor{}, &models.Book{}); err != nil {
		return err
	}
	return nil
}
