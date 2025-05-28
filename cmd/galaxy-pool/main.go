// Galaxy Node Pool
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/spf13/cobra"
)

var (
	// These will be set during build
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "none"
	GitBranch = ""
)

// PluginInterface defines what a plugin must implement
type PluginInterface interface {
	// GetCommands returns the commands provided by this plugin
	GetCommands() []*cobra.Command
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "galaxy-pool",
		Short: "Galaxy Node Pool - Manage your Galaxy Network nodes",
		Long: `Galaxy Node Pool is a high-performance node management system for the Galaxy Network.

This tool allows you to manage nodes, pools, and participate in the Galaxy Network.`,
		Version: fmt.Sprintf("%s\nBuild Date: %s\nCommit: %s\nBranch: %s", 
			Version, BuildDate, GitCommit, GitBranch),
	}

	// Add core commands
	rootCmd.AddCommand(versionCmd())
	
	// Add node, pool, and setup commands if they exist
	if nodeCmd := nodeCmd(); nodeCmd != nil {
		rootCmd.AddCommand(nodeCmd)
	}
	if poolCmd := poolCmd(); poolCmd != nil {
		rootCmd.AddCommand(poolCmd)
	}
	if setupCmd := setupCmd(); setupCmd != nil {
		rootCmd.AddCommand(setupCmd)
	}
	if buildCmd := buildCmd(); buildCmd != nil {
		rootCmd.AddCommand(buildCmd)
	}
	
	// Add testnet and domain commands
	rootCmd.AddCommand(testnetCmd())
	rootCmd.AddCommand(domainCmd())

	// Load plugins (enterprise features can be added here)
	loadPlugins(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadPlugins loads all plugins from the plugins directory
func loadPlugins(rootCmd *cobra.Command) {
	// Get plugin directory
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Warning: Could not determine executable path: %v\n", err)
		return
	}
	
	pluginDir := filepath.Join(filepath.Dir(execPath), "plugins")
	
	// Check if plugin directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// Plugin directory doesn't exist, that's okay
		return
	}
	
	// Read plugin directory
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		fmt.Printf("Warning: Could not read plugin directory: %v\n", err)
		return
	}
	
	// Load each plugin
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			loadPlugin(rootCmd, pluginPath)
		}
	}
}

// loadPlugin loads a single plugin and adds its commands to the root command
func loadPlugin(rootCmd *cobra.Command, path string) {
	// Open plugin
	p, err := plugin.Open(path)
	if err != nil {
		fmt.Printf("Warning: Could not open plugin %s: %v\n", path, err)
		return
	}
	
	// Look up the plugin's GetCommands symbol
	sym, err := p.Lookup("GetCommands")
	if err != nil {
		fmt.Printf("Warning: Plugin %s does not export GetCommands: %v\n", path, err)
		return
	}
	
	// Assert that the symbol is a function that returns []*cobra.Command
	getCommands, ok := sym.(func() []*cobra.Command)
	if !ok {
		fmt.Printf("Warning: Plugin %s's GetCommands has wrong type\n", path)
		return
	}
	
	// Add the plugin's commands to the root command
	for _, cmd := range getCommands() {
		rootCmd.AddCommand(cmd)
	}
}

func versionCmd() *cobra.Command {
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
