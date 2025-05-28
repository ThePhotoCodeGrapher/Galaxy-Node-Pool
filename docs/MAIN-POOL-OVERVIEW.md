# Galaxy Node Pool - Main Pool Overview

> **AI-ID**: CP-GAL-NODEPOOL-001  
> **Organization**: Castle Palette Cloud A.I.  
> **Website**: [https://hybridconnect.cloud/](https://hybridconnect.cloud/)  
> **Protocol**: .gal  
> **Version**: 1.0.0  
> **Last Updated**: 2025-05-28

## Purpose
This document defines the architecture, eligibility requirements, and operational parameters of the main Galaxy pool, which is separate from the open source components of the Galaxy Node Pool project.

## Main Pool Architecture

### Overview
The main Galaxy pool is a specialized implementation of the Galaxy Node Pool architecture that serves as the primary entry point for the Galaxy Network. It implements additional security, scaling, and eligibility verification mechanisms not present in the open source components.

### Components
- **Main Pool Server**: High-availability implementation of the Galaxy Node Pool
- **Eligibility Verification**: System for verifying node and organization eligibility
- **Staking Contract**: Integration with Stellar for staking requirements
- **Performance Monitoring**: Advanced monitoring and alerting system
- **Admin API**: Administrative interface for managing the main pool

### Separation from Open Source
The main pool implementation is intentionally separated from the open source components to:
1. Protect proprietary business logic
2. Ensure security of the main network
3. Allow for specialized optimizations
4. Implement advanced features not suitable for general release

## Eligibility Requirements

### Staking Requirements
Organizations and nodes must stake HCCO tokens to participate in the main pool:

- **Organizations**: Base stake of 1000 HCCO, adjusted by:
  - Node capacity (storage)
  - Processing power
  - GAS generation rate

- **Individual Nodes**: Base stake of 100 HCCO, adjusted by:
  - Node capacity
  - Processing power
  - Specialization type

### "SOLID" Status Requirements
To achieve "SOLID" status and be eligible to connect to the main pool:

1. **Experience**: Must have completed at least 500 AI worker tasks
2. **Reliability**: Must maintain 99.5% uptime over 30 days
3. **Quality**: Must maintain 4.8/5.0 or higher quality rating
4. **Previous Pool Operation**: Must have operated own pool for at least 14 days
5. **Verification**: Must pass identity verification through Stellar federation

### Probationary Period
New nodes joining the main pool enter a 7-day probationary period with:
- Reduced task assignments
- Enhanced monitoring
- Locked stake

## Implementation Details

### Eligibility Verification
```yaml
eligibility:
  solid_status:
    min_completed_tasks: 500
    min_uptime_percentage: 99.5
    min_quality_rating: 4.8
    min_pool_operation_days: 14
    verification_required: true
  staking:
    organization_base: 1000
    node_base: 100
    capacity_factor: 0.1
    power_factor: 0.2
    gas_factor: 0.5
  probation:
    duration_days: 7
    task_reduction_percentage: 50
```

### Scaling Configuration
```yaml
scaling:
  min_nodes: 10
  max_nodes: 1000
  scale_up_threshold: 0.8
  scale_down_threshold: 0.2
  scale_up_increment: 5
  scale_down_increment: 2
  cooldown_period: "5m"
```

## Value Proposition
- **For Nodes**: Access to high-value tasks and enhanced reputation
- **For Organizations**: Participation in the core Galaxy Network
- **For the Network**: High-quality, reliable node pool with strict requirements
- **For Users**: Consistent, high-performance AI services

## Integration with Open Source Components
The main pool uses the same core architecture as the open source Galaxy Node Pool, but with additional components and stricter requirements. Organizations can:

1. Start with the open source components to build their own pool
2. Gain experience and establish a track record
3. Meet the eligibility requirements
4. Connect to the main pool once "SOLID" status is achieved

## Relationship to Main Net
The main pool is a participant in the main net, not the main net itself:
- The main net is the federation layer connecting all pools
- The main pool is a specialized implementation of a Galaxy Node Pool
- All pools, including the main pool, connect to the main net

---

*This document is part of the Galaxy Node Pool project by Castle Palette Cloud A.I.*  
*For more information, visit [https://hybridconnect.cloud/](https://hybridconnect.cloud/)*
