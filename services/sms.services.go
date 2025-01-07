package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/peymanh/sms_gateway/models"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

var (
	normalChan  = make(chan MockedResponse, 200)
	premiumChan = make(chan MockedResponse, 500)
)

type SMSService struct {
	DB *gorm.DB
}

type MockedResponse struct {
	SMS    *models.SMSLog
	Status models.SMSResultStatus
}

func NewSMSService(DB *gorm.DB) *SMSService {
	return &SMSService{DB}
}

func (s *SMSService) SendSMS(ctx context.Context, user *models.User, log *models.SMSLog, receiver string, body string, language string) error {
	// Create a channel to receive the status
	var statusChan chan MockedResponse
	switch user.Class {
	case models.UserTypeNormal:
		statusChan = normalChan // Buffer of size 200 for normal users
	case models.UserTypePremium:
		statusChan = premiumChan // Buffer of size 500 for premium users
	default:
		return errors.New("invalid user type")
	}

	// Launch a goroutine for asynchronous SMS sending
	go func(log *models.SMSLog, receiver string, body string, language string, statusChan chan<- MockedResponse) {
		// defer close(statusChan) // Close the channel when done

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
		m := MockedResponse{
			SMS:    log,
			Status: state,
		}

		statusChan <- m
	}(log, receiver, body, language, statusChan)

	// Handle potential context cancellation
	select {
	case <-ctx.Done():
		log.Status = models.SMSResultError
		log.ErrorMessage = ctx.Err().Error()
		s.DB.Save(log)
		return ctx.Err()
	}
	return nil
}

func (s *SMSService) Listen() {
	for {
		if response, ok := <-normalChan; ok {
			sms := response.SMS
			log.Printf("normal chan received SMS:  %v", sms)
			sms.Status = response.Status
			s.DB.Save(response.SMS)
			continue
		}

		// Check if there's a message on premiumChan
		if response, ok := <-premiumChan; ok {
			sms := response.SMS
			log.Printf("normal chan received SMS:  %v", sms)
			sms.Status = response.Status
			s.DB.Save(response.SMS)
			continue
		}
	}
}
