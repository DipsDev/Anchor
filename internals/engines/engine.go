package engines

import (
	"anchor/internals/config"
	"fmt"
)

// EngineState is interface{} to allow each engine to save what they need
type EngineState interface{}

type Engine interface {
	Start() (EngineState, error)
	Stop() (EngineState, error)
}

type engineDefinition struct {
	engineFactory func(serviceConfig config.ServiceConfig, config config.EngineConfig, state EngineState) Engine
	configFactory func() config.EngineConfig
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: func(serviceConfig config.ServiceConfig, config config.EngineConfig, state EngineState) Engine {
			return newDocker(serviceConfig, config, state)
		},
		configFactory: func() config.EngineConfig { return &DockerEngineConfig{} },
	},
}

func Create(engineType string, serviceConfig config.ServiceConfig, config config.EngineConfig, state EngineState) (Engine, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.engineFactory(serviceConfig, config, state), nil
}

func Config(engineType string) (config.EngineConfig, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.configFactory(), nil
}
