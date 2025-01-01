package services

import (
	"context"
	"errors"
	"time"

	"github.com/peymanh/sms_gateway/models"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

var (
	normalChan  = make(chan models.SMSResultStatus, 200)
	premiumChan = make(chan models.SMSResultStatus, 500)
)

type SMSService struct {
	DB *gorm.DB
}

func NewSMSService(DB *gorm.DB) *SMSService {
	return &SMSService{DB}
}

func (s *SMSService) SendSMS(ctx context.Context, user *models.User, log *models.SMSLog, receiver string, body string, language string) error {
	// Create a channel to receive the status
	var statusChan chan models.SMSResultStatus
	switch user.Class {
	case models.UserTypeNormal:
		statusChan = normalChan // Buffer of size 200 for normal users
	case models.UserTypePremium:
		statusChan = premiumChan // Buffer of size 500 for premium users
	default:
		return errors.New("invalid user type")
	}

	// Launch a goroutine for asynchronous SMS sending
	go func(receiver string, body string, language string, statusChan chan<- models.SMSResultStatus) {
		defer close(statusChan) // Close the channel when done

		// Simulate SMS sending with a delay
		time.Sleep(time.Second * 2)

		// Decide on mock SMS status (success/failure)
		randomFloat := rand.Float64()
		var state models.SMSResultStatus
		if randomFloat < 0.9 {
			state = models.SMSResultSuccess
		} else {
			state = models.SMSResultFailed
		}

		statusChan <- state
	}(receiver, body, language, statusChan)

	// Handle potential context cancellation
	select {
	case <-ctx.Done():
		log.Status = models.SMSResultError
		log.ErrorMessage = ctx.Err().Error()
		s.DB.Save(log)
		return ctx.Err()
	case state := <-statusChan:
		log.Status = state
		s.DB.Save(log)
		return nil
	}
}
