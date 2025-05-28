package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Provider is responsible for loading and managing configuration
type Provider struct {
	viper      *viper.Viper
	configPath string
	env        string
	mu         sync.RWMutex
	watchers   map[string][]func(string, interface{})
}

// NewProvider creates a new configuration provider
func NewProvider(configPath, env string) *Provider {
	v := viper.New()
	
	return &Provider{
		viper:      v,
		configPath: configPath,
		env:        env,
		watchers:   make(map[string][]func(string, interface{})),
	}
}

// Load loads configuration from files and environment variables
func (p *Provider) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Set up viper
	p.viper.SetConfigFile(p.configPath)
	
	// Read the main config file
	if err := p.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Load environment-specific config if available
	if p.env != "" {
		ext := filepath.Ext(p.configPath)
		baseConfigPath := p.configPath[0:len(p.configPath)-len(ext)]
		envConfigPath := fmt.Sprintf("%s.%s%s", baseConfigPath, p.env, ext)
		
		if _, err := os.Stat(envConfigPath); err == nil {
			p.viper.SetConfigFile(envConfigPath)
			if err := p.viper.MergeInConfig(); err != nil {
				return fmt.Errorf("failed to merge environment config: %v", err)
			}
		}
	}

	// Load from environment variables
	p.viper.SetEnvPrefix("GALAXY")
	p.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	p.viper.AutomaticEnv()

	return nil
}

// Get retrieves a configuration value by key
func (p *Provider) Get(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.Get(key)
}

// GetString retrieves a string configuration value by key
func (p *Provider) GetString(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetString(key)
}

// GetInt retrieves an integer configuration value by key
func (p *Provider) GetInt(key string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetInt(key)
}

// GetBool retrieves a boolean configuration value by key
func (p *Provider) GetBool(key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetBool(key)
}

// GetStringSlice retrieves a string slice configuration value by key
func (p *Provider) GetStringSlice(key string) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetStringSlice(key)
}

// GetStringMap retrieves a string map configuration value by key
func (p *Provider) GetStringMap(key string) map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetStringMap(key)
}

// Set sets a configuration value
func (p *Provider) Set(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.viper.Set(key, value)
	
	// Notify watchers
	if watchers, exists := p.watchers[key]; exists {
		for _, watcher := range watchers {
			go watcher(key, value)
		}
	}
}

// Watch registers a function to be called when a configuration value changes
func (p *Provider) Watch(key string, callback func(string, interface{})) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if _, exists := p.watchers[key]; !exists {
		p.watchers[key] = make([]func(string, interface{}), 0)
	}
	
	p.watchers[key] = append(p.watchers[key], callback)
}

// UnmarshalKey unmarshals a specific key into a struct
func (p *Provider) UnmarshalKey(key string, rawVal interface{}) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.UnmarshalKey(key, rawVal)
}

// Unmarshal unmarshals the entire configuration into a struct
func (p *Provider) Unmarshal(rawVal interface{}) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.Unmarshal(rawVal)
}

// LoadConfig is a convenience function to load the entire configuration into a Config struct
func (p *Provider) LoadConfig() (*Config, error) {
	var config Config
	if err := p.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &config, nil
}
