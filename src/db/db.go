package db

import (
	"github.com/joaooliveira247/go_olist_challenge/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		return nil, err
	}
	return db, nil
}
