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

func resolveExecutionOrder(runnables []parser.Runnable) ([]string, error) {
	// Map tasks by name for easy lookup
	dependencies := make(map[string]parser.Runnable)
	for _, t := range runnables {
		dependencies[t.GetName()] = t
	}

	inDegree := make(map[string]int)
	adj := make(map[string][]string)

	for _, t := range dependencies {
		for _, dep := range t.GetDependencies() {
			// Check if dependency exists in the current environment
			if _, exists := dependencies[dep]; !exists {
				return nil, fmt.Errorf("failed while resolving dependencies: task `%s` depends on `%s`, which isn't in this environment", t.GetName(), dep)
			}
			adj[dep] = append(adj[dep], t.GetName())
			inDegree[t.GetName()]++
		}
	}

	// Kahn's Algorithm
	var queue []string
	for _, t := range dependencies {
		if inDegree[t.GetName()] == 0 {
			queue = append(queue, t.GetName())
		}
	}

	var result []string
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		result = append(result, curr)

		for _, neighbor := range adj[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(result) != len(runnables) {
		return nil, fmt.Errorf("failed while resolving dependencies: circular dependency detected")
	}

	return result, nil
}

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

func (r *Runtime) run(runnable parser.Runnable) error {
	switch v := runnable.(type) {
	case *parser.TaskBlock:
		return r.runTask(*v)
	default:
		return fmt.Errorf("unknown runnable type: %T", v)
	}

}
