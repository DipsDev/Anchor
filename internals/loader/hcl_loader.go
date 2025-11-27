package loader

import (
	"anchor/internals/config"
	"anchor/internals/engines"
	"errors"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"os"
)

type HclLoader struct {
}

func (l *HclLoader) Load(path string) (*config.Config, error) {
	var globalConfig config.Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := hclparse.NewParser()

	f, diags := parser.ParseHCL(data, path)
	if diags != nil && diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	diags = gohcl.DecodeBody(f.Body, nil, &globalConfig)
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	for i, env := range globalConfig.Environments {
		for j, service := range env.Services {
			engineConf, engineConfError := engines.Config(service.Engine)
			if engineConfError != nil {
				return nil, engineConfError
			}

			engineDiags := gohcl.DecodeBody(service.HclEngineConfig, nil, engineConf)
			if engineDiags.HasErrors() {
				return nil, errors.New(engineDiags.Error())
			}

			// must use indices because env and service is a copy of the real data
			globalConfig.Environments[i].Services[j].EngineConfig = engineConf
		}
	}

	return &globalConfig, nil
}

func NewHclLoader() *HclLoader {
	return &HclLoader{}
}
