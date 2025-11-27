package runtime

import (
	"anchor/internals/config"
	"anchor/internals/loader"
	"path/filepath"
)

const ConfigFilename = "Anchorfile"

type BaseConfig struct {
	LoaderName string
	Path       string
}

func loadConfig(rConfig BaseConfig) (*config.Config, error) {
	configLoader, err := loader.CreateLoader(rConfig.LoaderName)
	if err != nil {
		return nil, err
	}

	cnfg, err := configLoader.Load(filepath.Join(rConfig.Path, ConfigFilename))
	if err != nil {
		return nil, err
	}

	return cnfg, nil
}
