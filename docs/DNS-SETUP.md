# DNS and Networking Guide

## Recommended Records
| Type | Name                                           | Value (IP/Target)            | Purpose                    |
|------|------------------------------------------------|------------------------------|----------------------------|
| A    | main.galaxy.network.castlepalette.com          | [main pool IP]               | Main public pool           |
| A    | test.galaxy.network.castlepalette.com          | [test pool IP]               | Testnet pool               |
| A    | dev.galaxy.network.castlepalette.com           | [dev pool IP]                | Dev pool                   |
| A    | acme.pool.galaxy.network.castlepalette.com     | [acme pool IP]               | Acme's private pool        |
| A    | node.1.main.galaxy.network.castlepalette.com   | [node1 IP]                   | Specific node (optional)   |
| CNAME| *.main.galaxy.network.castlepalette.com        | main.galaxy.network.castlepalette.com | Wildcard for nodes |
| CNAME| *.pool.galaxy.network.castlepalette.com        | [pool root or wildcard IP]   | Wildcard for org pools     |
| A    | *.pool.galaxy.net.*.asia.hybridconnect.cloud   | [test pool IP]               | Testnet wildcard domain    |

## Wildcard Usage
- Use wildcards for scalable node and pool addition.
- Example: `*.main.galaxy.network.castlepalette.com` covers all mainnet nodes.
- For testnet environments, use the standardized wildcard domain pattern: `*.{pool_env}.pool.galaxy.net.{org_id}.asia.hybridconnect.cloud`

## Security
- Always use TLS/SSL for all endpoints.
- Restrict private pools via firewall or authentication.
