package cmd

import (
	"anchor/internals/runtime"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply <environment>",
	Short: "Apply the environment and start the required services",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execPath, err := os.Getwd()
		if err != nil {
			slog.Error(err.Error())
			return
		}

		applyConfig := runtime.ApplyConfig{
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

		err = runtime.ApplyEnvironmentCmd(applyConfig)
		if err != nil {
			slog.Error(err.Error())
		}

	},
}
