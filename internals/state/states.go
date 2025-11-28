package state

import (
	"fmt"
)

type ServiceState interface{}

type EnvironmentState struct {
	Services map[string]ServiceState
}

type State struct {
	Version      int
	Environments map[string]EnvironmentState
}

func generateEmptyEnvironmentState() EnvironmentState {
	return EnvironmentState{
		Services: make(map[string]ServiceState),
	}
}

// AddServiceState updates or creates a service under the given env name
func (st *State) AddServiceState(env string, serviceName string, service ServiceState) {
	if st.Environments == nil {
		st.Environments = make(map[string]EnvironmentState)
	}
	envState, ok := st.Environments[env]
	if !ok {
		st.Environments[env] = generateEmptyEnvironmentState()
		envState = st.Environments[env]
	}

	envState.Services[serviceName] = service

}

func (st *State) GetServiceState(env string, serviceName string) (*ServiceState, error) {
	envState, ok := st.Environments[env]
	if !ok {
		return nil, fmt.Errorf("environment %s does not exist", env)
	}
	serviceState, ok := envState.Services[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s does not exist", serviceName)
	}
	return &serviceState, nil
}

type mappedLoader struct {
	NewLoader func() Loader
}

var loaders = map[string]mappedLoader{
	"json": {
		NewLoader: func() Loader {
			return newJsonLoader()
		},
	},
}

type Loader interface {
	Load(path string) (*State, error)
	Write(path string, state State) error
}

func NewLoader(loaderName string) (Loader, error) {
	l, ok := loaders[loaderName]
	if !ok {
		return nil, fmt.Errorf("no loader with name %q", loaderName)
	}
	return l.NewLoader(), nil
}
