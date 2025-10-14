package config

import (
	"errors"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"os"
)

type HclLoader struct {
}

func (l *HclLoader) Load(path string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := hclparse.NewParser()

	f, diags := parser.ParseHCL(data, path)
	if diags != nil && diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	diags = gohcl.DecodeBody(f.Body, nil, &config)
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	return &config, nil
}

func NewHclLoader() *HclLoader {
	return &HclLoader{}
}
