package runtime

import (
	"anchor/internals/parser"
	"fmt"
	"github.com/briandowns/spinner"
	"time"
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
		return fmt.Errorf("couldn't start environment, environment '%s' not found", targetEnv)
	}

	for _, task := range env.Tasks {
		s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
		_ = s.Color("gray")
		s.Suffix = " Starting task " + task.Name
		s.Start()
		time.Sleep(4 * time.Second)
		s.Stop()
		fmt.Println("Started task " + task.Name)
	}

	return nil
}
