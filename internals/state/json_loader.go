package state

import (
	"encoding/json"
	"os"
)

type JsonLoader struct {
	Path string
}

func NewJsonStateLoader(path string) *JsonLoader {
	return &JsonLoader{
		Path: path,
	}
}

func (stl JsonLoader) Load() (*State, error) {
	data, err := os.ReadFile(stl.Path)
	if err != nil {
		return nil, err
	}

	state := &State{}
	marshallErr := json.Unmarshal(data, state)
	if marshallErr != nil {
		return nil, marshallErr
	}

	return state, nil
}

func (stl JsonLoader) Write(state State) error {
	data, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(stl.Path, data, 0644)
}
