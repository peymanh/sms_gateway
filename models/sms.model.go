package models

import (
	"time"

	"github.com/google/uuid"
)

type SMSResultStatus int

const (
	SMSResultPending SMSResultStatus = 0
	SMSResultSuccess SMSResultStatus = 1
	SMSResultFailed  SMSResultStatus = 2
	SMSResultError   SMSResultStatus = 3
)

type SMSLog struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Body         string          `gorm:"not null" json:"body,omitempty"`
	Language     string          `gorm:"not null" json:"language,omitempty"`
	Receiver     string          `gorm:"not null" json:"receiver,omitempty"`
	Cost         int             `gorm:"not null" json:"cost,omitempty"`
	UserID       uuid.UUID       `gorm:"not null" json:"user_id,omitempty"`
	User         User            `gorm:"constraint:OnUpdate:CASCADE;"`
	Status       SMSResultStatus `gorm:"not null;default:0" json:"status,omitempty"`
	ErrorMessage string          `gorm:"not null" json:"error_message,omitempty"`
	CreatedAt    time.Time       `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time       `gorm:"not null" json:"updated_at,omitempty"`
}

type SendInput struct {
	Body     string `json:"body" binding:"required"`
	Language string `json:"language" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
}
