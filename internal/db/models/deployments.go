package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB Interface for JSONB Field
type JSONB []map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

// Deployment is the struct for the deployment table
type Deployment struct {
	ID uuid.UUID `gorm:"type:uuid;"`
	// CreatedAt is the time the deployment was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the deployment was updated
	UpdatedAt time.Time `json:"updated_at"`
	// Name is the name of the deployment
	Name string `json:"name"`
	// Image is the image of the deployment
	Image string `json:"image"`
	// Status is the status of the deployment
	Status string `json:"status"`
	// ContainerID is the container id of the deployment
	ContainerID string `json:"container_id"`
	// EnvVars are the environment variables of the deployment
	EnvVars JSONB `json:"env_vars"`
	// CreatedBy is the user who created the deployment
	CreatedBy uuid.UUID `json:"created_by"`
	// CreatedBy is the user who created the deployment
	CreatedByUser User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// DeletedAt is the time the deployment was deleted
	DeletedAt *time.Time `json:"deleted_at"`
}
