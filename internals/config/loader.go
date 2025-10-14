package config

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

	// Engine dependent
	Image   *string `hcl:"image"`
	Command *string `hcl:"command"`
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
