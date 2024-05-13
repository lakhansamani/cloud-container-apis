package models

import (
	"time"

	"github.com/lakhansamani/cloud-container/graph/model"
)

// User is the struct for the user table
type User struct {
	ID string `json:"id" gorm:"primary_key"`
	// CreatedAt is the time the deployment was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the deployment was updated
	UpdatedAt time.Time `json:"updated_at"`
	// Email is the email of the user
	Email string `json:"email" gorm:"unique"`
	// FirstName is the first name of the user
	FirstName string `json:"first_name"`
	// LastName is the last name of the user
	LastName string `json:"last_name"`
	// VerifiedAt is the time the user was verified
	VerifiedAt *time.Time `json:"verified_at"`
}

func (u *User) ToAPI() *model.User {
	isVerified := false
	if u.VerifiedAt != nil {
		isVerified = true
	}
	return &model.User{
		ID:         u.ID,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		IsVerified: isVerified,
	}
}
