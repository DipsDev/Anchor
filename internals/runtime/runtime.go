package runtime

import (
	"anchor/internals/parser"
	"fmt"
)

type Runtime struct {
	config *parser.RootConfig
}

func New(config *parser.RootConfig) *Runtime {
	return &Runtime{
		config: config,
	}
}

func (r *Runtime) StartEnvironment(targetEnv string) error {
	env := r.config.GetEnvironment(targetEnv)
	if env == nil {
		return fmt.Errorf("environment '%s' not found", targetEnv)
	}

	for _, task := range env.Tasks {
		if err := r.runTask(task); err != nil {
			return err
		}
	}

	return nil
}
