package runtime

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"anchor/internals/state"
	"fmt"
	"log/slog"
	"path/filepath"
)

type EnvironmentStatusOptions struct {
	AnchorLoaderConfig LoadingConfig
	StateLoaderConfig  LoadingConfig
	Environment        string
}

type environmentModifier func(cnfg config.EnvironmentConfig, st state.State) (*state.State, error)

func modifyEnvironment(modifier environmentModifier, envConfig EnvironmentStatusOptions) error {
	cnfg, err := loadConfig(envConfig.AnchorLoaderConfig)
	if err != nil {
		return err
	}

	statePath := filepath.Join(envConfig.StateLoaderConfig.Path, stateFilename)
	stateLoader, err := state.NewLoader(envConfig.StateLoaderConfig.LoaderName)
	if err != nil {
		return err
	}

	loadedState, err := stateLoader.Load(statePath)
	if err != nil {
		return err
	}

	for _, env := range cnfg.Environments {
		if env.Name == envConfig.Environment {
			newState, err := modifier(env, *loadedState)
			if err != nil {
				return err
			}
			err = stateLoader.Write(statePath, *newState)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("environment %s not found", envConfig.Environment)
}

func getOrCreateServiceState(globalState *state.State, envName string, serviceName string, engineType string) (*state.ServiceState, error) {
	engineState, ok := globalState.GetServiceState(envName, serviceName)
	if !ok {
		newEngineState, err := engines.DefaultState(engineType)
		if err != nil {
			return nil, err
		}
		return &newEngineState, nil
	}
	return engineState, nil
}

func applyEnvironment(env config.EnvironmentConfig, globalState state.State) (*state.State, error) {
	slog.Info("applying environment", "name", env.Name)

	for _, service := range env.Services {
		engineState, err := getOrCreateServiceState(&globalState, env.Name, service.Name, service.Engine)
		if err != nil {
			return nil, err
		}

		engine, err := engines.Create(service.Engine, service, service.EngineConfig, *engineState)
		if err != nil {
			return nil, err
		}

		engineNewState, err := engine.Start()
		globalState.AddServiceState(env.Name, service.Name, engineNewState)

		if err != nil {
			return nil, err
		}
	}

	slog.Info("environment applied", "name", env.Name)
	return &globalState, nil
}

func stopEnvironment(env config.EnvironmentConfig, globalState state.State) (*state.State, error) {
	slog.Info("stopping environment", "name", env.Name)
	for _, service := range env.Services {
		engineState, err := getOrCreateServiceState(&globalState, env.Name, service.Name, service.Engine)
		if err != nil {
			return nil, err
		}

		engine, err := engines.Create(service.Engine, service, service.EngineConfig, *engineState)
		if err != nil {
			return nil, err
		}

		engineStopState, err := engine.Stop()
		globalState.AddServiceState(env.Name, service.Name, engineStopState)
		if err != nil {
			return nil, err
		}
	}
	slog.Info("environment stopped", "name", env.Name)
	return &globalState, nil
}
