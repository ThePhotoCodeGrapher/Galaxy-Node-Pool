// Galaxy Node Pool - Testnet Commands
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	
	"galaxy-node-pool/internal/testnet"
)

// testnetCmd creates a new testnet command
func testnetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Manage testnet pool",
	}

	// Add subcommands
	cmd.AddCommand(testnetInitCmd())
	cmd.AddCommand(testnetStartCmd())
	cmd.AddCommand(testnetStopCmd())
	cmd.AddCommand(testnetStatusCmd())
	cmd.AddCommand(testnetConfigCmd())
	cmd.AddCommand(testnetSslCmd())

	return cmd
}

// testnetInitCmd creates a command to initialize a testnet environment
func testnetInitCmd() *cobra.Command {
	var configPath string
	var orgID string
	var generateNginx bool
	var generateSSL bool

	cmd := &cobra.Command{
		Use:   "init [pool-name]",
		Short: "Initialize testnet environment",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			poolName := "test"
			if len(args) > 0 {
				poolName = args[0]
			}

			if orgID == "" {
				orgID = "testorg"
			}

			// Create testnet manager
			manager := testnet.NewManager("")

			// Initialize testnet environment
			listenAddr := "0.0.0.0:8080"
			grpcAddr := "0.0.0.0:9090"
			webAddr := "0.0.0.0:3000"

			if err := manager.Initialize(poolName, orgID, listenAddr, grpcAddr, webAddr); err != nil {
				fmt.Printf("Error initializing testnet environment: %v\n", err)
				return
			}

			// Generate Nginx configuration if requested
			if generateNginx {
				if err := manager.GenerateNginxConfig(poolName, orgID); err != nil {
					fmt.Printf("Error generating Nginx configuration: %v\n", err)
				}
			}

			// Generate SSL certificate if requested
			if generateSSL {
				if err := manager.GenerateSSLCertificate(poolName, orgID, false); err != nil {
					fmt.Printf("Error generating SSL certificate: %v\n", err)
				}
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to config directory")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID for the testnet domain (default: testorg)")
	cmd.Flags().BoolVar(&generateNginx, "nginx", false, "Generate Nginx configuration")
	cmd.Flags().BoolVar(&generateSSL, "ssl", false, "Generate self-signed SSL certificate")

	return cmd
}

// generateNginxConfig generates a Nginx configuration for the testnet using a template
func generateNginxConfig(nginxPath, poolName, orgID, apiPort, wsPort, webPort string) error {
	// Get the template path
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/root"
	}
	templateFile := filepath.Join(homeDir, "galaxy-node-pool", "templates", "nginx-testnet.conf.template")
	
	// Check if template exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", templateFile)
	}
	
	// Read template file
	templateData, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}
	
	// Get certificate paths
	certDir := filepath.Join(homeDir, ".galaxy", "certs")
	domain := fmt.Sprintf("*.%s.pool.galaxy.net.%s.asia.hybridconnect.cloud", poolName, orgID)
	certPath := filepath.Join(certDir, domain+".crt")
	keyPath := filepath.Join(certDir, domain+".key")
	
	// Check if certificates exist, otherwise use default snakeoil certs
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		certPath = "/etc/ssl/certs/ssl-cert-snakeoil.pem"
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		keyPath = "/etc/ssl/private/ssl-cert-snakeoil.key"
	}
	
	// Replace variables in template
	nginxConfig := string(templateData)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${pool_env}", poolName)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${org_id}", orgID)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${api_port}", apiPort)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ws_port}", wsPort)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${web_port}", webPort)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ssl_cert}", certPath)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ssl_key}", keyPath)
	
	// Write the configuration file
	return os.WriteFile(nginxPath, []byte(nginxConfig), 0644)
}

// testnetStartCmd creates a command to start a testnet pool
func testnetStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [pool-name]",
		Short: "Start a testnet pool",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			poolName := "test"
			if len(args) > 0 {
				poolName = args[0]
			}

			fmt.Printf("Starting testnet pool: %s\n", poolName)

			// Create testnet manager
			manager := testnet.NewManager("")

			// Start the testnet pool
			return manager.Start(poolName)
		},
	}

	return cmd
}

// testnetStopCmd creates a command to stop a testnet pool
func testnetStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [pool-name]",
		Short: "Stop a testnet pool",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			poolName := "test"
			if len(args) > 0 {
				poolName = args[0]
			}

			fmt.Printf("Stopping testnet pool: %s\n", poolName)

			// Create testnet manager
			manager := testnet.NewManager("")

			// Stop the testnet pool
			return manager.Stop(poolName)
		},
	}

	return cmd
}

// testnetStatusCmd creates a command to check the status of a testnet pool
func testnetStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [pool-name]",
		Short: "Check status of a testnet pool",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			poolName := "test"
			if len(args) > 0 {
				poolName = args[0]
			}

			fmt.Printf("Checking status of testnet pool: %s\n", poolName)

			// Create testnet manager
			manager := testnet.NewManager("")

			// Check the testnet pool status
			return manager.Status(poolName)
		},
	}

	return cmd
}

// testnetConfigCmd creates a command to manage testnet configuration
func testnetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage testnet configuration",
	}
	
	cmd.AddCommand(&cobra.Command{
		Use:   "show [pool-name]",
		Short: "Show testnet configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			poolName := args[0]
			fmt.Printf("Showing configuration for testnet pool: %s\n", poolName)
			// TODO: Implement configuration display
		},
	})
	
	cmd.AddCommand(&cobra.Command{
		Use:   "edit [pool-name]",
		Short: "Edit testnet configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			poolName := args[0]
			fmt.Printf("Editing configuration for testnet pool: %s\n", poolName)
			// TODO: Implement configuration editing
		},
	})
	
	return cmd
}

// Note: createTestnetConfig is defined in setup.go
