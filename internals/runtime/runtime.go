package runtime

import (
	"anchor/internals/config"
	"fmt"
	"path/filepath"
)

const CONFIG_FILENAME = "Anchorfile"

type ApplyConfig struct {
	Environment string
	LoaderName  string
	Path        string
}

func ApplyEnvironmentCmd(applyConfig ApplyConfig) error {
	loader, err := config.CreateLoader(applyConfig.LoaderName)
	if err != nil {
		return err
	}

	cnfg, err := loader.Load(filepath.Join(applyConfig.Path, CONFIG_FILENAME))
	if err != nil {
		return err
	}

	for _, env := range cnfg.Environments {
		if env.Name == applyConfig.Environment {
			return applyEnvironment(env)
		}
	}

	return fmt.Errorf("environment %s not found", applyConfig.Environment)

}
