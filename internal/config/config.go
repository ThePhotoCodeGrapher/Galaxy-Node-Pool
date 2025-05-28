package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Address string `mapstructure:"address"`
		TLS     struct {
			Enabled  bool   `mapstructure:"enabled"`
			CertFile string `mapstructure:"cert_file"`
			KeyFile  string `mapstructure:"key_file"`
		} `mapstructure:"tls"`
		MaxConnections int `mapstructure:"max_connections"`
		Resources      struct {
			CPULimit    string `mapstructure:"cpu_limit"`
			MemoryLimit string `mapstructure:"memory_limit"`
		} `mapstructure:"resources"`
		Location struct {
			Region      string `mapstructure:"region"`
			Datacenter  string `mapstructure:"datacenter"`
			Coordinates string `mapstructure:"coordinates"`
		} `mapstructure:"location"`
	} `mapstructure:"server"`

	// Main net configuration
	MainNet struct {
		RegistryAddress string `mapstructure:"registry_address"`
		RegistrationFee string `mapstructure:"registration_fee"`
		APIKey          string `mapstructure:"api_key"`
	} `mapstructure:"mainnet"`

	// Domain configuration
	Domain struct {
		DomainName  string `mapstructure:"domain_name"`
		DNSProvider string `mapstructure:"dns_provider"`
		DNSAPIKey   string `mapstructure:"dns_api_key"`
		SSLCertPath string `mapstructure:"ssl_cert_path"`
		SSLKeyPath  string `mapstructure:"ssl_key_path"`
	} `mapstructure:"domain"`

	// Registry configuration
	Registry struct {
		AllowPublicRegistration bool     `mapstructure:"allow_public_registration"`
		AllowedOrgs             []string `mapstructure:"allowed_orgs"`
		HealthCheckInterval     string   `mapstructure:"health_check_interval"`
		MaxNodes                int      `mapstructure:"max_nodes"`
		AutoDeregisterAfter     int      `mapstructure:"auto_deregister_after"`
		Plugins                 []struct {
			Name    string                 `mapstructure:"name"`
			Enabled bool                   `mapstructure:"enabled"`
			Config  map[string]interface{} `mapstructure:"config"`
		} `mapstructure:"plugins"`
	} `mapstructure:"registry"`

	// Logging configuration
	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
		File   string `mapstructure:"file"`
	} `mapstructure:"logging"`

	// Plugin system
	Plugins []struct {
		Name    string                 `mapstructure:"name"`
		Enabled bool                   `mapstructure:"enabled"`
		Config  map[string]interface{} `mapstructure:"config"`
	} `mapstructure:"plugins"`

	// Docker runtime settings
	Docker struct {
		RestartPolicy string   `mapstructure:"restart_policy"`
		NetworkMode   string   `mapstructure:"network_mode"`
		Environment   []string `mapstructure:"environment"`
	} `mapstructure:"docker"`
}

// LoadConfig loads the configuration from a file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// If a config file is provided, read it
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
	}

	// Read from environment variables
	v.SetEnvPrefix("GALAXY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal the config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

// setDefaults sets default values for the configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.address", "0.0.0.0:50051")
	v.SetDefault("server.tls.enabled", false)
	v.SetDefault("server.max_connections", 500)
	v.SetDefault("server.resources.cpu_limit", "2")
	v.SetDefault("server.resources.memory_limit", "2Gi")

	// Registry defaults
	v.SetDefault("registry.allow_public_registration", true)
	v.SetDefault("registry.health_check_interval", "30s")
	v.SetDefault("registry.max_nodes", 100)
	v.SetDefault("registry.auto_deregister_after", 3)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Docker defaults
	v.SetDefault("docker.restart_policy", "always")
	v.SetDefault("docker.network_mode", "bridge")
}

// GetPluginConfigs extracts plugin configurations from the main config
func GetPluginConfigs(config *Config) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	// Add global plugins
	for _, plugin := range config.Plugins {
		if plugin.Enabled {
			result[plugin.Name] = plugin.Config
		}
	}

	// Add registry plugins
	for _, plugin := range config.Registry.Plugins {
		if plugin.Enabled {
			result[plugin.Name] = plugin.Config
		}
	}

	return result
}
