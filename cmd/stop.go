package cmd

import (
	"anchor/internals/runtime"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop <environment>",
	Short: "Close the environment and stop its services",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execPath, err := os.Getwd()
		if err != nil {
			slog.Error(err.Error())
			return
		}

		applyConfig := runtime.EnvironmentStatusOptions{
			Environment: args[0],
			AnchorLoaderConfig: runtime.LoadingConfig{
				LoaderName: "hcl",
				Path:       execPath,
			},
			StateLoaderConfig: runtime.LoadingConfig{
				LoaderName: "json",
				Path:       execPath,
			},
		}

		err = runtime.StopEnvironmentCmd(applyConfig)
		if err != nil {
			slog.Error(err.Error())
		}

	},
}
