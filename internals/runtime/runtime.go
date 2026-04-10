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

	runnables := make([]parser.Runnable, 0)

	for _, t := range env.Tasks {
		runnables = append(runnables, &t)
	}

	execOrder, err := resolveExecutionOrder(runnables)
	if err != nil {
		return err
	}

	for _, runnableName := range execOrder {
		runnable := env.FindRunnable(runnableName)
		if err := r.run(runnable); err != nil {
			return err
		}
	}

	return nil
}
