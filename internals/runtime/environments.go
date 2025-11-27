package runtime

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"fmt"
	"log/slog"
)

func applyEnvironment(env config.EnvironmentConfig) error {
	slog.Info("Applying environment", "name", env.Name)
	for _, service := range env.Services {
		engine, err := engines.Create(service.Engine, service, service.EngineConfig)
		if err != nil {
			return err
		}

		_, err = engine.Start()
		if err != nil {
			return err
		}
	}

	slog.Info("Environment applied", "name", env.Name)
	return nil
}

type ApplyConfig struct {
	BaseConfig
	Environment string
}

func ApplyEnvironmentCmd(applyConfig ApplyConfig) error {
	cnfg, err := loadConfig(applyConfig.BaseConfig)
	if err != nil {
		return err
	}

	for _, env := range cnfg.Environments {
		if env.Name == applyConfig.Environment {
			return applyEnvironment(env)
		}
	}

	return fmt.Errorf("environment %s not found", applyConfig.Environment)

}
