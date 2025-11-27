package runtime

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"log/slog"
)

func applyEnvironment(env config.EnvironmentConfig) error {
	slog.Info("Applying environment", "name", env.Name)
	for _, service := range env.Services {
		engine, err := engines.Create(service.Engine, service.EngineConfig)
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
