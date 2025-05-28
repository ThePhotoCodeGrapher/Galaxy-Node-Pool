# Main Pool Eligibility Requirements

## Overview

This document outlines the staking and eligibility requirements for nodes and organizations to connect to the main Galaxy pool. These requirements ensure that only high-quality, reliable nodes are allowed to participate in the main pool, maintaining its performance and reputation.

## Staking Requirements

### Organization Staking

Organizations must stake a minimum amount of HCCO tokens to participate in the main pool. The staking amount is calculated based on:

1. **Node Capacity**: Higher capacity nodes require larger stakes
2. **Processing Power**: More powerful nodes require larger stakes
3. **GAS Generation Rate**: Nodes that generate more GAS require larger stakes

The staking formula is:

```
Required Stake = (Base Stake) × (Capacity Factor) × (Power Factor) × (GAS Generation Factor)
```

Where:
- **Base Stake**: Minimum stake required (currently 1000 HCCO)
- **Capacity Factor**: 1.0 + (0.1 × each TB of storage)
- **Power Factor**: 1.0 + (0.2 × each 100 compute units)
- **GAS Generation Factor**: 1.0 + (0.5 × daily GAS generation rate)

### Individual Node Staking

Individual nodes must stake a minimum of 100 HCCO tokens to join the main pool, plus additional tokens based on their advertised capabilities.

## Eligibility Requirements

### "SOLID" Status Requirements

To achieve "SOLID" status and be eligible to connect to the main pool, nodes or organizations must meet the following criteria:

1. **Experience**: Must have successfully completed at least 500 AI worker tasks
2. **Reliability**: Must maintain a 99.5% uptime over the last 30 days
3. **Quality**: Must maintain a 4.8/5.0 or higher quality rating
4. **Previous Pool Operation**: Must have operated their own pool for at least 14 days
5. **Verification**: Must pass identity verification through the Stellar federation

### Probationary Period

New nodes joining the main pool enter a 7-day probationary period during which:
- They receive fewer task assignments
- They are monitored more closely for performance issues
- Their stake is locked and cannot be withdrawn

### Maintaining Eligibility

To maintain eligibility in the main pool:
- Nodes must maintain a 99.5% uptime
- Quality ratings must remain above 4.5/5.0
- Stake must remain locked
- Monthly performance reviews must be passed

## Implementation

The eligibility verification process is implemented through:

1. **Stellar Smart Contract**: Verifies and locks the required stake
2. **Federation Protocol**: Verifies identity and previous pool operation
3. **Performance Monitoring**: Tracks uptime and quality metrics
4. **Automatic Evaluation**: Script that evaluates eligibility based on all factors

## Script Example

```javascript
// Example of the eligibility calculation script
function calculateRequiredStake(nodeCapacity, processingPower, gasGeneration) {
  const baseStake = 1000;
  const capacityFactor = 1.0 + (0.1 * (nodeCapacity / 1024)); // TB to GB
  const powerFactor = 1.0 + (0.2 * (processingPower / 100));
  const gasFactor = 1.0 + (0.5 * gasGeneration);
  
  return baseStake * capacityFactor * powerFactor * gasFactor;
}

function checkSolidStatus(completedTasks, uptime, qualityRating, poolOperationDays, verificationStatus) {
  if (completedTasks < 500) return false;
  if (uptime < 99.5) return false;
  if (qualityRating < 4.8) return false;
  if (poolOperationDays < 14) return false;
  if (verificationStatus !== 'verified') return false;
  
  return true;
}
```

## Benefits of Main Pool Participation

Organizations and nodes that meet these requirements and join the main pool receive:

1. **Higher Task Priority**: Preferred assignment of high-value tasks
2. **Enhanced Reputation**: "SOLID" status badge visible to clients
3. **Reduced Fees**: Lower transaction fees on the network
4. **Governance Rights**: Voting rights on network governance decisions
5. **Reward Multiplier**: Increased rewards for completed tasks

## Revocation of Eligibility

Eligibility can be revoked for:
- Sustained poor performance
- Malicious behavior
- Unstaking tokens
- Violating network policies

When eligibility is revoked, nodes are removed from the main pool and must re-qualify after a 30-day cooling-off period.
