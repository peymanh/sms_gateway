package services

import (
	"errors"
	"fmt"

	"github.com/peymanh/sms_gateway/models"
	"gorm.io/gorm"
)

type SMSService struct {
	DB *gorm.DB
}

func NewSMSService(db *gorm.DB) *SMSService {
	return &SMSService{DB: db}
}

func (s *SMSService) SendSMS(from, to, message string, userID int) error {
	// Fetch user from database
	user := &models.User{}
	if err := s.DB.First(user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check balance
	if user.Balance < 0.1 { // Assuming 0.1 is the SMS fee
		return errors.New("insufficient balance")
	}

	// Create SMS record
	sms := &models.SMS{
		From:    from,
		To:      to,
		Message: message,
		UserID:  userID,
	}

	// Save SMS to database
	if err := s.DB.Create(sms).Error; err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Update user balance
	user.Balance -= 0.1
	if err := s.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Simulate external SMS sending (replace with actual API call)
	fmt.Println("SMS sent successfully:", sms)

	return nil
}

func (s *SMSService) GetSMSLog(userID uint) ([]models.SMS, error) {
	var smsLog []models.SMS
	if err := s.DB.Where("user_id = ?", userID).Find(&smsLog).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch SMS log: %w", err)
	}
	return smsLog, nil
}
