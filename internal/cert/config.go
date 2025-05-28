// Galaxy Node Pool - Certificate Configuration
// AI-ID: CP-GAL-NODEPOOL-001
package cert

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"gopkg.in/yaml.v3"
)

// Config represents certificate configuration
type Config struct {
	Email       string            `yaml:"email"`
	DNSProvider string            `yaml:"dns_provider"`
	Credentials map[string]string `yaml:"credentials"`
}

// LoadConfig loads certificate configuration from file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".galaxy", "cert-config.yaml")
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &Config{
			Email:       "",
			DNSProvider: "",
			Credentials: make(map[string]string),
		}, nil
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	// Parse config file
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	
	return &config, nil
}

// SaveConfig saves certificate configuration to file
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".galaxy", "cert-config.yaml")
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	// Write config file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	return nil
}

// CreateCredentialsFile creates a credentials file for the DNS provider
func CreateCredentialsFile(config *Config, provider string) (string, error) {
	if config.Credentials == nil || len(config.Credentials) == 0 {
		return "", fmt.Errorf("no credentials found for DNS provider")
	}
	
	// Create credentials directory
	credsDir := filepath.Join(os.Getenv("HOME"), ".galaxy", "certs", "credentials")
	if err := os.MkdirAll(credsDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create credentials directory: %v", err)
	}
	
	// Create credentials file
	credsPath := filepath.Join(credsDir, provider+".ini")
	var content strings.Builder
	
	// Different providers have different credential formats
	switch provider {
	case "cloudflare":
		content.WriteString("# Cloudflare API credentials used by Certbot\n")
		content.WriteString("dns_cloudflare_email = " + config.Credentials["email"] + "\n")
		content.WriteString("dns_cloudflare_api_key = " + config.Credentials["api_key"] + "\n")
	case "route53":
		content.WriteString("[default]\n")
		content.WriteString("aws_access_key_id = " + config.Credentials["access_key"] + "\n")
		content.WriteString("aws_secret_access_key = " + config.Credentials["secret_key"] + "\n")
	default:
		// Generic format
		for key, value := range config.Credentials {
			content.WriteString(key + " = " + value + "\n")
		}
	}
	
	// Write credentials file with restricted permissions
	if err := os.WriteFile(credsPath, []byte(content.String()), 0600); err != nil {
		return "", fmt.Errorf("failed to write credentials file: %v", err)
	}
	
	return credsPath, nil
}
