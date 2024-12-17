package services

import (
	"context"
	"fmt"
	"microservicios-api/clients"
	"microservicios-api/dto"
	"strings"
)

var dockerClient = clients.NewDockerClient()

func GetServices(ctx context.Context) (dto.ServicesResponse, error) {
	containers, err := dockerClient.GetContainers(ctx)
	if err != nil {
		return dto.ServicesResponse{}, fmt.Errorf("error fetching containers: %w", err)
	}

	expectedServices := map[string]string{
		"search-api": "8082",
		"cursos-api": "8081",
		"users-api":  "8080",
	}

	var services []dto.Service

	for _, container := range containers {
		for serviceName, port := range expectedServices {
			for _, name := range container.Names {
				if strings.Contains(name, serviceName) && container.State == "running" {
					services = append(services, dto.Service{
						Name:      serviceName,
						Container: container.ID,
						Port:      port,
						State:     container.State,
					})
					break
				}
			}
		}
	}

	return dto.ServicesResponse{Services: services}, nil
}
