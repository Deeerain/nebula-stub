package service

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type DockerContainer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type DockerService interface {
	GetContainerList(ctx context.Context) ([]DockerContainer, error)
	Close()
}

type dockerService struct {
	apiClient *client.Client
}

func (dc *dockerService) GetContainerList(ctx context.Context) ([]DockerContainer, error) {
	if dc.apiClient == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	clientListResult, err := dc.apiClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %w", err)
	}

	result := make([]DockerContainer, 0, len(clientListResult.Items))

	for _, container := range clientListResult.Items {
		name := ""
		if len(container.Names) > 0 {
			name = container.Names[0]
		}

		result = append(result, DockerContainer{
			Name:   name,
			Status: container.Status,
		})
	}

	return result, nil
}

func (dc *dockerService) Close() {
	if dc.apiClient != nil {
		dc.apiClient.Close()
	}
}

func NewDockerService() (DockerService, error) {
	dockerApi, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &dockerService{
		apiClient: dockerApi,
	}, nil
}
