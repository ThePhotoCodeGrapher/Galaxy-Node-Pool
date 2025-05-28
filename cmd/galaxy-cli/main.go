package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	mainnetURL string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "galaxy-cli",
	Short: "Galaxy CLI for managing pools and nodes",
	Long:  `A command-line interface for the Galaxy Network to manage pools, nodes, and their interactions with the main net.`,
}

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Manage Galaxy pools",
	Long:  `Create, start, stop, and manage Galaxy Node Pools.`,
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage Galaxy nodes",
	Long:  `Create, start, stop, and manage Galaxy Nodes.`,
}

var poolStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Galaxy pool",
	Long:  `Start a Galaxy Node Pool with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Galaxy Pool...")
		fmt.Printf("Using config: %s\n", configFile)
		if mainnetURL != "" {
			fmt.Printf("Connecting to main net: %s\n", mainnetURL)
		}
		// Here would be the actual code to start the pool
		fmt.Println("Pool started successfully!")
	},
}

var poolRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a pool with the main net",
	Long:  `Register a Galaxy Node Pool with the main net registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mainnetURL == "" {
			log.Fatal("Main net URL is required for registration")
		}
		fmt.Printf("Registering pool with main net: %s\n", mainnetURL)
		// Here would be the actual code to register with the main net
		fmt.Println("Pool registered successfully!")
	},
}

var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Galaxy node",
	Long:  `Start a Galaxy Node with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Galaxy Node...")
		fmt.Printf("Using config: %s\n", configFile)
		// Here would be the actual code to start the node
		fmt.Println("Node started successfully!")
	},
}

var nodeRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a node with a pool",
	Long:  `Register a Galaxy Node with a specified pool.`,
	Run: func(cmd *cobra.Command, args []string) {
		poolURL, _ := cmd.Flags().GetString("pool")
		if poolURL == "" {
			log.Fatal("Pool URL is required for registration")
		}
		fmt.Printf("Registering node with pool: %s\n", poolURL)
		// Here would be the actual code to register with the pool
		fmt.Println("Node registered successfully!")
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/example.yaml", "Config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Pool command flags
	poolRegisterCmd.Flags().StringVarP(&mainnetURL, "mainnet", "m", "", "Main net registry URL (required)")
	poolRegisterCmd.MarkFlagRequired("mainnet")

	// Node command flags
	nodeRegisterCmd.Flags().String("pool", "", "Pool URL to register with (required)")
	nodeRegisterCmd.MarkFlagRequired("pool")

	// Add commands to root
	rootCmd.AddCommand(poolCmd)
	rootCmd.AddCommand(nodeCmd)

	// Add subcommands to pool
	poolCmd.AddCommand(poolStartCmd)
	poolCmd.AddCommand(poolRegisterCmd)

	// Add subcommands to node
	nodeCmd.AddCommand(nodeStartCmd)
	nodeCmd.AddCommand(nodeRegisterCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
