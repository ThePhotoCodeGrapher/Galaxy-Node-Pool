package module

import (
	"context"
	"fmt"
	"log"

	"galaxy-node-pool/internal/container"
	"galaxy-node-pool/internal/event"
	"galaxy-node-pool/internal/service"
)

// Module represents a self-contained component that can be loaded and unloaded
type Module interface {
	// Name returns the unique name of the module
	Name() string
	
	// Description returns a description of the module
	Description() string
	
	// Version returns the module version
	Version() string
	
	// Dependencies returns the names of modules this module depends on
	Dependencies() []string
	
	// Load initializes the module and registers its services
	Load(ctx context.Context, container *container.ServiceContainer, dispatcher *event.EventDispatcher) error
	
	// Unload cleans up the module's resources
	Unload(ctx context.Context) error
}

// ModuleManager manages the lifecycle of modules
type ModuleManager struct {
	container      *container.ServiceContainer
	dispatcher     *event.EventDispatcher
	serviceManager *service.ServiceManager
	modules        map[string]Module
	loadedModules  map[string]bool
}

// NewModuleManager creates a new module manager
func NewModuleManager(
	container *container.ServiceContainer,
	dispatcher *event.EventDispatcher,
	serviceManager *service.ServiceManager,
) *ModuleManager {
	return &ModuleManager{
		container:      container,
		dispatcher:     dispatcher,
		serviceManager: serviceManager,
		modules:        make(map[string]Module),
		loadedModules:  make(map[string]bool),
	}
}

// Register adds a module to the manager
func (m *ModuleManager) Register(module Module) error {
	name := module.Name()
	if _, exists := m.modules[name]; exists {
		return fmt.Errorf("module %s already registered", name)
	}

	m.modules[name] = module
	m.loadedModules[name] = false

	log.Printf("Module %s registered (version: %s)", name, module.Version())
	return nil
}

// Load loads a module and its dependencies
func (m *ModuleManager) Load(ctx context.Context, name string) error {
	// Check if module exists
	module, exists := m.modules[name]
	if !exists {
		return fmt.Errorf("module %s not found", name)
	}

	// Check if module is already loaded
	if m.loadedModules[name] {
		return nil
	}

	// Load dependencies first
	for _, dep := range module.Dependencies() {
		if err := m.Load(ctx, dep); err != nil {
			return fmt.Errorf("failed to load dependency %s for module %s: %v", dep, name, err)
		}
	}

	// Load the module
	log.Printf("Loading module: %s", name)
	if err := module.Load(ctx, m.container, m.dispatcher); err != nil {
		return fmt.Errorf("failed to load module %s: %v", name, err)
	}

	m.loadedModules[name] = true
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "module.loaded",
		Data: map[string]interface{}{
			"module": name,
		},
	})

	log.Printf("Module %s loaded successfully", name)
	return nil
}

// LoadAll loads all registered modules
func (m *ModuleManager) LoadAll(ctx context.Context) error {
	for name := range m.modules {
		if err := m.Load(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

// Unload unloads a module and its dependents
func (m *ModuleManager) Unload(ctx context.Context, name string) error {
	// Check if module exists
	module, exists := m.modules[name]
	if !exists {
		return fmt.Errorf("module %s not found", name)
	}

	// Check if module is already unloaded
	if !m.loadedModules[name] {
		return nil
	}

	// Find modules that depend on this one
	var dependents []string
	for modName, mod := range m.modules {
		for _, dep := range mod.Dependencies() {
			if dep == name && m.loadedModules[modName] {
				dependents = append(dependents, modName)
				break
			}
		}
	}

	// Unload dependents first
	for _, dep := range dependents {
		if err := m.Unload(ctx, dep); err != nil {
			return fmt.Errorf("failed to unload dependent module %s: %v", dep, err)
		}
	}

	// Unload the module
	log.Printf("Unloading module: %s", name)
	if err := module.Unload(ctx); err != nil {
		return fmt.Errorf("failed to unload module %s: %v", name, err)
	}

	m.loadedModules[name] = false
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "module.unloaded",
		Data: map[string]interface{}{
			"module": name,
		},
	})

	log.Printf("Module %s unloaded successfully", name)
	return nil
}

// UnloadAll unloads all loaded modules
func (m *ModuleManager) UnloadAll(ctx context.Context) error {
	// Find all loaded modules
	var loaded []string
	for name, isLoaded := range m.loadedModules {
		if isLoaded {
			loaded = append(loaded, name)
		}
	}

	// Unload in reverse dependency order
	// This is a simplified approach - a proper topological sort would be better
	for len(loaded) > 0 {
		for i := len(loaded) - 1; i >= 0; i-- {
			name := loaded[i]
			
			// Check if this module has any loaded dependents
			hasDependents := false
			for _, otherName := range loaded {
				if otherName == name {
					continue
				}
				
				otherModule := m.modules[otherName]
				for _, dep := range otherModule.Dependencies() {
					if dep == name {
						hasDependents = true
						break
					}
				}
				
				if hasDependents {
					break
				}
			}
			
			if !hasDependents {
				// Safe to unload
				if err := m.Unload(ctx, name); err != nil {
					return err
				}
				
				// Remove from loaded list
				loaded = append(loaded[:i], loaded[i+1:]...)
				break
			}
		}
	}

	return nil
}

// IsLoaded checks if a module is loaded
func (m *ModuleManager) IsLoaded(name string) bool {
	isLoaded, exists := m.loadedModules[name]
	return exists && isLoaded
}

// GetModule retrieves a module by name
func (m *ModuleManager) GetModule(name string) (Module, error) {
	module, exists := m.modules[name]
	if !exists {
		return nil, fmt.Errorf("module %s not found", name)
	}
	return module, nil
}

// GetAllModules returns all registered modules
func (m *ModuleManager) GetAllModules() map[string]Module {
	// Create a copy to avoid race conditions
	modules := make(map[string]Module, len(m.modules))
	for name, module := range m.modules {
		modules[name] = module
	}
	return modules
}

// GetLoadedModules returns all loaded modules
func (m *ModuleManager) GetLoadedModules() map[string]Module {
	// Create a copy to avoid race conditions
	modules := make(map[string]Module)
	for name, isLoaded := range m.loadedModules {
		if isLoaded {
			modules[name] = m.modules[name]
		}
	}
	return modules
}
