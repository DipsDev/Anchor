package state

type EngineState struct {
	Label  string
	Engine string
}

type DockerEngineState struct {
	EngineState
	Status string
}

type EnvironmentState struct {
	Engine []EngineState
}

type State struct {
	Version      int
	Environments []EnvironmentState
}

type Loader interface {
	Load() (*State, error)
	Write(state State) error
}
