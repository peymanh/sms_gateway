package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/peymanh/sms_gateway/models"
	"gorm.io/gorm"
)

var (
	normalChan  = make(chan bool, 200)
	premiumChan = make(chan bool, 500)
)

type SMSService struct {
	DB *gorm.DB
}

func NewSMSService(DB *gorm.DB) *SMSService {
	return &SMSService{DB}
}

func (s *SMSService) SendSMS(ctx context.Context, user *models.User, log *models.SMSLog, receiver string, body string, language string) error {
	// Generate a unique ID for this SMS
	smsID := uuid.New().String()

	// Create a channel to receive the status
	var statusChan chan bool
	switch user.Class {
	case models.UserTypeNormal:
		statusChan = normalChan // Buffer of size 10 for normal users
	case models.UserTypePremium:
		statusChan = premiumChan // Buffer of size 50 for premium users
	default:
		return errors.New("invalid user type")
	}

	// Launch a goroutine for asynchronous SMS sending
	go func(smsID string, receiver string, body string, language string, statusChan chan<- bool) {
		defer close(statusChan) // Close the channel when done

		// Simulate SMS sending with a delay
		time.Sleep(time.Second * 2)

		// Decide on mock SMS status (success/failure)
		success := true // Change this for testing failure

		if success {
			statusChan <- true
		} else {
			statusChan <- false
		}
	}(smsID, receiver, body, language, statusChan)

	// Handle potential context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case state := <-statusChan:
		log.Status = state
		s.DB.Save(log)
		return nil
	}
}
