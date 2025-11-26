package state

type EngineState struct {
	Label  string
	Engine string
}

type DockerEngineState struct {
	EngineState
	Status string
}

type State struct {
	Version      int
	EngineStates []EngineState
}

type Loader interface {
	Load() (*State, error)
	Write(state State) error
}
