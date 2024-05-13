package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/lakhansamani/cloud-container/graph/model"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
)

// Logout is the service for the logout mutation
// permission required: authenticated user
func (s *service) Logout(ctx context.Context) (*model.Response, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, err
	}
	// Get session from cookie
	cookie, err := gc.Request.Cookie(constants.SessionCookieName)
	if err != nil {
		log.Debug().Err(err).Msg("error getting session cookie")
		return nil, err
	}
	// Decrypt session token
	userID, nonce, err := utils.DecryptSession(cookie.Value)
	if err != nil {
		log.Debug().Err(err).Msg("error decrypting session token")
		return nil, err
	}
	// Get session from memory store
	_, err = s.MemoryStoreProvider.GetMfaSession(userID, nonce)
	if err != nil {
		log.Debug().Err(err).Msg("error getting session from memory store")
		return nil, err
	}
	// Delete session from memory store
	if err := s.MemoryStoreProvider.DeleteUserSession(userID, nonce); err != nil {
		log.Debug().Err(err).Msg("error deleting session from memory store")
		return nil, err
	}
	return &model.Response{
		Message: messages.LogoutMessage,
	}, nil
}
