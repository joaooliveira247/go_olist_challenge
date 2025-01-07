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

func DeleteAllTables(db *gorm.DB) error {
	rawQuery := `
do $$ declare
    r record;
begin
    for r in (select tablename from pg_tables where schemaname = 'public') loop
        execute 'drop table if exists ' || quote_ident(r.tablename) || ' cascade';
    end loop;
end $$;
		`
	if err := db.Exec(rawQuery).Error; err != nil {
		return err
	}
	return nil
}
