package engines

import (
	"anchor/internals/config"
	"anchor/internals/state"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

type ShellEngineConfig struct {
	Command string `hcl:"command"`
}

type shellEngineState struct {
	Pid int
}

func createDefaultShellState() state.ServiceState {
	return shellEngineState{Pid: 0}
}

type ShellEngine struct {
	config        ShellEngineConfig
	serviceConfig config.ServiceConfig
	state         shellEngineState
}

func newShellEngine(serviceConfig config.ServiceConfig, config config.EngineConfig, state state.ServiceState) Engine {
	return ShellEngine{
		config:        *config.(*ShellEngineConfig),
		serviceConfig: serviceConfig,
		state:         state.(shellEngineState),
	}
}

func isProcessRunning(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, nil
	}

	return true, nil

}

func (s ShellEngine) Start() (state.ServiceState, error) {
	slog.Info("executing shell command", "command", s.config.Command)
	if s.state.Pid > 0 {
		slog.Info("pid found, checking process status", "pid", s.state.Pid)
		running, err := isProcessRunning(s.state.Pid)
		if err != nil {
			return nil, err
		}
		if running {
			slog.Info("process is running, ignoring", "pid", s.state.Pid)
			return s.state, nil
		}
	}

	cmd := exec.Command(s.config.Command)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return &shellEngineState{Pid: cmd.Process.Pid}, nil
}

func (s ShellEngine) Stop() (state.ServiceState, error) {
	slog.Info("terminating shell command", "command", s.config.Command, "pid", s.state.Pid)
	process, err := os.FindProcess(s.state.Pid)
	if err != nil {
		return nil, err
	}

	err = process.Kill()
	if err != nil {
		return nil, err
	}

	return &shellEngineState{Pid: s.state.Pid}, nil

}
