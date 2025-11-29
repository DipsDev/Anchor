package runtime

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"anchor/internals/state"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
)

type EnvironmentStatusOptions struct {
	AnchorLoaderConfig LoadingConfig
	StateLoaderConfig  LoadingConfig
	Environment        string
}

type environmentModifier func(cnfg config.EnvironmentConfig, st state.State) (*state.State, error)

// region Environment utilities
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

// region execution order
func generateExectionParameters(env config.EnvironmentConfig, mappedServices map[string]config.ServiceConfig) (map[string]int, map[string][]string, error) {
	indegress := make(map[string]int)
	graph := make(map[string][]string)

	for _, service := range env.Services {
		indegress[service.Name] += len(service.DependsOn)
		for _, dependency := range service.DependsOn {
			_, ok := mappedServices[dependency]
			if !ok {
				return nil, nil, fmt.Errorf("dependency `%s` is not defined", dependency)
			}
			graph[dependency] = append(graph[dependency], service.Name)
		}
	}

	return indegress, graph, nil
}

func createMappedServiceConfigs(env config.EnvironmentConfig) map[string]config.ServiceConfig {
	res := make(map[string]config.ServiceConfig)
	for _, service := range env.Services {
		res[service.Name] = service
	}
	return res
}

func getEnvironmentExecutionOrder(env config.EnvironmentConfig) ([]config.ServiceConfig, error) {
	mappedServices := createMappedServiceConfigs(env)
	indegress, graph, err := generateExectionParameters(env, mappedServices)
	if err != nil {
		return nil, err
	}

	stack := make([]string, 0)
	res := make([]config.ServiceConfig, 0)

	for _, service := range env.Services {
		if indegress[service.Name] == 0 {
			stack = append(stack, service.Name)
		}
	}

	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]

		service, ok := mappedServices[node]
		if !ok {
			return res, fmt.Errorf("service %s not found", node)
		}

		res = append(res, service)

		for _, dependency := range graph[node] {
			indegress[dependency]--
			if indegress[dependency] == 0 {
				stack = append(stack, dependency)
			}
		}
	}

	if len(res) != len(env.Services) {
		return res, fmt.Errorf("environment %s has circular-dependencies, please check the `depends_on` attribute for typos and mismatches", env.Name)
	}

	return res, nil
}

//endregion

//endregion

func applyEnvironment(env config.EnvironmentConfig, globalState state.State) (*state.State, error) {
	slog.Info("applying environment", "name", env.Name)

	services, err := getEnvironmentExecutionOrder(env)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		engineState, err := getOrCreateServiceState(&globalState, env.Name, service.Name, service.Engine)
		if err != nil {
			return nil, err
		}

		engine, err := engines.Create(service.Engine, service, service.EngineConfig, *engineState)
		if err != nil {
			return nil, err
		}

		slog.Info("starting service", "name", service.Name)
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

	services, err := getEnvironmentExecutionOrder(env)
	if err != nil {
		return nil, err
	}

	slices.Reverse(services)
	for _, service := range services {
		engineState, err := getOrCreateServiceState(&globalState, env.Name, service.Name, service.Engine)
		if err != nil {
			return nil, err
		}

		engine, err := engines.Create(service.Engine, service, service.EngineConfig, *engineState)
		if err != nil {
			return nil, err
		}

		slog.Info("stopping service", "name", service.Name)
		engineStopState, err := engine.Stop()
		globalState.AddServiceState(env.Name, service.Name, engineStopState)
		if err != nil {
			return nil, err
		}
	}
	slog.Info("environment stopped", "name", env.Name)
	return &globalState, nil
}
