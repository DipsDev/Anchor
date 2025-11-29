package engines

import (
	"anchor/internals/config"
	"anchor/internals/state"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Engine interface {
	Start() (state.ServiceState, error)
	Stop() (state.ServiceState, error)
}

type engineDefinition struct {
	engineFactory func(serviceConfig config.ServiceConfig, config config.EngineConfig, state state.ServiceState) Engine
	configFactory func() config.EngineConfig
	stateFactory  func() state.ServiceState
}

var engines = map[string]engineDefinition{
	"docker": {
		engineFactory: newDocker,
		configFactory: func() config.EngineConfig { return &DockerEngineConfig{} },
		stateFactory:  createDefaultDockerState,
	},
	"shell": {
		engineFactory: newShellEngine,
		configFactory: func() config.EngineConfig { return &ShellEngineConfig{} },
		stateFactory:  createDefaultShellState,
	},
}

func Create(engineType string, serviceConfig config.ServiceConfig, config config.EngineConfig, state state.ServiceState) (Engine, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	convertedState, _ := DefaultState(engineType)
	err := mapstructure.Decode(state, &convertedState)
	if err != nil {
		return nil, err
	}

	return engine.engineFactory(serviceConfig, config, convertedState), nil
}

func Config(engineType string) (config.EngineConfig, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.configFactory(), nil
}

func DefaultState(engineType string) (state.ServiceState, error) {
	engine, ok := engines[engineType]
	if !ok {
		return nil, fmt.Errorf("wrong engine type provided, engine '%v' is not defined\n", engineType)
	}

	return engine.stateFactory(), nil

}
