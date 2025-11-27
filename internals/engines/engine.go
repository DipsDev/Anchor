package engines

import (
	"anchor/internals/config"
	"fmt"
)

type EngineExecutionResult string

type Engine interface {
	Start() (*EngineExecutionResult, error)
	Stop() (*EngineExecutionResult, error)
}

type engineDefinition struct {
	engineFactory func(serviceConfig config.ServiceConfig, config config.EngineConfig) Engine
	configFactory func() config.EngineConfig
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: func(serviceConfig config.ServiceConfig, config config.EngineConfig) Engine {
			return NewDockerEngine(serviceConfig, config)
		},
		configFactory: func() config.EngineConfig { return &DockerEngineConfig{} },
	},
}

func Create(engineType string, serviceConfig config.ServiceConfig, config config.EngineConfig) (Engine, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.engineFactory(serviceConfig, config), nil
}

func Config(engineType string) (config.EngineConfig, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.configFactory(), nil
}
