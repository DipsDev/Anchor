package engines

import (
	"anchor/internals/config"
	"fmt"
)

// EngineResultState is interface{} to allow each engine to save what they need
type EngineResultState interface{}

type Engine interface {
	Start() (EngineResultState, error)
	Stop() (EngineResultState, error)
}

type engineDefinition struct {
	engineFactory func(serviceConfig config.ServiceConfig, config config.EngineConfig) Engine
	configFactory func() config.EngineConfig
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: func(serviceConfig config.ServiceConfig, config config.EngineConfig) Engine {
			return newDocker(serviceConfig, config)
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
