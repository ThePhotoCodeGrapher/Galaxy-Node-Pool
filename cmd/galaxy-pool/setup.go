// Galaxy Node Pool - Setup Commands
// AI-ID: CP-GAL-NODEPOOL-001
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// setupCmd creates a new setup command
func setupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup Galaxy Node Pool environments",
		Long:  `Setup and configure Galaxy Node Pool environments.`,
	}

	// Add subcommands
	cmd.AddCommand(setupTestnetCmd())
	cmd.AddCommand(setupNginxCmd())
	cmd.AddCommand(setupFirewallCmd())

	return cmd
}

// setupTestnetCmd creates a command to set up a testnet environment
func setupTestnetCmd() *cobra.Command {
	var configDir string
	var listenAddr string
	var grpcAddr string
	var webAddr string

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Setup a testnet environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Setting up testnet environment...")

			// Create config directory if it doesn't exist
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %v", err)
			}

			// Create testnet configuration
			configPath := filepath.Join(configDir, "config.yaml")
			poolName := "testpool" // Default pool name for testnet
			if err := createTestnetConfig(configPath, poolName, listenAddr, grpcAddr, webAddr); err != nil {
				return fmt.Errorf("failed to create testnet config: %v", err)
			}

			fmt.Printf("Testnet configuration created at: %s\n", configPath)
			fmt.Println("To start the testnet pool, run:")
			fmt.Println("  galaxy-pool pool start testpool --testnet")

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&configDir, "config-dir", filepath.Join(os.Getenv("HOME"), ".galaxy", "testnet"), "Configuration directory")
	cmd.Flags().StringVar(&listenAddr, "listen-addr", "0.0.0.0:3000", "API listen address")
	cmd.Flags().StringVar(&grpcAddr, "grpc-addr", "0.0.0.0:3001", "gRPC listen address")
	cmd.Flags().StringVar(&webAddr, "web-addr", "0.0.0.0:8080", "Web interface listen address")

	return cmd
}

// setupNginxCmd creates a command to set up Nginx configuration
func setupNginxCmd() *cobra.Command {
	var outputPath string
	var serverName string
	var apiPort string
	var webPort string

	cmd := &cobra.Command{
		Use:   "nginx",
		Short: "Generate Nginx configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Generating Nginx configuration...")

			// Create Nginx configuration
			if err := createNginxConfig(outputPath, serverName, apiPort, webPort); err != nil {
				return fmt.Errorf("failed to create Nginx config: %v", err)
			}

			fmt.Printf("Nginx configuration created at: %s\n", outputPath)
			fmt.Println("To enable this configuration, run:")
			fmt.Println("  sudo ln -sf " + outputPath + " /etc/nginx/sites-enabled/")
			fmt.Println("  sudo nginx -t")
			fmt.Println("  sudo systemctl reload nginx")

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&outputPath, "output", "galaxy-pool-nginx.conf", "Output file path")
	cmd.Flags().StringVar(&serverName, "server-name", "pool.example.com", "Server name")
	cmd.Flags().StringVar(&apiPort, "api-port", "3000", "API port")
	cmd.Flags().StringVar(&webPort, "web-port", "8080", "Web interface port")

	return cmd
}

// setupFirewallCmd creates a command to set up firewall rules
func setupFirewallCmd() *cobra.Command {
	var apiPort string
	var grpcPort string
	var webPort string
	var generateOnly bool

	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "Generate firewall rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Generating firewall rules...")

			// Generate firewall rules
			rules := generateFirewallRules(apiPort, grpcPort, webPort)

			// Print the rules
			fmt.Println("Firewall Rules:")
			fmt.Println("---------------")
			for _, rule := range rules {
				fmt.Println(rule)
			}

			if !generateOnly {
				fmt.Println("\nTo apply these rules, you would need to run them with appropriate permissions.")
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&apiPort, "api-port", "3000", "API port")
	cmd.Flags().StringVar(&grpcPort, "grpc-port", "3001", "gRPC port")
	cmd.Flags().StringVar(&webPort, "web-port", "8080", "Web interface port")
	cmd.Flags().BoolVar(&generateOnly, "generate-only", true, "Only generate rules, don't apply them")

	return cmd
}

// createTestnetConfig creates a testnet configuration file
func createTestnetConfig(path, poolName, listenAddr, grpcAddr, webAddr string) error {
	config := `# Galaxy Node Pool Testnet Configuration
version: 1.0
environment: testnet

# Pool configuration
pool:
  name: {{.PoolName}}
  description: "Galaxy Node Pool Testnet"
  listen_address: "{{.ListenAddr}}"
  grpc_address: "{{.GrpcAddr}}"
  web_address: "{{.WebAddr}}"

# Node configuration
node:
  default_specialization: "developer"
  health_check_interval: 30
  timeout: 10

# Security configuration
security:
  require_authentication: true
  allow_anonymous_queries: true
`

	// Create template
	tmpl, err := template.New("config").Parse(config)
	if err != nil {
		return err
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template
	data := struct {
		PoolName   string
		ListenAddr string
		GrpcAddr   string
		WebAddr    string
	}{
		PoolName:   poolName,
		ListenAddr: listenAddr,
		GrpcAddr:   grpcAddr,
		WebAddr:    webAddr,
	}

	return tmpl.Execute(f, data)
}

// createNginxConfig creates an Nginx configuration file
func createNginxConfig(path, serverName, apiPort, webPort string) error {
	config := `# Galaxy Node Pool Nginx Configuration
server {
    listen 80;
    server_name {{.ServerName}};
    
    location / {
        return 301 https://$host$request_uri;
    }
    
    # Let's Encrypt verification
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }
}

server {
    listen 443 ssl http2;
    server_name {{.ServerName}};
    
    # SSL Configuration (placeholder - will be updated by certbot)
    ssl_certificate /etc/ssl/certs/ssl-cert-snakeoil.pem;
    ssl_certificate_key /etc/ssl/private/ssl-cert-snakeoil.key;
    
    # SSL Settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;
    
    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Logs
    access_log /var/log/nginx/{{.ServerName}}.access.log;
    error_log /var/log/nginx/{{.ServerName}}.error.log;
    
    # API Service
    location /api/ {
        proxy_pass http://localhost:{{.ApiPort}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Web Interface
    location / {
        proxy_pass http://localhost:{{.WebPort}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
`

	// Create template
	tmpl, err := template.New("nginx").Parse(config)
	if err != nil {
		return err
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template
	data := struct {
		ServerName string
		ApiPort    string
		WebPort    string
	}{
		ServerName: serverName,
		ApiPort:    apiPort,
		WebPort:    webPort,
	}

	return tmpl.Execute(f, data)
}

// generateFirewallRules generates firewall rules for the Galaxy Node Pool
func generateFirewallRules(apiPort, grpcPort, webPort string) []string {
	var rules []string

	// UFW rules
	rules = append(rules, "# UFW Rules")
	rules = append(rules, fmt.Sprintf("sudo ufw allow %s/tcp # API", apiPort))
	rules = append(rules, fmt.Sprintf("sudo ufw allow %s/tcp # gRPC", grpcPort))
	rules = append(rules, fmt.Sprintf("sudo ufw allow %s/tcp # Web", webPort))
	rules = append(rules, "sudo ufw allow 80/tcp # HTTP")
	rules = append(rules, "sudo ufw allow 443/tcp # HTTPS")

	// iptables rules
	rules = append(rules, "\n# iptables Rules")
	rules = append(rules, fmt.Sprintf("sudo iptables -A INPUT -p tcp --dport %s -j ACCEPT # API", apiPort))
	rules = append(rules, fmt.Sprintf("sudo iptables -A INPUT -p tcp --dport %s -j ACCEPT # gRPC", grpcPort))
	rules = append(rules, fmt.Sprintf("sudo iptables -A INPUT -p tcp --dport %s -j ACCEPT # Web", webPort))
	rules = append(rules, "sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT # HTTP")
	rules = append(rules, "sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT # HTTPS")

	return rules
}
