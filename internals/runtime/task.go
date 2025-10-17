package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Task interface {
	Run() error
}

type ShellTask struct {
	Command string
}

func (t *ShellTask) Run() error {
	args := strings.Split(t.Command, " ")
	var stderr bytes.Buffer

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.String() != "" {
			return err
		}
		return fmt.Errorf(stderr.String())
	}

	return nil

}
