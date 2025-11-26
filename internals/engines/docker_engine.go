package engines

import "fmt"

type DockerEngineConfig struct {
	Image   string `hcl:"image"`
	Network string `hcl:"network"`
}

type DockerEngine struct{}

func (m *DockerEngine) Init(config EngineConfig) error {
	dockerConfig := config.(*DockerEngineConfig)
	fmt.Printf("docker engine parse %v\n", dockerConfig.Image)
	return nil
}

func (m *DockerEngine) Execute() (EngineExecutionResult, error) {
	fmt.Println("docker engine execute")
	return "result", nil
}
