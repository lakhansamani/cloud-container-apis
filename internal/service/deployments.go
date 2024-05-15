package service

import (
	"context"
	"errors"

	"github.com/lakhansamani/cloud-container/graph/model"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
	"github.com/rs/zerolog/log"
)

// Deployments is the service for the deployments query
// permission required: authenticated user
func (s *service) Deployments(ctx context.Context, params *model.ListDeploymentsRequest) ([]*model.Deployment, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	userID := gc.Value("user_id").(string)
	limit := 100
	offset := 0
	if params != nil && params.Limit != nil {
		limit = *params.Limit
	}
	if params != nil && params.Offset != nil {
		offset = *params.Offset
	}
	depls, err := s.DatabaseClient.ListDeployments(userID, limit, offset)
	if err != nil {
		log.Debug().Err(err).Msg("error getting deployments")
		return nil, errors.New(messages.InternalServerError)
	}
	res := []*model.Deployment{}
	for _, depl := range depls {
		res = append(res, depl.ToAPI())
	}
	return res, nil
}
