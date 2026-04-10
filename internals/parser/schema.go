package parser

type RootConfig struct {
	Version string `hcl:"version,optional"`

	Environments []EnvironmentDecl `hcl:"environment,block"`
}

func (r *RootConfig) GetEnvironment(envName string) *EnvironmentDecl {
	for _, e := range r.Environments {
		if e.Name == envName {
			return &e
		}
	}
	return nil
}

type EnvironmentDecl struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,optional"`

	Tasks []TaskBlock `hcl:"task,block"`
}

func (e *EnvironmentDecl) FindRunnable(name string) Runnable {
	for _, t := range e.Tasks {
		if t.Name == name {
			return &t
		}
	}

	return nil
}

type Runnable interface {
	GetDependencies() []string
	GetName() string
}

type TaskBlock struct {
	Name      string   `hcl:"name,label"`
	Exec      string   `hcl:"exec"`
	DependsOn []string `hcl:"depends_on,optional"`
}

func (t *TaskBlock) GetDependencies() []string {
	return t.DependsOn
}

func (t *TaskBlock) GetName() string {
	return t.Name
}
