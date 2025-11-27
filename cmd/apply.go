package cmd

import (
	"anchor/internals/runtime"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply <environment>",
	Short: "Apply the environment and start the required services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		execPath, err := os.Getwd()
		if err != nil {
			return err
		}

		applyConfig := runtime.ApplyConfig{
			Environment: args[0],
			LoaderName:  "hcl",
			Path:        execPath,
		}

		return runtime.ApplyEnvironmentCmd(applyConfig)

	},
}
