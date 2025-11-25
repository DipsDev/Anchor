package engines

type EngineExecutionResult string

type AnchorEngine[T any] interface {
	Parse(service &ServiceConfig) (T, error)
	Execute() error
}

func Create(engineType string) (AnchorEngine, error) {
	return nil, fmt.Errorf("wrong engine type provided, go engine '%v'", engineType)
}