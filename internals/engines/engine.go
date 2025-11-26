package engines

import (
	"fmt"
)

type EngineExecutionResult string
type EngineConfig interface{}

type Engine interface {
	Init(engineConfig EngineConfig) error
	Execute() (EngineExecutionResult, error)
}

type engineDefinition struct {
	engineFactory func() Engine
	configFactory func() EngineConfig
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: func() Engine { return &DockerEngine{} },
		configFactory: func() EngineConfig { return &DockerEngineConfig{} },
	},
}

func Create(engineType string) (Engine, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.engineFactory(), nil
}

func Config(engineType string) (EngineConfig, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.configFactory(), nil
}
