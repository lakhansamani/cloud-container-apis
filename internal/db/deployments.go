package db

import (
	"github.com/lakhansamani/cloud-container/internal/db/models"
)

// CreateDeployment creates a new deployment
func (p *DBProvider) CreateDeployment(depl *models.Deployment) error {
	return p.DB.Create(&depl).Error
}

// DeleteDeployment deletes a deployment
func (p *DBProvider) DeleteDeployment(depl *models.Deployment) error {
	// set deleted_at to current time
	return p.DB.Delete(&depl).Error
}
