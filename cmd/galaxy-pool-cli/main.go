// Galaxy Node Pool CLI
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// These will be set during build time
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "galaxy-pool-cli",
		Short: "Galaxy Node Pool CLI - Manage your Galaxy Network nodes and pools",
		Long: `Galaxy Node Pool CLI provides tools to manage nodes and pools in the Galaxy Network.

This is the command line interface for interacting with the Galaxy Node Pool system.`,
		Version: fmt.Sprintf("%s\nCommit: %s\nBuild Time: %s", Version, Commit, BuildTime),
	}

	// Add commands
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewNodeCommand())
	rootCmd.AddCommand(NewPoolCommand())
	rootCmd.AddCommand(NewStakeCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
