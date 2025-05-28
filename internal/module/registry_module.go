package module

import (
	"context"
	"fmt"
	"log"

	"galaxy-node-pool/internal/container"
	"galaxy-node-pool/internal/event"
	"galaxy-node-pool/internal/registry"
	"galaxy-node-pool/internal/config"
)

// RegistryModule implements the Module interface for the registry component
type RegistryModule struct {
	name        string
	description string
	version     string
	registry    *registry.Registry
	config      *config.Config
}

// NewRegistryModule creates a new registry module
func NewRegistryModule(cfg *config.Config) *RegistryModule {
	return &RegistryModule{
		name:        "registry",
		description: "Node registry and discovery service",
		version:     "1.0.0",
		config:      cfg,
	}
}

// Name returns the module name
func (m *RegistryModule) Name() string {
	return m.name
}

// Description returns the module description
func (m *RegistryModule) Description() string {
	return m.description
}

// Version returns the module version
func (m *RegistryModule) Version() string {
	return m.version
}

// Dependencies returns module dependencies
func (m *RegistryModule) Dependencies() []string {
	// Registry has no module dependencies
	return []string{}
}

// Load initializes the module
func (m *RegistryModule) Load(ctx context.Context, container *container.ServiceContainer, dispatcher *event.EventDispatcher) error {
	log.Printf("Loading registry module...")

	// Get plugin manager from container
	pluginManagerInterface, err := container.Get("plugin_manager")
	if err != nil {
		return fmt.Errorf("failed to get plugin manager: %v", err)
	}
	
	pluginManager, ok := pluginManagerInterface.(*registry.PluginManager)
	if !ok {
		return fmt.Errorf("invalid plugin manager type")
	}

	// Create registry
	reg := registry.NewRegistry(m.config, pluginManager)
	m.registry = reg

	// Register registry with container
	if err := container.Register("registry", reg); err != nil {
		return fmt.Errorf("failed to register registry with container: %v", err)
	}

	// Subscribe to events
	dispatcher.Subscribe("node.registered", func(e event.Event) {
		nodeID, ok := e.Data["node_id"].(string)
		if !ok {
			return
		}
		log.Printf("Event: Node registered: %s", nodeID)
	})

	dispatcher.Subscribe("node.heartbeat", func(e event.Event) {
		nodeID, ok := e.Data["node_id"].(string)
		if !ok {
			return
		}
		log.Printf("Event: Node heartbeat: %s", nodeID)
	})

	// Start the registry
	if err := reg.Start(ctx); err != nil {
		return fmt.Errorf("failed to start registry: %v", err)
	}

	log.Printf("Registry module loaded successfully")
	return nil
}

// Unload cleans up the module
func (m *RegistryModule) Unload(ctx context.Context) error {
	log.Printf("Unloading registry module...")
	
	// Registry doesn't need explicit cleanup as it will be stopped
	// when the context is canceled
	
	log.Printf("Registry module unloaded successfully")
	return nil
}

// GetRegistry returns the registry instance
func (m *RegistryModule) GetRegistry() *registry.Registry {
	return m.registry
}
