package service

import (
	"context"

	gomail "gopkg.in/mail.v2"

	"github.com/lakhansamani/cloud-container/graph/model"
	"github.com/lakhansamani/cloud-container/internal/db"
	mp "github.com/lakhansamani/cloud-container/internal/memorystore/providers"
)

// Dependencies is a struct that contains all the dependencies of the service
// It is injected into Resolvers so that we have all the service dependencies available in the resolvers
type Dependencies struct {
	// DatabaseClient is the database client
	DatabaseClient *db.DBProvider
	// Mailer is the mailing client
	Mailer *gomail.Dialer
	// MemoryStoreProvider is the memory store provider
	MemoryStoreProvider mp.MemoryStoreProvider
}

// Service is the interface that all services must implement
type Service interface {
	// IsReady is called during the health check.
	// Only return true when everything is really ready.
	IsReady() bool
	// Signup is the service for the signup mutation
	// permission required: none
	Signup(ctx context.Context, params model.SignUpRequest) (*model.Response, error)
	// Login is the service for the login mutation
	// permission required: none
	Login(ctx context.Context, params model.LoginRequest) (*model.Response, error)
	// VerifyOTP is the service for the verify_otp mutation
	// permission required: none
	VerifyOTP(ctx context.Context, params model.VerifyOtpRequest) (*model.AuthResponse, error)
	// Session is the service for the session query
	// permission required: none
	Session(ctx context.Context) (*model.AuthResponse, error)
	// Logout is the service for the logout mutation
	// permission required: authenticated user
	Logout(ctx context.Context) (*model.Response, error)
	// Deployments is the service for the deployments query
	// permission required: authenticated user
	CreateDeployment(ctx context.Context, params *model.CreateDeploymentRequest) (*model.Deployment, error)
	// Deployment is the service for the deployment query
	// permission required: authenticated user
	DeleteDeployment(ctx context.Context, params *model.DeleteDeploymentRequest) (*model.Response, error)
	// Deployments is the service for the deployments query
	// permission required: authenticated user
	Deployments(ctx context.Context, params *model.ListDeploymentsRequest) ([]*model.Deployment, error)
	// Deployment is the service for the deployment query
	// permission required: authenticated user
	Deployment(ctx context.Context, params *model.GetDeploymentRequest) (*model.Deployment, error)
}

type service struct {
	// Dependencies is the dependencies of the service
	Dependencies
}

// NewService returns a new Service
func NewService(svcDeps *Dependencies) Service {
	return &service{
		Dependencies: *svcDeps,
	}
}

// IsReady is called during the health check.
// Only return true when everything is really ready.
func (s *service) IsReady() bool {
	return true
}
