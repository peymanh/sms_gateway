package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"unique;not null"`
	Balance  float64 `gorm:"default:0.0"`
}
