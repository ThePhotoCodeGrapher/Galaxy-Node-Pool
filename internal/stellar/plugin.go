package stellar

import (
	"context"
	"fmt"
	"log"
	"sync"

	"galaxy-node-pool/internal/plugin"
)

// StellarPlugin implements the FederationPlugin interface
// for Stellar blockchain integration
type StellarPlugin struct {
	client       *StellarClient
	config       StellarConfig
	initialized  bool
	mu           sync.RWMutex
	poolMetadata map[string]interface{}
}

// StellarConfig holds the configuration for the Stellar plugin
type StellarConfig struct {
	HorizonURL        string `mapstructure:"horizon_url"`
	NetworkPassphrase string `mapstructure:"network_passphrase"`
	PoolSeed          string `mapstructure:"pool_seed"`
	MainNetAccount    string `mapstructure:"mainnet_account"`
	RegistrationFee   string `mapstructure:"registration_fee"`
	PoolDomain        string `mapstructure:"pool_domain"`
	StakerRewardPerc  int    `mapstructure:"staker_reward_percentage"`
}

// NewStellarPlugin creates a new Stellar plugin instance
func NewStellarPlugin() *StellarPlugin {
	return &StellarPlugin{
		initialized:  false,
		poolMetadata: make(map[string]interface{}),
	}
}

// Name returns the plugin name
func (p *StellarPlugin) Name() string {
	return "stellar-federation"
}

// Initialize sets up the plugin with its configuration
func (p *StellarPlugin) Initialize(rawConfig map[string]interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Parse configuration
	var config StellarConfig
	
	// Extract config values with defaults
	if horizonURL, ok := rawConfig["horizon_url"].(string); ok {
		config.HorizonURL = horizonURL
	} else {
		config.HorizonURL = "https://horizon-testnet.stellar.org" // Default to testnet
	}
	
	if networkPassphrase, ok := rawConfig["network_passphrase"].(string); ok {
		config.NetworkPassphrase = networkPassphrase
	} else {
		config.NetworkPassphrase = "Test SDF Network ; September 2015" // Default to testnet
	}
	
	// Required fields
	if poolSeed, ok := rawConfig["pool_seed"].(string); ok {
		config.PoolSeed = poolSeed
	} else {
		return fmt.Errorf("pool_seed is required for Stellar plugin")
	}
	
	if mainNetAccount, ok := rawConfig["mainnet_account"].(string); ok {
		config.MainNetAccount = mainNetAccount
	} else {
		return fmt.Errorf("mainnet_account is required for Stellar plugin")
	}
	
	if registrationFee, ok := rawConfig["registration_fee"].(string); ok {
		config.RegistrationFee = registrationFee
	} else {
		config.RegistrationFee = "5" // Default 5 XLM
	}
	
	if poolDomain, ok := rawConfig["pool_domain"].(string); ok {
		config.PoolDomain = poolDomain
	} else {
		return fmt.Errorf("pool_domain is required for Stellar plugin")
	}
	
	if stakerRewardPerc, ok := rawConfig["staker_reward_percentage"].(int); ok {
		config.StakerRewardPerc = stakerRewardPerc
	} else {
		config.StakerRewardPerc = 70 // Default 70% to stakers
	}

	// Create Stellar client
	client, err := NewStellarClient(config.HorizonURL, config.NetworkPassphrase, config.PoolSeed)
	if err != nil {
		return fmt.Errorf("failed to create Stellar client: %v", err)
	}

	p.client = client
	p.config = config
	p.initialized = true

	log.Printf("Stellar federation plugin initialized with horizon: %s", config.HorizonURL)
	return nil
}

// Shutdown gracefully stops the plugin
func (p *StellarPlugin) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.initialized = false
	log.Println("Stellar federation plugin shutdown complete")
	return nil
}

// RegisterWithMainNet registers this pool with the main net
func (p *StellarPlugin) RegisterWithMainNet(mainNetURL string, poolMetadata map[string]interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.initialized {
		return fmt.Errorf("stellar plugin not initialized")
	}

	// Store pool metadata for future use
	p.poolMetadata = poolMetadata

	// Register with main net via Stellar transaction
	err := p.client.RegisterPoolWithMainNet(
		p.config.MainNetAccount,
		p.config.PoolDomain,
		p.config.RegistrationFee,
	)
	if err != nil {
		return fmt.Errorf("failed to register with main net: %v", err)
	}

	log.Printf("Pool registered with main net via Stellar. Fee: %s XLM", p.config.RegistrationFee)
	return nil
}

// DiscoverPools finds other pools on the main net
func (p *StellarPlugin) DiscoverPools(filter map[string]string) ([]map[string]interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.initialized {
		return nil, fmt.Errorf("stellar plugin not initialized")
	}

	// This would typically query the Stellar network for pool registration transactions
	// For now, return a mock response
	pools := []map[string]interface{}{
		{
			"pool_id":   "main.galaxy.network",
			"endpoint":  "main.galaxy.network:50051",
			"public":    true,
			"location":  "global",
			"node_count": 120,
		},
		{
			"pool_id":   "eu.galaxy.network",
			"endpoint":  "eu.galaxy.network:50051",
			"public":    true,
			"location":  "eu-central",
			"node_count": 45,
		},
	}

	return pools, nil
}

// SyncWithPeers synchronizes state with peer pools
func (p *StellarPlugin) SyncWithPeers() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.initialized {
		return fmt.Errorf("stellar plugin not initialized")
	}

	// This would sync state with peer pools
	// For now, just log the action
	log.Println("Syncing with peer pools via Stellar network")
	return nil
}

// VerifyNodePayment verifies a node's payment for registration
func (p *StellarPlugin) VerifyNodePayment(nodeID string, nodeAccount string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.initialized {
		return false, fmt.Errorf("stellar plugin not initialized")
	}

	// Check if the node has paid the registration fee
	err := p.client.ProcessNodeRegistrationFee(nodeAccount, nodeID, p.config.RegistrationFee)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DistributeRewards distributes rewards to stakers
func (p *StellarPlugin) DistributeRewards(totalFees string, stakerAccounts []string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.initialized {
		return fmt.Errorf("stellar plugin not initialized")
	}

	// Calculate reward per staker
	// This is a simplified calculation - in reality, you'd want to consider stake amounts
	rewardPerStaker := "1.0" // Simplified for this example

	// Distribute rewards
	err := p.client.DistributeStakerRewards(stakerAccounts, rewardPerStaker)
	if err != nil {
		return fmt.Errorf("failed to distribute rewards: %v", err)
	}

	log.Printf("Distributed rewards to %d stakers. Total: %s XLM", len(stakerAccounts), totalFees)
	return nil
}

// Ensure StellarPlugin implements the FederationPlugin interface
var _ plugin.FederationPlugin = (*StellarPlugin)(nil)
