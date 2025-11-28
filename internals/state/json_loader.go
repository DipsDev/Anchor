package state

import (
	"encoding/json"
	"log/slog"
	"os"
)

type JsonLoader struct{}

func newJsonLoader() *JsonLoader {
	return &JsonLoader{}
}

func (stl JsonLoader) Load(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Warn("encountered an error while trying to read statefile, ignoring...", "path", path, "error", err)
		return &State{}, nil
	}

	state := &State{}
	marshallErr := json.Unmarshal(data, state)
	if marshallErr != nil {
		return nil, marshallErr
	}

	return state, nil
}

func (stl JsonLoader) Write(path string, state State) error {
	data, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
