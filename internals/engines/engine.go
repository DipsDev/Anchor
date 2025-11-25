package engines

import (
	cnfg "anchor/internals/config"
	"fmt"
)

type EngineExecutionResult string

type AnchorEngine interface {
	Parse(service cnfg.ServiceConfig) error
	Execute() error
}

func Create(engineType string) (AnchorEngine, error) {
	return nil, fmt.Errorf("wrong engine type provided, go engine '%v'", engineType)
}
