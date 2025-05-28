package loader

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"galaxy-node-pool/internal/container"
)

// PluginLoader loads plugins from a directory and registers them with the container
type PluginLoader struct {
	container *container.ServiceContainer
	loadedPlugins map[string]interface{}
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader(container *container.ServiceContainer) *PluginLoader {
	return &PluginLoader{
		container: container,
		loadedPlugins: make(map[string]interface{}),
	}
}

// LoadPluginsFromDir loads all plugins from a directory
func (l *PluginLoader) LoadPluginsFromDir(dir string, configs map[string]map[string]interface{}) error {
	// Read all files in the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	// Load each plugin file
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".so") {
			continue
		}

		pluginPath := filepath.Join(dir, file.Name())
		pluginName := strings.TrimSuffix(file.Name(), ".so")

		// Check if we have a config for this plugin
		config, hasConfig := configs[pluginName]
		if !hasConfig {
			log.Printf("No configuration found for plugin %s, skipping", pluginName)
			continue
		}

		// Load the plugin
		if err := l.LoadPlugin(pluginPath, pluginName, config); err != nil {
			log.Printf("Failed to load plugin %s: %v", pluginName, err)
			continue
		}
	}

	return nil
}

// LoadPlugin loads a single plugin
func (l *PluginLoader) LoadPlugin(path, name string, config map[string]interface{}) error {
	// Open the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %v", path, err)
	}

	// Look up the plugin factory function
	newFunc, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'New' function: %v", path, err)
	}

	// Call the factory function to create the plugin instance
	constructor, ok := newFunc.(func(map[string]interface{}) (interface{}, error))
	if !ok {
		return fmt.Errorf("plugin %s has invalid 'New' function signature", path)
	}

	// Create the plugin instance
	instance, err := constructor(config)
	if err != nil {
		return fmt.Errorf("failed to create plugin instance for %s: %v", path, err)
	}

	// Register the plugin with the container
	if err := l.container.Register(name, instance); err != nil {
		return fmt.Errorf("failed to register plugin %s: %v", name, err)
	}

	// Store in loaded plugins map
	l.loadedPlugins[name] = instance
	log.Printf("Plugin %s loaded and registered", name)

	return nil
}

// GetLoadedPlugins returns all loaded plugins
func (l *PluginLoader) GetLoadedPlugins() map[string]interface{} {
	return l.loadedPlugins
}

// GetPlugin retrieves a loaded plugin by name
func (l *PluginLoader) GetPlugin(name string) (interface{}, error) {
	plugin, exists := l.loadedPlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return plugin, nil
}
