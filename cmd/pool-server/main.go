package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"galaxy-node-pool/internal/config"
	"galaxy-node-pool/internal/federation"
	"galaxy-node-pool/internal/plugin"
	"galaxy-node-pool/internal/registry"
	"galaxy-node-pool/internal/stellar"
	pb "galaxy-node-pool/proto/pool"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/example.yaml", "Path to configuration file")
	pluginDir := flag.String("plugins", "./plugins", "Directory containing plugins")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// Set up logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// Load configuration
	log.Printf("Loading configuration from %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create plugin manager
	log.Printf("Initializing plugin manager with directory: %s", *pluginDir)
	pluginManager := plugin.NewPluginManager()

	// Extract plugin configs
	pluginConfigs := config.GetPluginConfigs(cfg)

	// Initialize plugin manager
	if err := pluginManager.Initialize(*pluginDir, pluginConfigs); err != nil {
		log.Printf("Warning: Failed to initialize plugin manager: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, initiating shutdown", sig)
		cancel()
	}()

	// Prepare gRPC server options
	var opts []grpc.ServerOption
	if cfg.Server.TLS.Enabled {
		log.Printf("Setting up TLS with cert: %s, key: %s", cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		creds, err := credentials.NewServerTLSFromFile(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		if err != nil {
			log.Fatalf("Failed to setup TLS: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// Create and configure the registry
	log.Printf("Creating registry with max nodes: %d", cfg.Registry.MaxNodes)
	reg := registry.NewRegistry(cfg, pluginManager)

	// Start the registry (initializes plugins, starts health check loop)
	if err := reg.Start(ctx); err != nil {
		log.Fatalf("Failed to start registry: %v", err)
	}

	// Create gRPC server and register services
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRegistryServer(grpcServer, reg)

	// Start listening
	log.Printf("Starting Galaxy Node Pool server on %s (TLS: %v)", cfg.Server.Address, cfg.Server.TLS.Enabled)
	listener, err := registry.Listen(cfg.Server.Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down server...")
	
	// Graceful shutdown
	grpcServer.GracefulStop()
	log.Println("Server shutdown complete")
}
