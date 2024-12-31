package services

import (
	"fmt"

	"github.com/peymanh/sms_gateway/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(username string) (*models.User, error) {
	user := &models.User{Username: username}
	if err := s.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if err := s.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *UserService) TopUpBalance(user *models.User, amount float64) error {
	user.Balance += amount
	if err := s.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to top up balance: %w", err)
	}
	return nil
}
