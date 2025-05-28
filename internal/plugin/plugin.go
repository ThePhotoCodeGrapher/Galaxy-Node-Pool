package plugin

import (
	"fmt"
	"log"
	"plugin"
	"sync"
)

// PluginManager handles the loading and management of plugins
type PluginManager struct {
	plugins     map[string]interface{}
	mu          sync.RWMutex
	initialized bool
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins:     make(map[string]interface{}),
		initialized: false,
	}
}

// Initialize loads all plugins from the specified directory
func (pm *PluginManager) Initialize(pluginDir string, configs map[string]map[string]interface{}) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.initialized {
		return fmt.Errorf("plugin manager already initialized")
	}

	// Load plugins from directory
	// For each plugin file in pluginDir:
	//   1. Open the plugin
	//   2. Look up the "New" symbol
	//   3. Create a new instance with the plugin config
	//   4. Register the plugin

	pm.initialized = true
	log.Printf("Plugin manager initialized with directory: %s", pluginDir)
	return nil
}

// Register adds a plugin to the manager
func (pm *PluginManager) Register(name string, instance interface{}) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	pm.plugins[name] = instance
	log.Printf("Plugin registered: %s", name)
	return nil
}

// Get retrieves a plugin by name
func (pm *PluginManager) Get(name string) (interface{}, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// LoadPlugin loads a single plugin from a file
func (pm *PluginManager) LoadPlugin(path string, config map[string]interface{}) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %v", path, err)
	}

	newFunc, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'New' function: %v", path, err)
	}

	// Call the New function to create a plugin instance
	constructor, ok := newFunc.(func(map[string]interface{}) (interface{}, error))
	if !ok {
		return fmt.Errorf("plugin %s has invalid 'New' function signature", path)
	}

	instance, err := constructor(config)
	if err != nil {
		return fmt.Errorf("failed to create plugin instance for %s: %v", path, err)
	}

	// Extract plugin name from the instance
	nameMethod, err := p.Lookup("Name")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'Name' function: %v", path, err)
	}

	nameFunc, ok := nameMethod.(func() string)
	if !ok {
		return fmt.Errorf("plugin %s has invalid 'Name' function signature", path)
	}

	name := nameFunc()
	return pm.Register(name, instance)
}
