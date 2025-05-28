package plugin

import (
	"context"
	"net/http"
)

// Plugin is the base interface that all plugins must implement
type Plugin interface {
	// Name returns the unique name of the plugin
	Name() string
	
	// Initialize sets up the plugin with its configuration
	Initialize(config map[string]interface{}) error
	
	// Shutdown gracefully stops the plugin
	Shutdown(ctx context.Context) error
}

// AuthPlugin provides authentication and authorization services
type AuthPlugin interface {
	Plugin
	
	// Authenticate verifies credentials and returns a user ID or error
	Authenticate(credentials map[string]string) (string, error)
	
	// Authorize checks if a user has permission for an action
	Authorize(userID string, resource string, action string) (bool, error)
	
	// Middleware returns an HTTP middleware for auth
	Middleware() func(http.Handler) http.Handler
}

// MetricsPlugin provides metrics collection and reporting
type MetricsPlugin interface {
	Plugin
	
	// RecordMetric records a named metric with value and labels
	RecordMetric(name string, value float64, labels map[string]string) error
	
	// GetMetrics returns all current metrics
	GetMetrics() (map[string]interface{}, error)
	
	// Middleware returns an HTTP middleware for metrics
	Middleware() func(http.Handler) http.Handler
}

// StoragePlugin provides persistence for the registry
type StoragePlugin interface {
	Plugin
	
	// Store persists a key-value pair
	Store(key string, value interface{}) error
	
	// Retrieve gets a value by key
	Retrieve(key string) (interface{}, error)
	
	// Delete removes a key-value pair
	Delete(key string) error
	
	// List returns all keys with optional prefix
	List(prefix string) ([]string, error)
}

// RegistryPlugin hooks into the node registry operations
type RegistryPlugin interface {
	Plugin
	
	// OnNodeRegister is called when a node registers
	OnNodeRegister(nodeID string, metadata map[string]interface{}) error
	
	// OnNodeHeartbeat is called when a node sends a heartbeat
	OnNodeHeartbeat(nodeID string) error
	
	// OnNodeList is called when nodes are listed
	OnNodeList(filter map[string]string) error
	
	// OnNodeDeregister is called when a node is removed
	OnNodeDeregister(nodeID string) error
}

// FederationPlugin handles main net and cross-pool communication
type FederationPlugin interface {
	Plugin
	
	// RegisterWithMainNet registers this pool with the main net
	RegisterWithMainNet(mainNetURL string, poolMetadata map[string]interface{}) error
	
	// DiscoverPools finds other pools on the main net
	DiscoverPools(filter map[string]string) ([]map[string]interface{}, error)
	
	// SyncWithPeers synchronizes state with peer pools
	SyncWithPeers() error
	
	// VerifyNodePayment verifies a node's payment for registration (for blockchain-based federation)
	VerifyNodePayment(nodeID string, nodeAccount string) (bool, error)
	
	// DistributeRewards distributes rewards to stakers (for blockchain-based federation)
	DistributeRewards(totalFees string, stakerAccounts []string) error
}
