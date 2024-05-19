package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lakhansamani/container-orchestrator-apis/container"
	"github.com/rs/zerolog/log"

	"github.com/lakhansamani/cloud-container/graph/model"
	"github.com/lakhansamani/cloud-container/internal/db/models"
	"github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/utils"
)

// Deployments is the service for the deployments query
// permission required: authenticated user
func (s *service) CreateDeployment(ctx context.Context, params *model.CreateDeploymentRequest) (*model.Deployment, error) {
	if params.Name == "" {
		return nil, errors.New(messages.InvalidDeploymentRequestError)
	}
	if params.Image == "" {
		return nil, errors.New(messages.InvalidDeploymentRequestError)
	}
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg(messages.GinContextError)
		return nil, errors.New(messages.GinContextError)
	}
	userID := gc.Value("user_id").(string)
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing user id")
		return nil, errors.New(messages.InternalServerError)
	}
	// Create deployment in database
	depl := &models.Deployment{
		Name:      params.Name,
		Image:     params.Image,
		Status:    "pending",
		CreatedBy: userIDUUID,
		EnvVars:   params.EnvVars,
	}
	depl, err = s.DatabaseClient.CreateDeployment(depl)
	if err != nil {
		log.Debug().Err(err).Msg("error creating deployment")
		return nil, errors.New(messages.InternalServerError)
	}
	// Create container
	containerEnvVars := []*container.EnvVar{}
	for k, v := range depl.EnvVars {
		containerEnvVars = append(containerEnvVars, &container.EnvVar{
			Key:   k,
			Value: v.(string),
		})
	}
	newContainer, err := s.ContainerServiceClient.CreateContainer(ctx, &container.CreateContainerRequest{
		Image:   depl.Image,
		Name:    depl.ID.String(),
		EnvVars: containerEnvVars,
	})
	if err != nil {
		// Update deployment status
		depl.Status = fmt.Sprintf("failed: %s", err.Error())
		_, err := s.DatabaseClient.UpdateDeployment(depl)

		log.Debug().Err(err).Msg("error creating container")
		return nil, errors.New(messages.InternalServerError)
	}
	// Update deployment status
	depl.Status = newContainer.GetStatus()
	depl.ContainerID = newContainer.GetContainerId()
	go func() {
		// Wait for container to be ready or failed
		for {
			<-time.After(5 * time.Second)
			containerInfo, err := s.ContainerServiceClient.GetContainer(context.Background(), &container.GetContainerRequest{
				ContainerId: newContainer.GetContainerId(),
			})
			if err != nil {
				log.Debug().Err(err).Msg("error getting container")
				// continue to wait
				continue
			}
			if containerInfo.GetStatus() == "running" || strings.Contains(containerInfo.GetStatus(), "failed") || containerInfo.GetStatus() == "exited" {
				depl.Status = containerInfo.GetStatus()
				_, err := s.DatabaseClient.UpdateDeployment(depl)
				if err != nil {
					log.Debug().Err(err).Msg("error updating deployment")
				}
				break
			}
		}
	}()
	return depl.ToAPI(), nil
}
