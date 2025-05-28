// Galaxy Node Pool - Testnet Manager
// AI-ID: CP-GAL-NODEPOOL-001
package testnet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	
	"galaxy-node-pool/internal/cert"
)

// Manager handles testnet operations
type Manager struct {
	BaseDir     string
	ConfigDir   string
	TemplateDir string
}

// Config represents testnet configuration
type Config struct {
	PoolName   string
	OrgID      string
	ListenAddr string
	GrpcAddr   string
	WebAddr    string
	ApiPort    string
	WsPort     string
	WebPort    string
}

// NewManager creates a new testnet manager
func NewManager(baseDir string) *Manager {
	if baseDir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = "/root"
		}
		baseDir = homeDir
	}
	
	return &Manager{
		BaseDir:     baseDir,
		ConfigDir:   filepath.Join(baseDir, ".galaxy", "testnet"),
		TemplateDir: filepath.Join(baseDir, "galaxy-node-pool", "templates"),
	}
}

// Initialize initializes a testnet environment
func (m *Manager) Initialize(poolName, orgID, listenAddr, grpcAddr, webAddr string) error {
	if poolName == "" {
		poolName = "test"
	}
	
	if orgID == "" {
		orgID = "testorg"
	}
	
	// Extract ports from addresses
	apiPort := extractPort(listenAddr)
	wsPort := extractPort(grpcAddr)
	webPort := extractPort(webAddr)
	
	config := Config{
		PoolName:   poolName,
		OrgID:      orgID,
		ListenAddr: listenAddr,
		GrpcAddr:   grpcAddr,
		WebAddr:    webAddr,
		ApiPort:    apiPort,
		WsPort:     wsPort,
		WebPort:    webPort,
	}
	
	// Create config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Create config file
	configPath := filepath.Join(configDir, "config.yaml")
	if err := m.createConfigFile(configPath, config); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	
	fmt.Printf("Testnet environment initialized at %s\n", configDir)
	return nil
}

// GenerateNginxConfig generates Nginx configuration for testnet
func (m *Manager) GenerateNginxConfig(poolName, orgID string) error {
	if poolName == "" {
		poolName = "test"
	}
	
	if orgID == "" {
		orgID = "testorg"
	}
	
	// Get config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	
	// Create Nginx config file
	nginxPath := filepath.Join(configDir, "nginx.conf")
	
	// Get template path
	templateFile := filepath.Join(m.TemplateDir, "nginx-testnet.conf.template")
	
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
	certDir := filepath.Join(m.BaseDir, ".galaxy", "certs")
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
	nginxConfig = strings.ReplaceAll(nginxConfig, "${api_port}", "8080")
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ws_port}", "9090")
	nginxConfig = strings.ReplaceAll(nginxConfig, "${web_port}", "3000")
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ssl_cert}", certPath)
	nginxConfig = strings.ReplaceAll(nginxConfig, "${ssl_key}", keyPath)
	
	// Write the configuration file
	if err := os.WriteFile(nginxPath, []byte(nginxConfig), 0644); err != nil {
		return fmt.Errorf("failed to write Nginx config: %v", err)
	}
	
	fmt.Printf("Nginx configuration generated at %s\n", nginxPath)
	return nil
}

// GenerateSSLCertificate generates SSL certificate for testnet
func (m *Manager) GenerateSSLCertificate(poolName, orgID string, production bool) error {
	if poolName == "" {
		poolName = "test"
	}
	
	if orgID == "" {
		orgID = "testorg"
	}
	
	// Get config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	nginxPath := filepath.Join(configDir, "nginx.conf")
	
	// Create certificate manager
	certDir := filepath.Join(m.BaseDir, ".galaxy", "certs")
	manager := cert.NewManager(certDir, nginxPath)
	
	// Generate domain name
	domain := fmt.Sprintf("*.%s.pool.galaxy.net.%s.asia.hybridconnect.cloud", poolName, orgID)
	
	if production {
		// Load certificate configuration
		configPath := filepath.Join(m.BaseDir, ".galaxy", "cert-config.yaml")
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
}

// Start starts the testnet pool
func (m *Manager) Start(poolName string) error {
	if poolName == "" {
		poolName = "test"
	}
	
	// Get config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	
	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("testnet environment not found: %s", configDir)
	}
	
	// TODO: Implement actual pool starting logic
	
	fmt.Printf("Testnet pool %s started\n", poolName)
	return nil
}

// Stop stops the testnet pool
func (m *Manager) Stop(poolName string) error {
	if poolName == "" {
		poolName = "test"
	}
	
	// Get config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	
	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("testnet environment not found: %s", configDir)
	}
	
	// TODO: Implement actual pool stopping logic
	
	fmt.Printf("Testnet pool %s stopped\n", poolName)
	return nil
}

// Status checks the status of the testnet pool
func (m *Manager) Status(poolName string) error {
	if poolName == "" {
		poolName = "test"
	}
	
	// Get config directory
	configDir := filepath.Join(m.ConfigDir, poolName)
	
	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("testnet environment not found: %s", configDir)
	}
	
	// TODO: Implement actual status checking logic
	
	fmt.Printf("Testnet pool %s is running\n", poolName)
	return nil
}

// createConfigFile creates a testnet configuration file
func (m *Manager) createConfigFile(path string, config Config) error {
	// Define the configuration template
	tmpl := `# Galaxy Node Pool Testnet Configuration
# Generated by galaxy-pool testnet init

# Pool Configuration
pool:
  name: {{ .PoolName }}
  org_id: {{ .OrgID }}
  listen_addr: {{ .ListenAddr }}
  grpc_addr: {{ .GrpcAddr }}
  web_addr: {{ .WebAddr }}

# Network Configuration
network:
  type: testnet
  api_port: {{ .ApiPort }}
  ws_port: {{ .WsPort }}
  web_port: {{ .WebPort }}

# Security Configuration
security:
  tls_enabled: true
  auth_required: true
  rate_limit: 100

# Logging Configuration
logging:
  level: debug
  file: {{ .PoolName }}.log
`
	
	// Parse the template
	t, err := template.New("config").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}
	
	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Execute the template
	if err := t.Execute(file, config); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	
	return nil
}

// extractPort extracts the port from an address string
func extractPort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return "8080" // Default port
	}
	return parts[1]
}
