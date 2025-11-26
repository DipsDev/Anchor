package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "anchor",
	Short: "Anchor is a powerful workflow tool to streamline development with ease",
	Long:  "A powerful workflow tool built to impower you dev experience.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
