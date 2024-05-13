package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/google/uuid"
	"github.com/lakhansamani/cloud-container/graph/model"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
	"github.com/rs/zerolog/log"
)

// Session is the service for the session query.
// permission required: none
func (s *service) Session(ctx context.Context) (*model.AuthResponse, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	host := utils.GetHost(gc)
	hostname, _ := utils.GetHostParts(host)
	// Get session from cookie
	cookie, err := gc.Request.Cookie(constants.SessionCookieName)
	if err != nil {
		log.Debug().Err(err).Msg("error getting session cookie")
		return nil, errors.New(messages.InvalidSessionError)
	}
	sessionValue, err := url.PathUnescape(cookie.Value)
	if err != nil {
		log.Debug().Err(err).Msg("error unescaping mfa session value")
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	// Decrypt session token
	userID, nonce, err := utils.DecryptSession(sessionValue)
	if err != nil {
		log.Debug().Err(err).Msg("error decrypting session token")
		return nil, errors.New(messages.InvalidSessionError)
	}
	// Get session from memory store
	_, err = s.MemoryStoreProvider.GetUserSession(userID, nonce)
	if err != nil {
		log.Debug().Err(err).Msg("error getting session from memory store")
		return nil, errors.New(messages.InvalidSessionError)
	}
	// Delete session from memory store
	if err := s.MemoryStoreProvider.DeleteUserSession(userID, nonce); err != nil {
		log.Debug().Err(err).Msg("error deleting session from memory store")
		// continue
	}
	// Set session in memory store
	nonce = uuid.NewString()
	session, err := utils.GenerateSession(userID, nonce)
	if err != nil {
		log.Debug().Err(err).Msg("error generating session")
		return nil, errors.New(messages.ErrorGeneratingSession)
	}
	s.MemoryStoreProvider.SetUserSession(userID, nonce, session)
	gc.SetCookie(constants.SessionCookieName, session, 60*60*24*120, "/", hostname, true, true)
	user, err := s.DatabaseClient.GetUserByID(userID)
	if err != nil {
		log.Debug().Err(err).Msg("error getting user from database")
		return nil, errors.New(messages.InternalServerError)
	}
	return &model.AuthResponse{
		Message: messages.LoginSuccessMessage,
		User:    user.ToAPI(),
	}, nil
}
