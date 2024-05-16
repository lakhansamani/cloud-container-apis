package cmd

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/lakhansamani/cloud-container/internal/global"
	"github.com/lakhansamani/cloud-container/internal/server"
)

const (
	// envPath is the path to the .env file
	envPath = ".env"
	// portEnvKey is the key for the port env variable
	portEnvKey = "PORT"
	// databaseURLKey is the key for the database URL env variable
	databaseURLKey = "DATABASE_URL"
	// redisURL is the key for the redis URL env variable
	redisURL = "REDIS_URL"
	// smtpHostKey is the key for the smtp host env variable
	smtpHostKey = "SMTP_HOST"
	// smtpPortKey is the key for the smtp port env variable
	smtpPortKey = "SMTP_PORT"
	// smtpUsernameKey is the key for the smtp username env variable
	smtpUsernameKey = "SMTP_USERNAME"
	// smtpPasswordKey is the key for the smtp password env variable
	smtpPasswordKey = "SMTP_PASSWORD"
	// smtpSenderEmailKey is the key for the smtp sender email env variable
	smtpSenderEmailKey = "SMTP_SENDER_EMAIL"
	// smtpSenderNameKey is the key for the smtp sender name env variable
	smtpSenderNameKey = "SMTP_SENDER_NAME"
)

var (
	// RootCmd is the root (and only) command of this service
	// TODO change this to your docker image name
	RootCmd = &cobra.Command{
		Use:   "api",
		Short: "The api Service",
		Run:   runRootCmd,
	}

	rootArgs struct {
		// Version of the service
		version string
		// Log level
		logLevel string
		// Server configuration
		server server.Config
	}
)

// SetVersion stores the given version
func SetVersion(version, build string) {
	rootArgs.version = version
}

// SetEnvVars sets the environment variables
func SetEnvVars() {
	// Set global variables and required env variables
	global.Port = os.Getenv(portEnvKey)
	if global.Port == "" {
		global.Port = "3000"
	}
	global.DatabaseURL = os.Getenv(databaseURLKey)
	if global.DatabaseURL == "" {
		log.Fatal().Msg("DATABASE_URL not set")
	}
	global.ContainerOrchestratorServiceURL = os.Getenv("CONTAINER_ORCHESTRATOR_SERVICE_URL")
	if global.ContainerOrchestratorServiceURL == "" {
		log.Fatal().Msg("CONTAINER_ORCHESTRATOR_SERVICE_URL not set")
	}
	global.RedisURL = os.Getenv(redisURL)
	if global.RedisURL == "" {
		log.Fatal().Msg("REDIS_URL not set")
	}
	global.SMTPHost = os.Getenv(smtpHostKey)
	if global.SMTPHost == "" {
		log.Fatal().Msg("SMTP_HOST not set")
	}
	global.SMTPPort = os.Getenv(smtpPortKey)
	if global.SMTPPort == "" {
		log.Fatal().Msg("SMTP_PORT not set")
	}
	global.SMTPUsername = os.Getenv(smtpUsernameKey)
	if global.SMTPUsername == "" {
		log.Fatal().Msg("SMTP_USERNAME not set")
	}
	global.SMTPPassword = os.Getenv(smtpPasswordKey)
	if global.SMTPPassword == "" {
		log.Fatal().Msg("SMTP_PASSWORD not set")
	}
	global.SMTPSenderEmail = os.Getenv(smtpSenderEmailKey)
	if global.SMTPSenderEmail == "" {
		log.Fatal().Msg("SMTP_SENDER_EMAIL not set")
	}
	global.SMTPSenderName = os.Getenv(smtpSenderNameKey)
	if global.SMTPSenderName == "" {
		log.Fatal().Msg("SMTP_SENDER_NAME not set")
	}

}

func init() {
	// Load env variables from .env if present
	godotenv.Load(envPath)
	// Set env variables
	SetEnvVars()
	// Setup flags
	f := RootCmd.Flags()
	portInt, _ := strconv.ParseInt(global.Port, 0, 64)
	// Logging flags
	f.StringVar(&rootArgs.logLevel, "log-level", "debug", "Minimum log level")
	// Server flags
	f.IntVar(&rootArgs.server.Port, "http-port", int(portInt), "Port to listen on for HTTP requests")
}

func runRootCmd(cmd *cobra.Command, args []string) {
	// Setup logging
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFieldName = "time"
	log := zerolog.New(cmd.OutOrStderr()).With().Timestamp().Logger()
	// Set log level
	logLevel, err := zerolog.ParseLevel(rootArgs.logLevel)
	if err != nil {
		// Default to debug if the log level is invalid
		logLevel = zerolog.DebugLevel
	}
	log.Level(logLevel)
	// Setup server
	srv, err := server.New(log, rootArgs.server)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create server")
	}
	// Run server
	ctx := cmd.Context()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return srv.Run(ctx) })
	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}
