package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/lakhansamani/cloud-container/graph/model"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
)

// Login is the service for the login mutation.
// permission required: none
func (s *service) Login(ctx context.Context, params model.LoginRequest) (*model.Response, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	// Check if user exists
	user, err := s.DatabaseClient.GetUserByEmail(params.Email)
	if err != nil {
		return nil, errors.New(messages.UserNotFoundError)
	}
	// Set MFA session in memory store
	mfaSession := uuid.NewString()
	otp := fmt.Sprintf("%d", rand.Intn(999999))
	mfaSessionExpiresIn := time.Now().Add(time.Minute * 2).Unix()
	s.MemoryStoreProvider.SetMfaSession(user.ID, mfaSession, otp, mfaSessionExpiresIn)
	host := utils.GetHost(gc)
	hostname, _ := utils.GetHostParts(host)
	gc.SetCookie(constants.MfaSessionCookieName, fmt.Sprintf("%s:%s", user.ID, mfaSession), 60*2, "/", hostname, true, true)
	// TODO send OTP to user
	log.Debug().Msgf("OTP: %s", otp)
	return &model.Response{
		Message: messages.OTPSentMessage,
	}, nil
}
