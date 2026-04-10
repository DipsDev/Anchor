package runtime

import "fmt"

type TaskError struct {
	TaskName string
	Command  string
	Err      error
	Output   string
}

func (e *TaskError) Error() string {
	if e.Output != "" {
		return fmt.Sprintf("task '%s' failed:\n  command: %s\n  output: %v", e.TaskName, e.Command, e.Err)
	}

	return fmt.Sprintf("task '%s' failed: %v", e.TaskName, e.Err)
}
