package engines

import (
	"anchor/internals/config"
	"context"
	"fmt"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"log/slog"
	"time"
)

const containerWaitingTimeout = 30 * time.Second

type DockerResultState struct {
	Pid string
}

type DockerEngineConfig struct {
	Image   string   `hcl:"image"`
	Network *string  `hcl:"network,optional"`
	Env     []string `hcl:"env,optional"`
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

type waitingContainerPredicate func(inspectResult client.ContainerInspectResult) bool

func waitForContainer(conn *dockerConnection, containerId string, predicate waitingContainerPredicate) error {
	timeout := time.NewTimer(containerWaitingTimeout)
	interval := time.NewTicker(1 * time.Second)
	defer interval.Stop()
	for {
		select {
		case <-timeout.C:
			return fmt.Errorf("timed out waiting for container to start")
		case <-interval.C:
			inspectResult, err := conn.client.ContainerInspect(conn.ctx, containerId, client.ContainerInspectOptions{})
			if err != nil {
				return err
			}
			if predicate(inspectResult) {
				return nil
			}
		}
	}
}

func newDocker(serviceConfig config.ServiceConfig, config config.EngineConfig) DockerEngine {
	dockerConfig := config.(*DockerEngineConfig)
	return DockerEngine{
		Config:        *dockerConfig,
		serviceConfig: serviceConfig,
	}
}

func (de DockerEngine) Start() (EngineResultState, error) {
	conn, err := de.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.client.Close()

	slog.Info("pulling docker image", "image", de.Config.Image)
	pullResponse, err := conn.client.ImagePull(conn.ctx, de.Config.Image, client.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	err = pullResponse.Wait(conn.ctx)
	if err != nil {
		return nil, err
	}

	resp, err := conn.client.ContainerCreate(conn.ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Env: de.Config.Env,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
		Platform:         nil,
		Name:             de.serviceConfig.Name,
		Image:            de.Config.Image,
	})
	if err != nil {
		return nil, err
	}

	slog.Info("starting container", "id", resp.ID)
	_, err = conn.client.ContainerStart(conn.ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	slog.Info("waiting for container to run", "id", resp.ID, "timeout", containerWaitingTimeout)
	err = waitForContainer(conn, resp.ID, func(inspectResult client.ContainerInspectResult) bool {
		return inspectResult.Container.State.Running
	})
	if err != nil {
		return nil, err
	}

	return &DockerResultState{
		Pid: resp.ID,
	}, nil
}

func (de DockerEngine) Stop() (EngineResultState, error) {
	fmt.Println("Stopping docker container", de.Config.Image)
	return nil, nil
}
