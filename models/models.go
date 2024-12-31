package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
