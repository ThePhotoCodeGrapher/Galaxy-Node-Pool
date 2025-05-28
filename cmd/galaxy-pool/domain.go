// Galaxy Node Pool - Domain Commands
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

func domainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage Galaxy domains",
		Long:  `Manage Galaxy domain configurations for pools.`,
	}

	// Add subcommands
	cmd.AddCommand(domainListCmd())
	cmd.AddCommand(domainAttachCmd())
	cmd.AddCommand(domainDetachCmd())
	cmd.AddCommand(domainInfoCmd())
	cmd.AddCommand(domainVerifyCmd())
	cmd.AddCommand(domainSslCmd())

	return cmd
}

func domainListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all attached domains",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing all attached domains...")
			// Implementation will go here
		},
	}
}

func domainAttachCmd() *cobra.Command {
	var poolName string
	var orgId string
	var environment string

	cmd := &cobra.Command{
		Use:   "attach [domain]",
		Short: "Attach a domain to a pool",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			
			// Validate domain format for wildcard pattern
			if !strings.Contains(domain, "pool.galaxy.net") {
				fmt.Printf("Error: Domain must follow the pattern: *.pool.galaxy.net.*.asia.hybridconnect.cloud\n")
				return
			}
			
			fmt.Printf("Attaching domain: %s\n", domain)
			fmt.Printf("Pool: %s\n", poolName)
			fmt.Printf("Organization: %s\n", orgId)
			fmt.Printf("Environment: %s\n", environment)
			
			// Implementation will go here
		},
	}
	
	cmd.Flags().StringVarP(&poolName, "pool", "p", "", "Pool name to attach the domain to")
	cmd.Flags().StringVarP(&orgId, "org", "o", "", "Organization ID")
	cmd.Flags().StringVarP(&environment, "env", "e", "dev", "Environment (dev, test, main)")
	cmd.MarkFlagRequired("pool")
	cmd.MarkFlagRequired("org")
	
	return cmd
}

func domainDetachCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detach [domain]",
		Short: "Detach a domain from a pool",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			fmt.Printf("Detaching domain: %s\n", domain)
			// Implementation will go here
		},
	}
}

func domainInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [domain]",
		Short: "Get information about a domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			fmt.Printf("Domain information for: %s\n", domain)
			// Implementation will go here
		},
	}
}

func domainVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [domain]",
		Short: "Verify domain ownership",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			fmt.Printf("Verifying domain ownership: %s\n", domain)
			// Implementation will go here
		},
	}
}

func domainSslCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Manage SSL certificates for domains",
	}
	
	var dnsProvider string
	var email string
	var testMode bool
	var configPath string
	var nginxConfig string
	
	generateCmd := &cobra.Command{
		Use:   "generate [domain]",
		Short: "Generate SSL certificate for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			fmt.Printf("Generating SSL certificate for domain: %s\n", domain)
			
			// Check if domain is a wildcard domain
			isWildcard := strings.Contains(domain, "*")
			
			// Create certificate manager
			certDir := filepath.Join(os.Getenv("HOME"), ".galaxy", "certs")
			manager := cert.NewManager(certDir, nginxConfig)
			
			// For testnet, we can use self-signed certificates if in test mode
			if testMode {
				fmt.Println("Using test mode with self-signed certificates")
				return manager.GenerateSelfSigned(domain)
			}
			
			// Load certificate configuration
			config, err := cert.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load certificate configuration: %v", err)
			}
			
			// Use provided email and DNS provider if not specified in config
			if email != "" {
				config.Email = email
			}
			if dnsProvider != "" {
				config.DNSProvider = dnsProvider
			}
			
			// For wildcard domains, we need to use DNS challenge
			if isWildcard {
				if config.DNSProvider == "" {
					return fmt.Errorf("DNS provider is required for wildcard certificates")
				}
				return manager.GenerateWithLetsEncrypt(domain, config.Email, true, config.DNSProvider)
			}
			
			// For regular domains, we can use HTTP challenge
			return manager.GenerateWithLetsEncrypt(domain, config.Email, false, "")
		},
	}
	
	// Add flags
	generateCmd.Flags().StringVar(&dnsProvider, "dns-provider", "", "DNS provider for wildcard certificates (e.g., cloudflare, route53)")
	generateCmd.Flags().StringVar(&email, "email", "", "Email address for Let's Encrypt notifications")
	generateCmd.Flags().BoolVar(&testMode, "test", false, "Use test mode with self-signed certificates for testnet")
	generateCmd.Flags().StringVar(&configPath, "config", "", "Path to certificate configuration file")
	generateCmd.Flags().StringVar(&nginxConfig, "nginx-config", "/etc/nginx/sites-available/galaxy-pool.conf", "Path to Nginx configuration file")
	
	cmd.AddCommand(generateCmd)
	
	cmd.AddCommand(&cobra.Command{
		Use:   "renew [domain]",
		Short: "Renew SSL certificate for a domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			fmt.Printf("Renewing SSL certificate for domain: %s\n", domain)
			// Implementation will go here
		},
	})
	
	cmd.AddCommand(&cobra.Command{
		Use:   "status [domain]",
		Short: "Check SSL certificate status for a domain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			fmt.Printf("SSL certificate status for domain: %s\n", domain)
			// Implementation will go here
		},
	})
	
	return cmd
}
