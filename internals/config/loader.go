package config

type Config struct {
	Environments []Environment `hcl:"environment,block"`
}
type Environment struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description"`

	// Processes
	Services []Service `hcl:"service,block"`
	Tasks    []Task    `hcl:"task,block"`
}

type Service struct {
	Name        string      `hcl:"name,label"`
	Engine      string      `hcl:"engine"`
	HealthCheck HealthCheck `hcl:"health_check,block"`
	DependsOn   []string    `hcl:"depends_on,optional"`

	// Engine dependent
	Image   *string `hcl:"image"`
	Command *string `hcl:"command"`
}

type Task struct {
	Name      string   `hcl:"name,label"`
	Command   string   `hcl:"command"`
	DependsOn []string `hcl:"depends_on,optional"`
}

type HealthCheck struct {
	ConnectionType string `hcl:"type"`
	Target         string `hcl:"target"`

	// Optional
	Timeout *string `hcl:"timeout"`
}

type Loader interface {
	Load(path string) (*Config, error)
}
