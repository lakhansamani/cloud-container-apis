package db

import (
	"github.com/lakhansamani/cloud-container/internal/db/models"
	"github.com/lakhansamani/cloud-container/internal/global"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBProvider struct {
	// DB is the database connection
	DB *gorm.DB
}

func NewDBProvider() (*DBProvider, error) {
	db, err := gorm.Open(postgres.Open(global.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
		return nil, err
	}
	err = db.AutoMigrate(&models.User{}, &models.Deployment{})
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating models")
		return nil, err
	}
	return &DBProvider{
		DB: db,
	}, nil
}
