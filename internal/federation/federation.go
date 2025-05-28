package federation

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"galaxy-node-pool/internal/config"
	"galaxy-node-pool/internal/plugin"
)

// Federation manages the pool's connection to the main net and other pools
type Federation struct {
	config        *config.Config
	pluginManager *plugin.PluginManager
	federationPlugin plugin.FederationPlugin
	mu            sync.RWMutex
	isRegistered  bool
	lastSyncTime  time.Time
	peerPools     []map[string]interface{}
}

// NewFederation creates a new federation manager
func NewFederation(cfg *config.Config, pluginMgr *plugin.PluginManager) (*Federation, error) {
	return &Federation{
		config:        cfg,
		pluginManager: pluginMgr,
		isRegistered:  false,
		peerPools:     make([]map[string]interface{}, 0),
	}, nil
}

// Initialize sets up the federation manager
func (f *Federation) Initialize(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Find and initialize the federation plugin
	for _, pluginCfg := range f.config.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		plg, err := f.pluginManager.Get(pluginCfg.Name)
		if err != nil {
			continue
		}

		// Check if it's a federation plugin
		if fedPlugin, ok := plg.(plugin.FederationPlugin); ok {
			f.federationPlugin = fedPlugin
			log.Printf("Federation plugin found: %s", pluginCfg.Name)
			break
		}
	}

	if f.federationPlugin == nil {
		return fmt.Errorf("no federation plugin found, cannot connect to main net")
	}

	return nil
}

// RegisterWithMainNet registers this pool with the main net
func (f *Federation) RegisterWithMainNet() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.federationPlugin == nil {
		return fmt.Errorf("no federation plugin available")
	}

	if f.isRegistered {
		log.Println("Pool already registered with main net")
		return nil
	}

	// Prepare pool metadata
	poolMetadata := map[string]interface{}{
		"domain":     f.config.Domain.DomainName,
		"address":    f.config.Server.Address,
		"location":   f.config.Server.Location,
		"public":     f.config.Registry.AllowPublicRegistration,
		"max_nodes":  f.config.Registry.MaxNodes,
	}

	// Register with main net
	err := f.federationPlugin.RegisterWithMainNet(f.config.MainNet.RegistryAddress, poolMetadata)
	if err != nil {
		return fmt.Errorf("failed to register with main net: %v", err)
	}

	f.isRegistered = true
	log.Printf("Pool registered with main net: %s", f.config.MainNet.RegistryAddress)
	return nil
}

// DiscoverPools finds other pools on the main net
func (f *Federation) DiscoverPools(filter map[string]string) ([]map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.federationPlugin == nil {
		return nil, fmt.Errorf("no federation plugin available")
	}

	// Discover pools
	pools, err := f.federationPlugin.DiscoverPools(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to discover pools: %v", err)
	}

	f.peerPools = pools
	return pools, nil
}

// SyncWithPeers synchronizes state with peer pools
func (f *Federation) SyncWithPeers() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.federationPlugin == nil {
		return fmt.Errorf("no federation plugin available")
	}

	// Sync with peers
	err := f.federationPlugin.SyncWithPeers()
	if err != nil {
		return fmt.Errorf("failed to sync with peers: %v", err)
	}

	f.lastSyncTime = time.Now()
	log.Printf("Synced with peer pools at %s", f.lastSyncTime.Format(time.RFC3339))
	return nil
}

// StartSyncLoop starts a background loop to periodically sync with peers
func (f *Federation) StartSyncLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Federation sync loop stopped")
				return
			case <-ticker.C:
				if err := f.SyncWithPeers(); err != nil {
					log.Printf("Error syncing with peers: %v", err)
				}
			}
		}
	}()

	log.Printf("Federation sync loop started with interval: %s", interval)
}

// IsRegistered returns whether the pool is registered with the main net
func (f *Federation) IsRegistered() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.isRegistered
}

// GetPeerPools returns the list of discovered peer pools
func (f *Federation) GetPeerPools() []map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.peerPools
}

// GetLastSyncTime returns the time of the last successful sync
func (f *Federation) GetLastSyncTime() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.lastSyncTime
}
