package engines

import "fmt"

type DockerEngineConfig struct {
	Image   string `hcl:"image"`
	Network string `hcl:"network"`
}

type DockerEngine struct {
	Config DockerEngineConfig
}

func NewDockerEngine(config EngineConfig) DockerEngine {
	dockerConfig := config.(*DockerEngineConfig)
	return DockerEngine{
		Config: *dockerConfig,
	}
}

func (m DockerEngine) Start() (EngineExecutionResult, error) {
	fmt.Println("docker engine execute")
	return "up", nil
}

func (m DockerEngine) Stop() (EngineExecutionResult, error) {
	fmt.Println("Stopping docker container", m.Config.Image)
	return "down", nil
}
