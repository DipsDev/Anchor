package config

import (
	"anchor/internals/engines"
	"fmt"
	"github.com/hashicorp/hcl/v2"
)

type Config struct {
	Environments []EnvironmentConfig `hcl:"environment,block"`
}
type EnvironmentConfig struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description"`

	// Processes
	Services []ServiceConfig `hcl:"service,block"`
	Tasks    []TaskConfig    `hcl:"task,block"`
}

type ServiceConfig struct {
	Name        string            `hcl:"name,label"`
	Engine      string            `hcl:"engine"`
	HealthCheck HealthCheckConfig `hcl:"health_check,block"`
	DependsOn   []string          `hcl:"depends_on,optional"`

	// concrete engine config
	EngineConfig engines.EngineConfig

	// loader related config
	HclEngineConfig hcl.Body `hcl:",remain"`
}

type TaskConfig struct {
	Name      string   `hcl:"name,label"`
	Command   string   `hcl:"command"`
	DependsOn []string `hcl:"depends_on,optional"`
}

type HealthCheckConfig struct {
	ConnectionType string `hcl:"type"`
	Target         string `hcl:"target"`

	// Optional
	Timeout *string `hcl:"timeout"`
}

type Loader interface {
	Load(path string) (*Config, error)
}

type mappedLoader struct {
	CreateLoader func() (Loader, error)
}

var loaders = map[string]mappedLoader{
	"hcl": {
		CreateLoader: func() (Loader, error) {
			hclLoader := NewHclLoader()
			return hclLoader, nil
		},
	},
}

func CreateLoader(loaderName string) (Loader, error) {
	loader, ok := loaders[loaderName]
	if !ok {
		return nil, fmt.Errorf("loader %s does not exist", loaderName)
	}

	return loader.CreateLoader()
}
