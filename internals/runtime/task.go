package runtime

import (
	"bytes"
	"os/exec"
)

type Task interface {
	Run() (stderr string, err error)
}

type ShellTask struct {
	Command string
}

func (t *ShellTask) Run() (stderr string, err error) {
	var stderrBytes bytes.Buffer
	cmd := exec.Command(t.Command)
	cmd.Stderr = &stderrBytes
	err = cmd.Run()
	if err != nil {
		return stderrBytes.String(), err
	}
	return stderrBytes.String(), nil

}
