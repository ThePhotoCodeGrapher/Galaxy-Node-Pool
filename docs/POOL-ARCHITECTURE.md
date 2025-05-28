# Galaxy Node Pool Architecture

## Core Architecture

### What is a Pool?
A pool is a registry and discovery service for Galaxy Network Nodes. Pools can be public (main/test/dev) or private (enterprise/org). Nodes register with pools, advertise their specialization, and can be discovered by clients (VS Code extension, job submitters, AI engineers).

### Modular Design
The Galaxy Node Pool is built with a modular architecture that separates core functionality from extensions:

1. **Core Components**:
   - Node registration and management
   - Pool configuration and operation
   - Basic discovery mechanisms
   - Command-line interface

2. **Plugin Architecture**:
   - Extensible plugin system for additional features
   - Clean separation between open-source and enterprise features
   - Dynamic loading of plugins at runtime
   - Standardized plugin interface

## Pool Types

### Public vs Private Pools
- **Public Pools**: Anyone can register nodes and submit jobs (e.g., mainnet, testnet).
- **Private Pools**: Only authorized nodes/users can register or submit jobs (e.g., enterprise/org pools).

### Node Specialization
Nodes advertise their specialization (developer, storage, analytics, etc.) when registering with a pool.

## Operational Flow

### Discovery Flow
1. Node registers with the pool (gRPC API).
2. Pool tracks node health, specialization, and metadata.
3. Clients query the pool for available nodes.
4. Jobs are routed to nodes based on specialization, availability, and pool privacy.

### CLI Management
1. Administrators use the CLI to manage pools and nodes
2. Commands are organized by functionality
3. Enterprise features are loaded via plugins

## Security
- All communication via TLS.
- Authentication for private pools (API keys, mTLS, etc.).
- Role-based access control for job visibility and execution.

## Gelato Protocol Pattern
- Modular, Dockerized, config-driven.
- Community/enterprise split for features and access.
- Extensible for future node/pool types.
