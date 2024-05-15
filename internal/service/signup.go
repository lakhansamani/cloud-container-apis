package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/lakhansamani/cloud-container/graph/model"
	constants "github.com/lakhansamani/cloud-container/internal/contants"
	"github.com/lakhansamani/cloud-container/internal/db/models"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
)

// Signup is the service for the signup mutation
// permission required: none
func (s *service) Signup(ctx context.Context, params model.SignUpRequest) (*model.Response, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	email := strings.ToLower(strings.TrimSpace(params.Email))
	// Validate email
	if email == "" {
		return nil, errors.New(messages.InvalidEmailAddressError)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New(messages.InvalidEmailAddressError)
	}
	// Check if user exists
	if u, err := s.DatabaseClient.GetUserByEmail(email); err == nil && u != nil {
		return nil, errors.New(messages.UserAlreadyExistsError)
	}
	// Create user
	user := &models.User{
		Email:     email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}
	user, err = s.DatabaseClient.CreateUser(user)
	if err != nil {
		return nil, err
	}
	// Set MFA session in memory store
	mfaSession := uuid.NewString()
	otp := fmt.Sprintf("%d", rand.Intn(999999))
	mfaSessionExpiresIn := time.Now().Add(time.Minute * 2).Unix()
	s.MemoryStoreProvider.SetMfaSession(user.ID.String(), mfaSession, otp, mfaSessionExpiresIn)
	host := utils.GetHost(gc)
	hostname, _ := utils.GetHostParts(host)
	gc.SetCookie(constants.MfaSessionCookieName, fmt.Sprintf("%s:%s", user.ID, mfaSession), 60*2, "/", hostname, true, true)
	// TODO send OTP to user
	log.Debug().Msgf("OTP: %s", otp)
	return &model.Response{
		Message: messages.OTPSentMessage,
	}, nil
}
