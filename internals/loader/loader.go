package loader

import (
	"anchor/internals/config"
	"fmt"
)

type Loader interface {
	Load(path string) (*config.Config, error)
}

type mappedLoader struct {
	CreateLoader func() (Loader, error)
}

var loaders = map[string]mappedLoader{
	"hcl": {
		CreateLoader: func() (Loader, error) {
			hclLoader := NewHclLoader()
			return hclLoader, nil
		},
	},
}

func CreateLoader(loaderName string) (Loader, error) {
	loader, ok := loaders[loaderName]
	if !ok {
		return nil, fmt.Errorf("loader %s does not exist", loaderName)
	}

	return loader.CreateLoader()
}
