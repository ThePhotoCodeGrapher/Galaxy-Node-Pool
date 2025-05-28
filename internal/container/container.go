package container

import (
	"fmt"
	"sync"
)

// ServiceContainer is a dependency injection container for services
type ServiceContainer struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewServiceContainer creates a new service container
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services: make(map[string]interface{}),
	}
}

// Register adds a service to the container
func (c *ServiceContainer) Register(name string, service interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	c.services[name] = service
	return nil
}

// Get retrieves a service from the container
func (c *ServiceContainer) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	return service, nil
}

// GetTyped retrieves a service from the container and casts it to the expected type
func (c *ServiceContainer) GetTyped(name string, target interface{}) error {
	service, err := c.Get(name)
	if err != nil {
		return err
	}

	// Use type assertion to check if service can be assigned to target
	targetPtr, ok := target.(*interface{})
	if !ok {
		return fmt.Errorf("target must be a pointer to interface{}")
	}

	*targetPtr = service
	return nil
}

// Has checks if a service exists in the container
func (c *ServiceContainer) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.services[name]
	return exists
}

// Remove removes a service from the container
func (c *ServiceContainer) Remove(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.services[name]; !exists {
		return fmt.Errorf("service %s not found", name)
	}

	delete(c.services, name)
	return nil
}

// GetAll returns all registered service names
func (c *ServiceContainer) GetAll() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.services))
	for name := range c.services {
		names = append(names, name)
	}

	return names
}
