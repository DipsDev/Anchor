package runtime

import (
	"anchor/internals/parser"
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"os/exec"
	"strings"
	"time"
)

func (r *Runtime) runTask(task parser.TaskBlock) error {
	s := spinner.New(spinner.CharSets[26], 500*time.Millisecond)
	s.Prefix = fmt.Sprintf("running [%s] ", task.Name)
	s.Start()

	var stderrBuf bytes.Buffer

	args := strings.Split(task.Exec, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = &stderrBuf

	err := cmd.Start()
	if err != nil {
		s.Stop()
		fmt.Printf("running [%s] ... failed\n", task.Name)
		return &TaskError{
			TaskName: task.Name,
			Command:  task.Exec,
			Err:      err,
			Output:   stderrBuf.String(),
		}
	}

	err = cmd.Wait()
	if err != nil {
		s.Stop()
		fmt.Printf("running [%s] ... failed\n", task.Name)
		return &TaskError{
			TaskName: task.Name,
			Command:  task.Exec,
			Err:      err,
			Output:   stderrBuf.String(),
		}
	}

	s.Stop()
	fmt.Printf("running [%s] ... done\n", task.Name)
	return nil
}
