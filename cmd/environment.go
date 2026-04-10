package cmd

import (
	"anchor/internals/parser"
	"anchor/internals/runtime"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var envStartCmd = &cobra.Command{
	Use:   "start <environment>",
	Short: "Start an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		currDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unexpected error occured: %v", err)
		}

		configFile := filepath.Join(currDir, "Anchorfile")
		parser := parser.New()

		config, err := parser.ParseFile(configFile)
		if err != nil {
			return fmt.Errorf("parsing Anchorfile failed: %v", err)
		}

		r := runtime.New(config)
		return r.StartEnvironment(args[0])

	},
}
