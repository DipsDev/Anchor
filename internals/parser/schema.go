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

type TaskBlock struct {
	Name string `hcl:"name,label"`
	Exec string `hcl:"exec"`
}
