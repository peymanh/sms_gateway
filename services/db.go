package services

import (
	"log"

	"github.com/peymanh/sms_gateway/config"
	"github.com/peymanh/sms_gateway/models"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	db, err := models.ConnectDB(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
