package graph

import (
	"strconv"

	container "github.com/lakhansamani/container-orchestrator-apis/container"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gomail "gopkg.in/mail.v2"

	"github.com/lakhansamani/cloud-container/internal/db"
	"github.com/lakhansamani/cloud-container/internal/global"
	"github.com/lakhansamani/cloud-container/internal/memorystore"
	"github.com/lakhansamani/cloud-container/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service service.Service
}

func NewResolver() *Resolver {
	// Initialize database
	db, err := db.NewDBProvider()
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing database")
	}
	// Initialize mailer
	smtPort, _ := strconv.Atoi(global.SMTPPort)
	mailer := gomail.NewDialer(global.SMTPHost, smtPort, global.SMTPUsername, global.SMTPPassword)
	// Initialize memory store
	memoryStore, err := memorystore.NewMemoryStore()
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing memory store")
	}
	conn, err := grpc.NewClient(global.ContainerOrchestratorServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing container orchestrator service")
	}
	containerClient := container.NewContainerServiceClient(conn)
	svc := service.NewService(&service.Dependencies{
		DatabaseClient:         db,
		Mailer:                 mailer,
		MemoryStoreProvider:    memoryStore,
		ContainerServiceClient: containerClient,
	})
	return &Resolver{
		Service: svc,
	}
}
