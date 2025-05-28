package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Galaxy Node Pool CLI\n")
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Commit:     %s\n", Commit)
			fmt.Printf("Build Time: %s\n", BuildTime)
			fmt.Printf("AI-ID:      CP-GAL-NODEPOOL-001\n")
		},
	}
}
