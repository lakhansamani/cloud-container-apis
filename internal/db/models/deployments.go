package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
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
	ID string `json:"id" gorm:"primary_key"`
	// CreatedAt is the time the deployment was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the deployment was updated
	UpdatedAt time.Time `json:"updated_at"`
	// Name is the name of the deployment
	Name string `json:"name"`
	// Image is the image of the deployment
	Image string `json:"image"`
	// Port is the port of the deployment
	Port string `json:"port"`
	// EnvironmentVariables are the environment variables of the deployment
	EnvironmentVariables JSONB `json:"environment"`
}
