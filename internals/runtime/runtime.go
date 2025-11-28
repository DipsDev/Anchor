package runtime

import (
	"anchor/internals/config"
	"anchor/internals/loader"
	"path/filepath"
)

const anchorFilename = "Anchorfile"
const stateFilename = ".Anchorstate"

type LoadingConfig struct {
	LoaderName string
	Path       string
}

func loadConfig(rConfig LoadingConfig) (*config.Config, error) {
	configLoader, err := loader.NewLoader(rConfig.LoaderName)
	if err != nil {
		return nil, err
	}

	cnfg, err := configLoader.Load(filepath.Join(rConfig.Path, anchorFilename))
	if err != nil {
		return nil, err
	}

	return cnfg, nil
}
