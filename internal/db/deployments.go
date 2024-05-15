package db

import (
	"time"

	"github.com/google/uuid"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/db/models"
)

// CreateDeployment creates a new deployment
func (p *DBProvider) CreateDeployment(depl *models.Deployment) (*models.Deployment, error) {
	depl.CreatedAt = time.Now()
	depl.UpdatedAt = time.Now()
	depl.ID = uuid.New()
	err := p.DB.Create(&depl).Error
	if err != nil {
		return nil, err
	}
	// Create user deployment entry for the owner
	userDepl := &models.UserDeployments{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserID:       depl.CreatedBy,
		DeploymentID: depl.ID,
		UserRole:     constants.RoleTypeDeploymentOwner,
	}
	err = p.DB.Create(&userDepl).Error
	if err != nil {
		return nil, err
	}
	return depl, nil
}

// ListDeployments lists all deployments for a user based on pagination
func (p *DBProvider) ListDeployments(user_id string, limit int, offset int) ([]*models.Deployment, error) {
	var depls []*models.Deployment
	err := p.DB.Where("created_by = ? AND deleted_at IS NULL", user_id).Limit(limit).Offset(offset).Find(&depls).Error
	return depls, err
}

// GetDeploymentByID gets a deployment by id
func (p *DBProvider) GetDeploymentByID(id string, user_id string) (*models.Deployment, error) {
	var depl models.Deployment
	err := p.DB.Where("id = ? AND created_by = ?", id, user_id).First(&depl).Error
	return &depl, err
}

// UpdateDeployment updates a deployment
func (p *DBProvider) UpdateDeployment(depl *models.Deployment) error {
	depl.UpdatedAt = time.Now()
	return p.DB.Save(&depl).Error
}

// DeleteDeployment deletes a deployment
func (p *DBProvider) DeleteDeployment(depl *models.Deployment) error {
	// set deleted_at to current time
	t := time.Now()
	depl.DeletedAt = &t
	depl.UpdatedAt = t
	depl.Status = "removed"
	return p.DB.Save(&depl).Error
}
