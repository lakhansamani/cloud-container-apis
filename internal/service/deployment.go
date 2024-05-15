package service

import (
	"context"
	"errors"

	"github.com/lakhansamani/cloud-container/graph/model"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
	"github.com/rs/zerolog/log"
)

// Deployment is the service for the deployment query
// permission required: authenticated user
func (s *service) Deployment(ctx context.Context, params *model.GetDeploymentRequest) (*model.Deployment, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	userID := gc.Value("user_id").(string)
	depl, err := s.DatabaseClient.GetDeploymentByID(params.ID, userID)
	if err != nil {
		log.Debug().Err(err).Msg("error getting deployment")
		return nil, errors.New(messages.InternalServerError)
	}
	return depl.ToAPI(), nil
}
