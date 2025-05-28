// Galaxy Node Pool - Testnet SSL Command
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	
	"galaxy-node-pool/internal/cert"
)

// testnetSslCmd creates a command for managing SSL certificates for testnet
func testnetSslCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Manage SSL certificates for testnet",
	}
	
	// Add subcommands
	cmd.AddCommand(testnetSslGenerateCmd())
	
	return cmd
}

// testnetSslGenerateCmd creates a command for generating SSL certificates for testnet
func testnetSslGenerateCmd() *cobra.Command {
	var poolName string
	var orgID string
	var useProduction bool
	var nginxConfigPath string
	
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate SSL certificate for testnet domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Construct the domain pattern for testnet
			if poolName == "" {
				poolName = "test"
			}
			
			if orgID == "" {
				orgID = "testorg"
			}
			
			// Construct the domain for the testnet pool
			domain := fmt.Sprintf("*.%s.pool.galaxy.net.%s.asia.hybridconnect.cloud", poolName, orgID)
			
			fmt.Printf("Generating SSL certificate for testnet domain: %s\n", domain)
			
			// Create certificate manager
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				homeDir = "/root"
			}
			
			certDir := filepath.Join(homeDir, ".galaxy", "certs")
			
			// Use the specified Nginx config or the testnet one
			if nginxConfigPath == "" {
				nginxConfigPath = filepath.Join(homeDir, "galaxy-node-pool", "testnet-nginx.conf")
			}
			
			manager := cert.NewManager(certDir, nginxConfigPath)
			
			// For testnet, we use self-signed certificates by default
			// In production mode, we use Let's Encrypt
			if useProduction {
				// Load certificate configuration
				configPath := filepath.Join(homeDir, ".galaxy", "cert-config.yaml")
				config, err := cert.LoadConfig(configPath)
				if err != nil {
					return fmt.Errorf("failed to load certificate configuration: %v", err)
				}
				
				// Check if DNS provider is configured
				if config.DNSProvider == "" {
					return fmt.Errorf("DNS provider is required for wildcard certificates. Configure it in %s", configPath)
				}
				
				// Generate certificate using Let's Encrypt
				return manager.GenerateWithLetsEncrypt(domain, config.Email, true, config.DNSProvider)
			}
			
			// Generate self-signed certificate for testnet
			return manager.GenerateSelfSigned(domain)
		},
	}
	
	// Add flags
	cmd.Flags().StringVar(&poolName, "pool-name", "", "Pool name for the testnet domain (default: test)")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID for the testnet domain (default: testorg)")
	cmd.Flags().BoolVar(&useProduction, "production", false, "Use production-grade certificates from Let's Encrypt")
	cmd.Flags().StringVar(&nginxConfigPath, "nginx-config", "", "Path to Nginx configuration file")
	
	return cmd
}
