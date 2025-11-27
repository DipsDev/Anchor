package engines

import (
	"anchor/internals/config"
	"context"
	"fmt"
	"github.com/moby/moby/client"
	"log/slog"
)

type DockerEngineConfig struct {
	Image   string `hcl:"image"`
	Network string `hcl:"network"`
}

type DockerEngine struct {
	Config        DockerEngineConfig
	serviceConfig config.ServiceConfig
}

type dockerConnection struct {
	ctx    context.Context
	client *client.Client
}

func (de DockerEngine) createConnection() (*dockerConnection, error) {
	ctx := context.Background()
	cli, err := client.New()
	if err != nil {
		return nil, err
	}

	conn := &dockerConnection{
		ctx:    ctx,
		client: cli,
	}
	return conn, nil

}

func NewDockerEngine(serviceConfig config.ServiceConfig, config config.EngineConfig) DockerEngine {
	dockerConfig := config.(*DockerEngineConfig)
	return DockerEngine{
		Config:        *dockerConfig,
		serviceConfig: serviceConfig,
	}
}

func (de DockerEngine) Start() (*EngineExecutionResult, error) {
	conn, err := de.createConnection()
	if err != nil {
		return nil, err
	}

	slog.Info("pulling image", "image", de.Config.Image)
	_, err = conn.client.ImagePull(conn.ctx, de.Config.Image, client.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	resp, err := conn.client.ContainerCreate(conn.ctx, client.ContainerCreateOptions{
		Config:           nil,
		HostConfig:       nil,
		NetworkingConfig: nil,
		Platform:         nil,
		Name:             de.serviceConfig.Name,
		Image:            de.Config.Image,
	})
	if err != nil {
		return nil, err
	}

	_, err = conn.client.ContainerStart(conn.ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (de DockerEngine) Stop() (*EngineExecutionResult, error) {
	fmt.Println("Stopping docker container", de.Config.Image)
	return nil, nil
}
