// Galaxy Node Pool - Build Commands
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// buildCmd creates a new build command
func buildCmd() *cobra.Command {
	var outputDir string
	var binaryName string
	var ldflags string
	var targetOS string
	var targetArch string
	var installBinary bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build Galaxy Node Pool binaries",
		Long:  `Build Galaxy Node Pool binaries for various platforms.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Building Galaxy Node Pool...")

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}

			// Get version information
			version, err := getVersion()
			if err != nil {
				return fmt.Errorf("failed to get version: %v", err)
			}

			// Get git information
			gitCommit, gitBranch, err := getGitInfo()
			if err != nil {
				fmt.Printf("Warning: Failed to get git info: %v\n", err)
				gitCommit = "unknown"
				gitBranch = "unknown"
			}

			// Get build date
			buildDate := time.Now().UTC().Format(time.RFC3339)

			// Set ldflags
			if ldflags == "" {
				ldflags = fmt.Sprintf("-X main.Version=%s -X main.BuildDate=%s -X main.GitCommit=%s -X main.GitBranch=%s",
					version, buildDate, gitCommit, gitBranch)
			}

			// Build the binary
			outputPath := filepath.Join(outputDir, binaryName)
			if targetOS != runtime.GOOS || targetArch != runtime.GOARCH {
				outputPath = fmt.Sprintf("%s-%s-%s", outputPath, targetOS, targetArch)
				if targetOS == "windows" {
					outputPath += ".exe"
				}
			}

			fmt.Printf("Building binary for %s/%s...\n", targetOS, targetArch)
			buildCmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputPath, "./cmd/galaxy-pool")
			buildCmd.Env = append(os.Environ(),
				fmt.Sprintf("GOOS=%s", targetOS),
				fmt.Sprintf("GOARCH=%s", targetArch),
				"CGO_ENABLED=0",
			)
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr

			if err := buildCmd.Run(); err != nil {
				return fmt.Errorf("build failed: %v", err)
			}

			// Make the binary executable
			if targetOS != "windows" {
				if err := os.Chmod(outputPath, 0755); err != nil {
					return fmt.Errorf("failed to make binary executable: %v", err)
				}
			}

			// Install the binary if requested
			if installBinary {
				gopath := os.Getenv("GOPATH")
				if gopath == "" {
					gopath = filepath.Join(os.Getenv("HOME"), "go")
				}
				installPath := filepath.Join(gopath, "bin", binaryName)
				if err := copyFile(outputPath, installPath); err != nil {
					return fmt.Errorf("failed to install binary: %v", err)
				}
				fmt.Printf("Binary installed to: %s\n", installPath)
			}

			fmt.Printf("Build complete!\n")
			fmt.Printf("Binary location: %s\n", outputPath)
			fmt.Printf("Version: %s (%s on %s)\n", version, gitCommit, gitBranch)

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&outputDir, "output-dir", "bin", "Output directory for binaries")
	cmd.Flags().StringVar(&binaryName, "name", "galaxy-pool", "Binary name")
	cmd.Flags().StringVar(&ldflags, "ldflags", "", "Go build ldflags")
	cmd.Flags().StringVar(&targetOS, "os", runtime.GOOS, "Target OS (e.g., linux, darwin, windows)")
	cmd.Flags().StringVar(&targetArch, "arch", runtime.GOARCH, "Target architecture (e.g., amd64, arm64)")
	cmd.Flags().BoolVar(&installBinary, "install", false, "Install binary to GOPATH/bin")

	// Add subcommands
	cmd.AddCommand(buildAllCmd())

	return cmd
}

// buildAllCmd creates a command to build for all platforms
func buildAllCmd() *cobra.Command {
	var outputDir string
	var binaryName string

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Build for all platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Building Galaxy Node Pool for all platforms...")

			// Define platforms
			platforms := []struct {
				os   string
				arch string
			}{
				{"linux", "amd64"},
				{"linux", "arm64"},
				{"darwin", "amd64"},
				{"darwin", "arm64"},
				{"windows", "amd64"},
			}

			// Build for each platform
			for _, platform := range platforms {
				// Get version information
				version, err := getVersion()
				if err != nil {
					return fmt.Errorf("failed to get version: %v", err)
				}

				// Get git information
				gitCommit, gitBranch, err := getGitInfo()
				if err != nil {
					fmt.Printf("Warning: Failed to get git info: %v\n", err)
					gitCommit = "unknown"
					gitBranch = "unknown"
				}

				// Get build date
				buildDate := time.Now().UTC().Format(time.RFC3339)

				// Set ldflags
				ldflags := fmt.Sprintf("-X main.Version=%s -X main.BuildDate=%s -X main.GitCommit=%s -X main.GitBranch=%s",
					version, buildDate, gitCommit, gitBranch)

				// Determine output path
				outputPath := filepath.Join(outputDir, binaryName)
				outputPath = fmt.Sprintf("%s-%s-%s", outputPath, platform.os, platform.arch)
				if platform.os == "windows" {
					outputPath += ".exe"
				}

				fmt.Printf("Building for %s/%s...\n", platform.os, platform.arch)
				buildCmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputPath, "./cmd/galaxy-pool")
				buildCmd.Env = append(os.Environ(),
					fmt.Sprintf("GOOS=%s", platform.os),
					fmt.Sprintf("GOARCH=%s", platform.arch),
					"CGO_ENABLED=0",
				)
				buildCmd.Stdout = os.Stdout
				buildCmd.Stderr = os.Stderr

				if err := buildCmd.Run(); err != nil {
					return fmt.Errorf("failed to build for %s/%s: %v", platform.os, platform.arch, err)
				}
			}

			fmt.Println("All builds completed successfully!")
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&outputDir, "output-dir", "bin", "Output directory for binaries")
	cmd.Flags().StringVar(&binaryName, "name", "galaxy-pool", "Binary name")

	return cmd
}

// getVersion gets the version from the VERSION file or defaults to "dev"
func getVersion() (string, error) {
	versionFile := "VERSION"
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		return "dev", nil
	}

	versionBytes, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(versionBytes))
	if version == "" {
		return "dev", nil
	}

	return version, nil
}

// getGitInfo gets the git commit and branch information
func getGitInfo() (string, string, error) {
	// Get git commit
	commitCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	commitOutput, err := commitCmd.Output()
	if err != nil {
		return "", "", err
	}
	commit := strings.TrimSpace(string(commitOutput))

	// Get git branch
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		return commit, "", err
	}
	branch := strings.TrimSpace(string(branchOutput))

	return commit, branch, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Read the source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write to the destination file
	return os.WriteFile(dst, data, 0755)
}
