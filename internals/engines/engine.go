package engines

import (
	"fmt"
)

type EngineExecutionResult string
type EngineConfig interface{}

type Engine interface {
	Start() (EngineExecutionResult, error)
	Stop() (EngineExecutionResult, error)
}

type engineDefinition struct {
	engineFactory func(config EngineConfig) Engine
	configFactory func() EngineConfig
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: func(config EngineConfig) Engine { return NewDockerEngine(config) },
		configFactory: func() EngineConfig { return &DockerEngineConfig{} },
	},
}

func Create(engineType string, config EngineConfig) (Engine, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.engineFactory(config), nil
}

func Config(engineType string) (EngineConfig, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.configFactory(), nil
}
