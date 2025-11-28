package engines

import (
	"anchor/internals/config"
	"anchor/internals/state"
	"context"
	"fmt"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"log/slog"
	"math/rand"
	"net/netip"
	"strings"
	"time"
)

const containerWaitingTimeout = 30 * time.Second

type DockerEngineConfig struct {
	Image   string   `hcl:"image"`
	Network *string  `hcl:"network,optional"`
	Env     []string `hcl:"env,optional"`
	Ports   []string `hcl:"ports,optional"`
}

type dockerEngineState struct {
	Pid string
}

func createDefaultDockerState() state.ServiceState {
	return dockerEngineState{Pid: ""}
}

type DockerEngine struct {
	config        DockerEngineConfig
	serviceConfig config.ServiceConfig
	state         dockerEngineState
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
	interval := time.NewTicker(4 * time.Second)
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

func newDocker(serviceConfig config.ServiceConfig, config config.EngineConfig, serviceState state.ServiceState) Engine {
	return DockerEngine{
		config:        *config.(*DockerEngineConfig),
		serviceConfig: serviceConfig,
		state:         serviceState.(dockerEngineState),
	}
}

func generatePortMap(mapping []string) (network.PortMap, error) {
	mp := make(network.PortMap)
	for _, strPort := range mapping {
		ports := strings.Split(strPort, ":")
		hostPort := ports[0]
		containerPort, err := network.ParsePort(ports[1])
		if err != nil {
			return nil, err
		}

		portBindings := make([]network.PortBinding, 0)
		binding := network.PortBinding{
			HostIP:   netip.Addr{},
			HostPort: hostPort,
		}
		mp[containerPort] = append(portBindings, binding)
	}
	return mp, nil
}

func (de DockerEngine) Start() (state.ServiceState, error) {
	conn, err := de.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.client.Close()

	// if the container is already created, just run it
	if de.state.Pid != "" {
		slog.Info("starting container", "id", de.state.Pid)
		_, err = conn.client.ContainerStart(conn.ctx, de.state.Pid, client.ContainerStartOptions{})
		if err != nil {
			return nil, err
		}

		slog.Info("waiting for container to run", "id", de.state.Pid, "timeout", containerWaitingTimeout)
		err = waitForContainer(conn, de.state.Pid, func(inspectResult client.ContainerInspectResult) bool {
			return inspectResult.Container.State.Running
		})
		if err != nil {
			return nil, err
		}

		return de.state, nil
	}

	slog.Info("pulling docker image", "image", de.config.Image)
	pullResponse, err := conn.client.ImagePull(conn.ctx, de.config.Image, client.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	err = pullResponse.Wait(conn.ctx)
	if err != nil {
		return nil, err
	}

	portBindings, err := generatePortMap(de.config.Ports)
	if err != nil {
		return nil, err
	}

	resp, err := conn.client.ContainerCreate(conn.ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Env: de.config.Env,
		},
		HostConfig: &container.HostConfig{
			PortBindings: portBindings,
		},
		NetworkingConfig: &network.NetworkingConfig{},
		Platform:         nil,
		Name:             fmt.Sprintf("anchor-%s%d", de.serviceConfig.Name, rand.Intn(300)),
		Image:            de.config.Image,
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

	return &dockerEngineState{
		Pid: resp.ID,
	}, nil
}

func (de DockerEngine) Stop() (state.ServiceState, error) {
	conn, err := de.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.client.Close()

	slog.Info("stopping docker container", "id", de.state.Pid)
	_, err = conn.client.ContainerStop(conn.ctx, de.state.Pid, client.ContainerStopOptions{})
	if err != nil {
		return nil, err
	}

	slog.Info("waiting for container to stop", "id", de.state.Pid, "timeout", containerWaitingTimeout)
	err = waitForContainer(conn, de.state.Pid, func(inspectResult client.ContainerInspectResult) bool {
		return !inspectResult.Container.State.Running
	})
	if err != nil {
		return nil, err
	}

	return de.state, nil
}
