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
func (s *service) DeleteDeployment(ctx context.Context, params *model.DeleteDeploymentRequest) (*model.Response, error) {
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
	if depl == nil {
		return nil, errors.New(messages.DeploymentNotFoundError)
	}
	err = s.DatabaseClient.DeleteDeployment(depl)
	if err != nil {
		log.Debug().Err(err).Msg("error deleting deployment")
		return nil, errors.New(messages.InternalServerError)
	}
	return &model.Response{
		Message: messages.DeploymentDeletedMessage,
	}, nil
}
