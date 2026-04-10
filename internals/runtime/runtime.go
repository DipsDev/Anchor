package runtime

import (
	"anchor/internals/parser"
	"fmt"
	"github.com/briandowns/spinner"
	"os/exec"
	"strings"
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
		return fmt.Errorf("environment '%s' not found", targetEnv)
	}

	for _, task := range env.Tasks {
		s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
		_ = s.Color("gray")
		s.Suffix = " Starting task " + task.Name
		s.Start()

		args := strings.Split(task.Exec, " ")
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Start()
		if err != nil {
			return err
		}

		err = cmd.Wait()
		if err != nil {
			return err
		}

		s.Stop()
		fmt.Println("Started task " + task.Name)
	}

	return nil
}
