package config

import "github.com/hashicorp/hcl/v2/hclsimple"

type HclLoader struct {
}

func (l *HclLoader) Load(path string) (*Config, error) {
	var config *Config
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewHclLoader() *HclLoader {
	return &HclLoader{}
}
