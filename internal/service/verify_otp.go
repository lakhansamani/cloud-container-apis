package service

import (
	"context"
	"errors"
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
	mfaSession, err := gc.Request.Cookie(constants.MfaSessionCookieName)
	if err != nil {
		log.Debug().Err(err).Msg("error getting mfa session cookie")
		return nil, errors.New(messages.InvalidMfaSessionError)
	}
	// Split session
	splitSession := strings.Split(mfaSession.Value, ":")
	if len(splitSession) != 2 {
		log.Debug().Msg(messages.InvalidMfaSessionError)
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
	// Get user from database
	user, err := s.DatabaseClient.GetUserByID(userID)
	if err != nil {
		log.Debug().Err(err).Msg("error getting user from database")
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
	session, err := utils.GenerateSession(user, nonce)
	if err != nil {
		log.Debug().Err(err).Msg("error generating session")
		return nil, errors.New(messages.ErrorGeneratingSession)
	}
	s.MemoryStoreProvider.SetUserSession(user.ID, nonce, session)
	gc.SetCookie(constants.SessionCookieName, session, -1, "/", "", true, true)
	return &model.AuthResponse{
		User:    user.ToAPI(),
		Message: messages.OTPSentMessage,
	}, nil
}
