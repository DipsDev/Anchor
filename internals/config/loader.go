package config

type Config struct {
	Environments []Environment `hcl:"environment,block"`
}
type Environment struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description"`
}

type ProcessBlock struct {
	DependsOn []string `hcl:"depends_on"`
}

type Service struct {
	ProcessBlock
	Name   string `hcl:"name,label"`
	Engine string `hcl:"engine"`
}

type Task struct {
	ProcessBlock
	Name    string `hcl:"name,label"`
	Command string `hcl:"command"`
}

type HealthCheck struct {
	ConnectionType string `hcl:"type"`
	Target         string `hcl:"target"`

	// Optional
	Timeout *string `hcl:"timeout,optional"`
}

type Loader interface {
	Load() (*Config, error)
}
