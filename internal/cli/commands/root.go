// Galaxy Node Pool - Root Command
// AI-ID: CP-GAL-NODEPOOL-001
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// These will be set during build
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "none"
	GitBranch = ""
)

// NewRootCmd creates the root command for the Galaxy Pool CLI
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "galaxy-pool",
		Short: "Galaxy Node Pool - Manage your Galaxy Network nodes",
		Long: `Galaxy Node Pool is a high-performance node management system for the Galaxy Network.

This tool allows you to manage nodes, pools, and participate in the Galaxy Network.`,
		Version: fmt.Sprintf("%s\nBuild Date: %s\nCommit: %s\nBranch: %s", 
			Version, BuildDate, GitCommit, GitBranch),
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewRootCmd()
	
	// Add all commands
	addCommands(rootCmd)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// addCommands adds all child commands to the root command
func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(NewVersionCmd())
}
