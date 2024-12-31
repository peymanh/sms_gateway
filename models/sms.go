package models

import "gorm.io/gorm"

type SMS struct {
	gorm.Model
	From    string `gorm:"not null"`
	To      string `gorm:"not null"`
	Message string `gorm:"not null"`
	UserID  int    `gorm:"not null"`
}
