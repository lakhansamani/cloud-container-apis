package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lakhansamani/cloud-container/graph/model"
	"github.com/lakhansamani/cloud-container/internal/db/models"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
	"github.com/rs/zerolog/log"
)

// Deployments is the service for the deployments query
// permission required: authenticated user
func (s *service) CreateDeployment(ctx context.Context, params *model.CreateDeploymentRequest) (*model.Deployment, error) {
	if params.Name == "" {
		return nil, errors.New(messages.InvalidDeploymentRequestError)
	}
	if params.Image == "" {
		return nil, errors.New(messages.InvalidDeploymentRequestError)
	}
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	userID := gc.Value("user_id").(string)
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing user id")
		return nil, errors.New(messages.InternalServerError)
	}
	depl := &models.Deployment{
		Name:      params.Name,
		Image:     params.Image,
		Status:    "pending",
		CreatedBy: userIDUUID,
		EnvVars:   params.EnvVars,
	}
	depl, err = s.DatabaseClient.CreateDeployment(depl)
	if err != nil {
		log.Debug().Err(err).Msg("error creating deployment")
		return nil, errors.New(messages.InternalServerError)
	}
	return depl.ToAPI(), nil
}
