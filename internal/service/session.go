package service

import (
	"context"

	"github.com/lakhansamani/cloud-container/graph/model"
)

// Session is the service for the session query.
// permission required: none
func (s *service) Session(ctx context.Context, params *model.SessionQueryInput) (*model.AuthResponse, error) {
	return nil, nil
}
