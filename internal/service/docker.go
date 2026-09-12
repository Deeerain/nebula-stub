package service

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

type DockerContainer struct {
	Name  string
	Satus string
}

type DockerService interface {
	GetContainerList(ctx context.Context) ([]DockerContainer, error)
	Close()
}

type dockerService struct {
	apiClient *client.Client
}

func (dc *dockerService) GetContainerList(ctx context.Context) ([]DockerContainer, error) {
	var result []DockerContainer

	clientListResult, err := dc.apiClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %w", err)
	}

	result = make([]DockerContainer, 0)

	for _, container := range clientListResult.Items {
		result = append(result, DockerContainer{
			Name:  container.Names[0],
			Satus: container.Status,
		})
	}

	return result, nil
}

func (dc *dockerService) Close() {
	dc.apiClient.Close()
}

func NewDockerService() (DockerService, error) {
	var err error

	dockerApi, err := client.New(client.FromEnv)

	return &dockerService{
		apiClient: dockerApi,
	}, err
}
