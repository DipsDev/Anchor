package parser

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"os"
)

type Parser struct {
	hclparser *hclparse.Parser
}

func New() *Parser {
	return &Parser{hclparse.NewParser()}
}

func (p *Parser) ParseFile(path string) (*RootConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, diags := p.hclparser.ParseHCL(data, path)
	if diags.HasErrors() {
		return nil, diags
	}

	var config RootConfig
	diags = gohcl.DecodeBody(f.Body, nil, &config)
	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil

}
