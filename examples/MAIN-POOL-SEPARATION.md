# Main Pool Separation

## Overview

This document explains the separation between the open source Galaxy Node Pool components and the main pool implementation. The main pool implementation is not included in the open source repository for security, operational, and business reasons.

## Architecture Separation

The Galaxy Node Pool codebase is divided into three distinct parts:

1. **Open Source Components**
   - Core registry functionality
   - Plugin system
   - Module framework
   - Configuration system
   - Event system
   - Service container
   - Documentation (Gelato Protocol)

2. **Enterprise Components** (excluded from open source)
   - Stellar blockchain integration
   - Federation protocol
   - Staker rewards
   - Cross-pool discovery

3. **Main Pool Implementation** (excluded from open source)
   - Actual implementation of the main Galaxy pool
   - Specific configurations
   - Deployment scripts
   - Monitoring and alerting
   - Scaling logic
   - Admin tools

## Directory Structure

The main pool implementation is kept in separate directories that are excluded from the git repository:

```
/internal/main-pool/     # Main pool implementation code
/configs/main-pool/      # Main pool configuration files
/cmd/main-pool-server/   # Main pool server entry point
/docs/main-pool/         # Main pool documentation
```

## How to Use the Open Source Components

The open source components provide everything needed to run your own Galaxy Node Pool, but not the specific implementation of the main Galaxy pool. This allows:

1. **Transparency**: The core protocol and architecture are open source
2. **Extensibility**: Anyone can build their own pool using the framework
3. **Security**: The main pool's specific implementation details remain private
4. **Business Logic**: Proprietary features can be added to the main pool

## Relationship to Main Net

- All pools (including your own) connect to the main net
- The main pool is just one of many pools in the network
- The main pool has no special privileges in the protocol
- The main net registry is separate from any specific pool

## Creating Your Own Pool

To create your own pool using the open source components:

1. Clone the repository
2. Configure your pool settings
3. Implement any custom plugins you need
4. Deploy your pool
5. Register with the main net

Your pool will be a full peer in the Galaxy Network, with the same capabilities as the main pool.

## Main Pool vs. Main Net

It's important to understand the distinction:

- **Main Pool**: A specific implementation of a Galaxy Node Pool (not open source)
- **Main Net**: The network registry that all pools connect to (enterprise component)

The main net is what provides the federation capabilities, while the main pool is just one participant in that federation.
