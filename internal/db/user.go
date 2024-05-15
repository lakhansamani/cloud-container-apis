package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lakhansamani/cloud-container/internal/db/models"
)

// GetUserByID returns a user by ID
func (p *DBProvider) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := p.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail returns a user by email
func (p *DBProvider) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := p.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (p *DBProvider) CreateUser(user *models.User) (*models.User, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := p.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user
func (p *DBProvider) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	return p.DB.Save(&user).Error
}
