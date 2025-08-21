package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "harborArk",
	Short: "harborArk is a NAS application with Gin and Cobra",
	Long: `harborArk is a NAS application that combines a Gin web server 
with Cobra CLI functionality.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error executing command:", err)
		os.Exit(1)
	}
}
