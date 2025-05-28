package registry

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"galaxy-node-pool/internal/config"
	"galaxy-node-pool/internal/plugin"
	pb "galaxy-node-pool/proto/pool"
)

// Registry implements the gRPC registry server with plugin support
type Registry struct {
	pb.UnimplementedRegistryServer
	mu             sync.RWMutex
	nodes          map[string]*pb.NodeInfo
	pluginManager  *plugin.PluginManager
	config         *config.Config
	missedHeartbeats map[string]int
	maxNodes       int
}

// NewRegistry creates a new registry server with the given configuration
func NewRegistry(cfg *config.Config, pluginMgr *plugin.PluginManager) *Registry {
	return &Registry{
		nodes:          make(map[string]*pb.NodeInfo),
		pluginManager:  pluginMgr,
		config:         cfg,
		missedHeartbeats: make(map[string]int),
		maxNodes:       cfg.Registry.MaxNodes,
	}
}

// Start initializes the registry and starts background tasks
func (r *Registry) Start(ctx context.Context) error {
	// Start health check goroutine
	go r.healthCheckLoop(ctx)

	// Initialize registry plugins
	for _, pluginCfg := range r.config.Registry.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		plg, err := r.pluginManager.Get(pluginCfg.Name)
		if err != nil {
			log.Printf("Warning: Registry plugin %s not found: %v", pluginCfg.Name, err)
			continue
		}

		// Check if it's a registry plugin
		if regPlugin, ok := plg.(plugin.RegistryPlugin); ok {
			if err := regPlugin.Initialize(pluginCfg.Config); err != nil {
				log.Printf("Warning: Failed to initialize registry plugin %s: %v", pluginCfg.Name, err)
			}
		}
	}

	log.Printf("Registry started with max nodes: %d", r.maxNodes)
	return nil
}

// healthCheckLoop periodically checks node health and removes unhealthy nodes
func (r *Registry) healthCheckLoop(ctx context.Context) {
	interval, err := time.ParseDuration(r.config.Registry.HealthCheckInterval)
	if err != nil {
		log.Printf("Invalid health check interval: %v, using default of 30s", err)
		interval = 30 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.checkNodeHealth()
		}
	}
}

// checkNodeHealth identifies and removes unhealthy nodes
func (r *Registry) checkNodeHealth() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for nodeID, missedCount := range r.missedHeartbeats {
		if missedCount >= r.config.Registry.AutoDeregisterAfter {
			// Node has missed too many heartbeats, deregister it
			log.Printf("Deregistering unhealthy node: %s (missed %d heartbeats)", nodeID, missedCount)
			
			// Call plugins for node deregistration
			for _, pluginCfg := range r.config.Registry.Plugins {
				if !pluginCfg.Enabled {
					continue
				}

				plg, err := r.pluginManager.Get(pluginCfg.Name)
				if err != nil {
					continue
				}

				if regPlugin, ok := plg.(plugin.RegistryPlugin); ok {
					regPlugin.OnNodeDeregister(nodeID)
				}
			}

			// Remove the node
			delete(r.nodes, nodeID)
			delete(r.missedHeartbeats, nodeID)
		} else {
			// Increment missed heartbeat count
			r.missedHeartbeats[nodeID] = missedCount + 1
		}
	}
}

// RegisterNode handles node registration requests
func (r *Registry) RegisterNode(ctx context.Context, req *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if we've reached the maximum number of nodes
	if len(r.nodes) >= r.maxNodes {
		return &pb.RegisterNodeResponse{
			Success: false,
			Message: fmt.Sprintf("Maximum number of nodes (%d) reached", r.maxNodes),
		}, nil
	}

	// Check if this is a private pool and the org is allowed
	if !r.config.Registry.AllowPublicRegistration && len(r.config.Registry.AllowedOrgs) > 0 {
		allowed := false
		for _, org := range r.config.Registry.AllowedOrgs {
			if org == req.Org {
				allowed = true
				break
			}
		}

		if !allowed {
			return &pb.RegisterNodeResponse{
				Success: false,
				Message: fmt.Sprintf("Organization %s not allowed in this private pool", req.Org),
			}, nil
		}
	}

	// Call plugins for node registration
	for _, pluginCfg := range r.config.Registry.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		plg, err := r.pluginManager.Get(pluginCfg.Name)
		if err != nil {
			continue
		}

		if regPlugin, ok := plg.(plugin.RegistryPlugin); ok {
			// Convert request to metadata map
			metadata := map[string]interface{}{
				"node_id": req.NodeId,
				"specialization": req.Specialization,
				"endpoint": req.Endpoint,
				"org": req.Org,
				"private_node": req.PrivateNode,
			}

			if err := regPlugin.OnNodeRegister(req.NodeId, metadata); err != nil {
				return &pb.RegisterNodeResponse{Success: false, Message: err.Error()}, nil
			}
		}
	}

	// Create and store the node
	node := &pb.NodeInfo{
		NodeId:         req.NodeId,
		Specialization: req.Specialization,
		Endpoint:       req.Endpoint,
		Org:            req.Org,
		PrivateNode:    req.PrivateNode,
		Status:         "healthy",
		RegisteredAt:   time.Now().Unix(),
	}

	r.nodes[req.NodeId] = node
	r.missedHeartbeats[req.NodeId] = 0 // Initialize heartbeat tracking

	log.Printf("Registered node: %s (%s) from org: %s", req.NodeId, req.Specialization, req.Org)
	return &pb.RegisterNodeResponse{Success: true, Message: "Node registered successfully"}, nil
}

// Heartbeat handles node heartbeat requests
func (r *Registry) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, ok := r.nodes[req.NodeId]
	if !ok {
		return &pb.HeartbeatResponse{Alive: false, Message: "Node not registered"}, nil
	}

	// Update node status
	node.Status = "healthy"
	node.LastHeartbeatAt = time.Now().Unix()
	
	// Reset missed heartbeat counter
	r.missedHeartbeats[req.NodeId] = 0

	// Call plugins for heartbeat
	for _, pluginCfg := range r.config.Registry.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		plg, err := r.pluginManager.Get(pluginCfg.Name)
		if err != nil {
			continue
		}

		if regPlugin, ok := plg.(plugin.RegistryPlugin); ok {
			regPlugin.OnNodeHeartbeat(req.NodeId)
		}
	}

	return &pb.HeartbeatResponse{Alive: true, Message: "Heartbeat acknowledged"}, nil
}

// ListNodes handles requests to list available nodes
func (r *Registry) ListNodes(ctx context.Context, req *pb.ListNodesRequest) (*pb.ListNodesResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Prepare filter for plugins
	filter := map[string]string{}
	if req.Specialization != "" {
		filter["specialization"] = req.Specialization
	}
	if req.Org != "" {
		filter["org"] = req.Org
	}

	// Call plugins for node listing
	for _, pluginCfg := range r.config.Registry.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		plg, err := r.pluginManager.Get(pluginCfg.Name)
		if err != nil {
			continue
		}

		if regPlugin, ok := plg.(plugin.RegistryPlugin); ok {
			regPlugin.OnNodeList(filter)
		}
	}

	// Filter nodes based on request
	var result []*pb.NodeInfo
	for _, node := range r.nodes {
		// Skip unhealthy nodes
		if node.Status != "healthy" {
			continue
		}

		// Apply specialization filter
		if req.Specialization != "" && node.Specialization != req.Specialization {
			continue
		}

		// Apply organization filter
		if req.Org != "" && node.Org != req.Org {
			continue
		}

		result = append(result, node)
	}

	return &pb.ListNodesResponse{Nodes: result}, nil
}

// GetNodeCount returns the current number of registered nodes
func (r *Registry) GetNodeCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.nodes)
}

// GetNodeByID retrieves a specific node by ID
func (r *Registry) GetNodeByID(nodeID string) (*pb.NodeInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	node, exists := r.nodes[nodeID]
	return node, exists
}

// Listen utility for main.go
func Listen(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
