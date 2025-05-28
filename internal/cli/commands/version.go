// Galaxy Node Pool - Version Command
// AI-ID: CP-GAL-NODEPOOL-001
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Galaxy Node Pool\n")
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Build Date: %s\n", BuildDate)
			fmt.Printf("Git Commit: %s\n", GitCommit)
			fmt.Printf("Git Branch: %s\n", GitBranch)
			fmt.Printf("AI-ID:      CP-GAL-NODEPOOL-001\n")
		},
	}
}
