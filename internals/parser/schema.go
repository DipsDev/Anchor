package parser

type RootConfig struct {
	Version string `hcl:"version,optional"`

	Environments []EnvironmentDecl `hcl:"environment,block"`
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
