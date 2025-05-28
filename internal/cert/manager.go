// Galaxy Node Pool - Certificate Manager
// AI-ID: CP-GAL-NODEPOOL-001
package cert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles certificate operations for Galaxy Node Pool
type Manager struct {
	CertDir     string
	NginxConfig string
}

// NewManager creates a new certificate manager
func NewManager(certDir, nginxConfig string) *Manager {
	if certDir == "" {
		certDir = filepath.Join(os.Getenv("HOME"), ".galaxy", "certs")
	}
	
	return &Manager{
		CertDir:     certDir,
		NginxConfig: nginxConfig,
	}
}

// GenerateSelfSigned generates a self-signed certificate for testnet environments
func (m *Manager) GenerateSelfSigned(domain string) error {
	// Create cert directory if it doesn't exist
	if err := os.MkdirAll(m.CertDir, 0755); err != nil {
		return fmt.Errorf("failed to create cert directory: %v", err)
	}
	
	// Generate a private key
	keyPath := filepath.Join(m.CertDir, domain+".key")
	cmd := exec.Command("openssl", "genrsa", "-out", keyPath, "2048")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}
	
	// Generate a CSR configuration
	csrConfigPath := filepath.Join(m.CertDir, domain+".cnf")
	csrConfig := fmt.Sprintf(`[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext

[dn]
CN = %s

[req_ext]
subjectAltName = @alt_names

[alt_names]
`, domain)
	
	// Handle wildcard domains
	if strings.HasPrefix(domain, "*.") {
		baseDomain := strings.TrimPrefix(domain, "*.")
		csrConfig += fmt.Sprintf("DNS.1 = %s\nDNS.2 = *.%s\n", baseDomain, baseDomain)
	} else {
		csrConfig += fmt.Sprintf("DNS.1 = %s\n", domain)
	}
	
	if err := os.WriteFile(csrConfigPath, []byte(csrConfig), 0644); err != nil {
		return fmt.Errorf("failed to write CSR config: %v", err)
	}
	
	// Generate a CSR
	csrPath := filepath.Join(m.CertDir, domain+".csr")
	cmd = exec.Command("openssl", "req", "-new", "-key", keyPath, "-out", csrPath, "-config", csrConfigPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate CSR: %v", err)
	}
	
	// Generate a self-signed certificate
	certPath := filepath.Join(m.CertDir, domain+".crt")
	cmd = exec.Command("openssl", "x509", "-req", "-days", "365", "-in", csrPath, "-signkey", keyPath, "-out", certPath, "-extensions", "req_ext", "-extfile", csrConfigPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate self-signed certificate: %v", err)
	}
	
	fmt.Printf("Self-signed certificate generated at: %s\n", certPath)
	fmt.Printf("Private key generated at: %s\n", keyPath)
	
	// Update Nginx configuration if provided
	if m.NginxConfig != "" {
		if err := m.UpdateNginxConfig(domain, certPath, keyPath); err != nil {
			return fmt.Errorf("failed to update Nginx config: %v", err)
		}
	}
	
	return nil
}

// GenerateWithLetsEncrypt generates a certificate using Let's Encrypt
func (m *Manager) GenerateWithLetsEncrypt(domain, email string, wildcard bool, dnsProvider string) error {
	// Check if certbot is installed
	if _, err := exec.LookPath("certbot"); err != nil {
		return fmt.Errorf("certbot is not installed: %v", err)
	}
	
	var cmd *exec.Cmd
	
	if wildcard {
		// For wildcard certificates, we need to use DNS challenge
		if dnsProvider == "" {
			return fmt.Errorf("DNS provider is required for wildcard certificates")
		}
		
		// Install the DNS plugin if needed
		pluginPackage := fmt.Sprintf("python3-certbot-dns-%s", dnsProvider)
		installCmd := exec.Command("apt", "list", "--installed", pluginPackage)
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Warning: DNS plugin for %s might not be installed\n", dnsProvider)
		}
		
		// Generate certificate using DNS challenge
		cmd = exec.Command("certbot", "certonly",
			"--dns-"+dnsProvider,
			"--agree-tos",
			"--email", email,
			"-d", domain)
	} else {
		// For regular certificates, we can use HTTP challenge
		cmd = exec.Command("certbot", "certonly",
			"--webroot",
			"--webroot-path", "/var/www/html",
			"--agree-tos",
			"--email", email,
			"-d", domain)
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate certificate: %v", err)
	}
	
	// Update Nginx configuration if provided
	if m.NginxConfig != "" {
		certPath := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain)
		keyPath := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain)
		if err := m.UpdateNginxConfig(domain, certPath, keyPath); err != nil {
			return fmt.Errorf("failed to update Nginx config: %v", err)
		}
	}
	
	return nil
}

// UpdateNginxConfig updates the Nginx configuration with the certificate paths
func (m *Manager) UpdateNginxConfig(domain, certPath, keyPath string) error {
	// Read the Nginx configuration
	config, err := os.ReadFile(m.NginxConfig)
	if err != nil {
		return fmt.Errorf("failed to read Nginx config: %v", err)
	}
	
	// Replace the certificate paths
	configStr := string(config)
	configStr = strings.Replace(configStr, "ssl_certificate /etc/ssl/certs/ssl-cert-snakeoil.pem;", fmt.Sprintf("ssl_certificate %s;", certPath), -1)
	configStr = strings.Replace(configStr, "ssl_certificate_key /etc/ssl/private/ssl-cert-snakeoil.key;", fmt.Sprintf("ssl_certificate_key %s;", keyPath), -1)
	
	// Write the updated configuration
	if err := os.WriteFile(m.NginxConfig, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write Nginx config: %v", err)
	}
	
	fmt.Printf("Nginx configuration updated with certificate paths\n")
	fmt.Printf("To apply the changes, run: sudo nginx -t && sudo systemctl reload nginx\n")
	
	return nil
}
