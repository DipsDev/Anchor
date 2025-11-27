package runtime

import (
	"anchor/internals/config"
	"path/filepath"
)

const CONFIG_FILENAME = "Anchorfile"

type BaseConfig struct {
	LoaderName string
	Path       string
}

func loadConfig(rConfig BaseConfig) (*config.Config, error) {
	loader, err := config.CreateLoader(rConfig.LoaderName)
	if err != nil {
		return nil, err
	}

	cnfg, err := loader.Load(filepath.Join(rConfig.Path, CONFIG_FILENAME))
	if err != nil {
		return nil, err
	}

	return cnfg, nil
}
