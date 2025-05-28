package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"galaxy-node-pool/internal/container"
	"galaxy-node-pool/internal/event"
)

// ServiceState represents the state of a service
type ServiceState int

const (
	// ServiceStopped indicates the service is stopped
	ServiceStopped ServiceState = iota
	// ServiceStarting indicates the service is starting
	ServiceStarting
	// ServiceRunning indicates the service is running
	ServiceRunning
	// ServiceStopping indicates the service is stopping
	ServiceStopping
	// ServiceFailed indicates the service failed to start or run
	ServiceFailed
)

// Service defines the interface that all services must implement
type Service interface {
	// Name returns the unique name of the service
	Name() string
	
	// Start starts the service
	Start(ctx context.Context) error
	
	// Stop stops the service
	Stop(ctx context.Context) error
	
	// State returns the current state of the service
	State() ServiceState
	
	// Dependencies returns the names of services this service depends on
	Dependencies() []string
}

// ServiceManager manages the lifecycle of all services
type ServiceManager struct {
	container     *container.ServiceContainer
	dispatcher    *event.EventDispatcher
	services      map[string]Service
	states        map[string]ServiceState
	dependencies  map[string][]string
	dependents    map[string][]string
	startOrder    []string
	stopOrder     []string
	mu            sync.RWMutex
}

// NewServiceManager creates a new service manager
func NewServiceManager(container *container.ServiceContainer, dispatcher *event.EventDispatcher) *ServiceManager {
	return &ServiceManager{
		container:    container,
		dispatcher:   dispatcher,
		services:     make(map[string]Service),
		states:       make(map[string]ServiceState),
		dependencies: make(map[string][]string),
		dependents:   make(map[string][]string),
		startOrder:   make([]string, 0),
		stopOrder:    make([]string, 0),
	}
}

// Register adds a service to the manager
func (m *ServiceManager) Register(service Service) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := service.Name()
	if _, exists := m.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	// Register the service
	m.services[name] = service
	m.states[name] = ServiceStopped
	m.dependencies[name] = service.Dependencies()
	
	// Register with the container
	if err := m.container.Register(name, service); err != nil {
		return fmt.Errorf("failed to register service %s with container: %v", name, err)
	}

	// Update dependents map
	for _, dep := range service.Dependencies() {
		if _, exists := m.dependents[dep]; !exists {
			m.dependents[dep] = make([]string, 0)
		}
		m.dependents[dep] = append(m.dependents[dep], name)
	}

	// Recalculate start and stop order
	if err := m.calculateOrder(); err != nil {
		return err
	}

	log.Printf("Service %s registered with dependencies: %v", name, service.Dependencies())
	return nil
}

// Start starts all services in dependency order
func (m *ServiceManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, name := range m.startOrder {
		service := m.services[name]
		
		// Skip services that are already running
		if m.states[name] == ServiceRunning {
			continue
		}
		
		// Update state
		m.states[name] = ServiceStarting
		
		// Dispatch event
		m.dispatcher.Dispatch(event.Event{
			Name: "service.starting",
			Data: map[string]interface{}{
				"service": name,
			},
		})
		
		// Start the service
		log.Printf("Starting service: %s", name)
		if err := service.Start(ctx); err != nil {
			m.states[name] = ServiceFailed
			m.dispatcher.Dispatch(event.Event{
				Name: "service.failed",
				Data: map[string]interface{}{
					"service": name,
					"error":   err.Error(),
				},
			})
			return fmt.Errorf("failed to start service %s: %v", name, err)
		}
		
		// Update state
		m.states[name] = ServiceRunning
		
		// Dispatch event
		m.dispatcher.Dispatch(event.Event{
			Name: "service.started",
			Data: map[string]interface{}{
				"service": name,
			},
		})
	}

	return nil
}

// Stop stops all services in reverse dependency order
func (m *ServiceManager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastError error

	for _, name := range m.stopOrder {
		service := m.services[name]
		
		// Skip services that are already stopped
		if m.states[name] == ServiceStopped {
			continue
		}
		
		// Update state
		m.states[name] = ServiceStopping
		
		// Dispatch event
		m.dispatcher.Dispatch(event.Event{
			Name: "service.stopping",
			Data: map[string]interface{}{
				"service": name,
			},
		})
		
		// Stop the service
		log.Printf("Stopping service: %s", name)
		if err := service.Stop(ctx); err != nil {
			lastError = err
			m.states[name] = ServiceFailed
			m.dispatcher.Dispatch(event.Event{
				Name: "service.failed",
				Data: map[string]interface{}{
					"service": name,
					"error":   err.Error(),
				},
			})
			log.Printf("Failed to stop service %s: %v", name, err)
			continue
		}
		
		// Update state
		m.states[name] = ServiceStopped
		
		// Dispatch event
		m.dispatcher.Dispatch(event.Event{
			Name: "service.stopped",
			Data: map[string]interface{}{
				"service": name,
			},
		})
	}

	return lastError
}

// StartService starts a specific service and its dependencies
func (m *ServiceManager) StartService(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if service exists
	service, exists := m.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	// Start dependencies first
	for _, dep := range service.Dependencies() {
		depService, exists := m.services[dep]
		if !exists {
			return fmt.Errorf("dependency %s for service %s not found", dep, name)
		}

		if m.states[dep] != ServiceRunning {
			if err := m.startServiceLocked(ctx, depService); err != nil {
				return err
			}
		}
	}

	// Start the service
	return m.startServiceLocked(ctx, service)
}

// startServiceLocked starts a service (must be called with lock held)
func (m *ServiceManager) startServiceLocked(ctx context.Context, service Service) error {
	name := service.Name()
	
	// Skip if already running
	if m.states[name] == ServiceRunning {
		return nil
	}
	
	// Update state
	m.states[name] = ServiceStarting
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "service.starting",
		Data: map[string]interface{}{
			"service": name,
		},
	})
	
	// Start the service
	log.Printf("Starting service: %s", name)
	if err := service.Start(ctx); err != nil {
		m.states[name] = ServiceFailed
		m.dispatcher.Dispatch(event.Event{
			Name: "service.failed",
			Data: map[string]interface{}{
				"service": name,
				"error":   err.Error(),
			},
		})
		return fmt.Errorf("failed to start service %s: %v", name, err)
	}
	
	// Update state
	m.states[name] = ServiceRunning
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "service.started",
		Data: map[string]interface{}{
			"service": name,
		},
	})
	
	return nil
}

// StopService stops a specific service and its dependents
func (m *ServiceManager) StopService(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if service exists
	service, exists := m.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	// Stop dependents first
	for _, dep := range m.dependents[name] {
		depService, exists := m.services[dep]
		if !exists {
			continue
		}

		if m.states[dep] == ServiceRunning {
			if err := m.stopServiceLocked(ctx, depService); err != nil {
				return err
			}
		}
	}

	// Stop the service
	return m.stopServiceLocked(ctx, service)
}

// stopServiceLocked stops a service (must be called with lock held)
func (m *ServiceManager) stopServiceLocked(ctx context.Context, service Service) error {
	name := service.Name()
	
	// Skip if already stopped
	if m.states[name] == ServiceStopped {
		return nil
	}
	
	// Update state
	m.states[name] = ServiceStopping
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "service.stopping",
		Data: map[string]interface{}{
			"service": name,
		},
	})
	
	// Stop the service
	log.Printf("Stopping service: %s", name)
	if err := service.Stop(ctx); err != nil {
		m.states[name] = ServiceFailed
		m.dispatcher.Dispatch(event.Event{
			Name: "service.failed",
			Data: map[string]interface{}{
				"service": name,
				"error":   err.Error(),
			},
		})
		return fmt.Errorf("failed to stop service %s: %v", name, err)
	}
	
	// Update state
	m.states[name] = ServiceStopped
	
	// Dispatch event
	m.dispatcher.Dispatch(event.Event{
		Name: "service.stopped",
		Data: map[string]interface{}{
			"service": name,
		},
	})
	
	return nil
}

// GetService retrieves a service by name
func (m *ServiceManager) GetService(name string) (Service, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	service, exists := m.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// GetServiceState returns the state of a service
func (m *ServiceManager) GetServiceState(name string) (ServiceState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[name]
	if !exists {
		return ServiceStopped, fmt.Errorf("service %s not found", name)
	}

	return state, nil
}

// GetAllServices returns all registered services
func (m *ServiceManager) GetAllServices() map[string]Service {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	services := make(map[string]Service, len(m.services))
	for name, service := range m.services {
		services[name] = service
	}

	return services
}

// calculateOrder calculates the start and stop order of services based on dependencies
func (m *ServiceManager) calculateOrder() error {
	// Reset orders
	m.startOrder = make([]string, 0, len(m.services))
	m.stopOrder = make([]string, 0, len(m.services))

	// Create a copy of dependencies for topological sort
	deps := make(map[string]map[string]bool)
	for name, dependencies := range m.dependencies {
		deps[name] = make(map[string]bool)
		for _, dep := range dependencies {
			if _, exists := m.services[dep]; !exists {
				return fmt.Errorf("service %s depends on unknown service %s", name, dep)
			}
			deps[name][dep] = true
		}
	}

	// Topological sort for start order
	for len(deps) > 0 {
		// Find services with no dependencies
		ready := make([]string, 0)
		for name, dependencies := range deps {
			if len(dependencies) == 0 {
				ready = append(ready, name)
			}
		}

		if len(ready) == 0 {
			// Circular dependency detected
			return fmt.Errorf("circular dependency detected in services")
		}

		// Add ready services to start order
		for _, name := range ready {
			m.startOrder = append(m.startOrder, name)
			delete(deps, name)

			// Remove this service from dependencies of other services
			for _, dependencies := range deps {
				delete(dependencies, name)
			}
		}
	}

	// Stop order is reverse of start order
	m.stopOrder = make([]string, len(m.startOrder))
	for i, name := range m.startOrder {
		m.stopOrder[len(m.startOrder)-1-i] = name
	}

	return nil
}
