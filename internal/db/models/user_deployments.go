package models

import (
	"time"

	"github.com/google/uuid"
)

// UserDeployments is the struct for the user_deployments table
type UserDeployments struct {
	ID uuid.UUID `gorm:"type:uuid;"`
	// CreatedAt is the time the deployment was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the deployment was updated
	UpdatedAt time.Time `json:"updated_at"`
	// UserID is the user id
	UserID uuid.UUID `json:"user_id"`
	// User is the user
	User User `gorm:"column:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// DeploymentID is the deployment id
	DeploymentID uuid.UUID `json:"deployment_id"`
	// Deployment is the deployment
	Deployment Deployment `gorm:"column:deployment_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// UserRole is the role of the user
	UserRole string `json:"user_role"`
}
