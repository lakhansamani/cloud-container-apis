package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/lakhansamani/cloud-container/graph/model"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
)

// VerifyOTP is the service for the verify_otp mutation
// permission required: none
func (s *service) VerifyOTP(ctx context.Context, params model.VerifyOtpRequest) (*model.AuthResponse, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	host := utils.GetHost(gc)
	hostname, _ := utils.GetHostParts(host)
	mfaSession, err := gc.Request.Cookie(constants.MfaSessionCookieName)
	if err != nil {
		log.Debug().Err(err).Msg("error getting mfa session cookie")
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	// Split session
	mfaSessionValue, err := url.PathUnescape(mfaSession.Value)
	if err != nil {
		log.Debug().Err(err).Msg("error unescaping mfa session value")
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	splitSession := strings.Split(mfaSessionValue, ":")
	if len(splitSession) != 2 {
		log.Debug().Interface("mfa_session", mfaSession).Msg(messages.InvalidMfaSessionError)
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	// Get user id from session
	userID := splitSession[0]
	mfaSessionToken := splitSession[1]
	// Get otp from memory store
	otp, err := s.MemoryStoreProvider.GetMfaSession(userID, mfaSessionToken)
	if err != nil {
		log.Debug().Err(err).Msg("error getting mfa session from memory store")
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	if otp != params.Otp {
		log.Debug().Msg(messages.InvalidOtpError)
		return nil, errors.New(messages.InvalidOtpError)
	}
	// Delete mfa session from memory store
	if err := s.MemoryStoreProvider.DeleteMfaSession(userID, mfaSessionToken); err != nil {
		log.Debug().Err(err).Msg("error deleting mfa session from memory store")
		return nil, errors.New(messages.InternalServerError)
	}
	gc.SetCookie(constants.MfaSessionCookieName, "", -1, "/", hostname, true, true)
	// Get user from database
	user, err := s.DatabaseClient.GetUserByID(userID)
	if err != nil {
		log.Debug().Err(err).Str("user_id", user.ID).Msg("error getting user from database")
		return nil, errors.New(messages.UserNotFoundError)
	}
	if user.VerifiedAt == nil {
		// Update user verified_at
		t := time.Now()
		user.VerifiedAt = &t
		if err := s.DatabaseClient.UpdateUser(user); err != nil {
			log.Debug().Err(err).Msg("error updating user")
			// continue
		}
	}
	// Set session in memory store
	nonce := uuid.NewString()
	session, err := utils.GenerateSession(user.ID, nonce)
	if err != nil {
		log.Debug().Err(err).Msg("error generating session")
		return nil, errors.New(messages.ErrorGeneratingSession)
	}
	s.MemoryStoreProvider.SetUserSession(user.ID, nonce, session)
	gc.SetCookie(constants.SessionCookieName, session, 60*60*24*120, "/", hostname, true, true)
	return &model.AuthResponse{
		User:    user.ToAPI(),
		Message: messages.LoginSuccessMessage,
	}, nil
}
