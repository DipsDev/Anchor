package runtime

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"anchor/internals/state"
	"fmt"
	"log/slog"
	"path/filepath"
)

func applyEnvironment(env config.EnvironmentConfig, globalState state.State) (*state.State, error) {
	slog.Info("Applying environment", "name", env.Name)

	for _, service := range env.Services {
		engineState, err := globalState.GetServiceState(env.Name, service.Name)
		if err != nil {
			return nil, err
		}

		engine, err := engines.Create(service.Engine, service, service.EngineConfig, engineState)
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

type ApplyConfig struct {
	AnchorLoaderConfig LoadingConfig
	StateLoaderConfig  LoadingConfig
	Environment        string
}

func ApplyEnvironmentCmd(applyConfig ApplyConfig) error {
	cnfg, err := loadConfig(applyConfig.AnchorLoaderConfig)
	if err != nil {
		return err
	}

	statePath := filepath.Join(applyConfig.StateLoaderConfig.Path, stateFilename)
	stateLoader, err := state.NewLoader(applyConfig.StateLoaderConfig.LoaderName)
	if err != nil {
		return err
	}

	loadedState, err := stateLoader.Load(statePath)
	if err != nil {
		return err
	}

	for _, env := range cnfg.Environments {
		if env.Name == applyConfig.Environment {
			newState, err := applyEnvironment(env, *loadedState)
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

	return fmt.Errorf("environment %s not found", applyConfig.Environment)

}
