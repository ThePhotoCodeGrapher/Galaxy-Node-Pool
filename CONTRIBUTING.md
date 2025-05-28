# Contributing to Galaxy Node Pool

Thank you for your interest in contributing to the Galaxy Node Pool project! This document outlines the process for contributing to our project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Security](#security)
- [Enterprise Features](#enterprise-features)
- [Mainnet Integration](#mainnet-integration)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.20+
- Docker 20.10+
- Git
- Make

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/galaxy-node-pool.git
   cd galaxy-node-pool
   ```
3. Set up the pre-commit hooks:
   ```bash
   make setup
   ```
4. Build the project:
   ```bash
   make build
   ```

## Development Workflow

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Make your changes
3. Run tests and linters:
   ```bash
   make test
   make lint
   ```
4. Commit your changes with a descriptive message
5. Push to your fork and open a pull request

## Pull Request Process

1. Ensure all tests are passing
2. Update the documentation as needed
3. Include unit tests for new features
4. Ensure your code follows the project's coding standards
5. Submit the pull request with a clear description of the changes

## Coding Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` before committing
- Keep functions small and focused
- Write clear and concise commit messages
- Document all exported functions and types

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests (requires Docker)
make test-integration

# Run e2e tests
make test-e2e
```

### Writing Tests

- Write tests for all new features and bug fixes
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for at least 80% test coverage

## Documentation

### Updating Documentation

1. Update the relevant `.md` files in the `docs/` directory
2. Ensure all new features are documented
3. Update examples if they are affected by your changes

### Generating Documentation

```bash
# Generate API documentation
make docs
```

## Security

### Reporting Security Issues

Please report security issues to security@hybridconnect.cloud. Do not file a public issue.

### Security Best Practices

- Never commit secrets or credentials
- Validate all inputs
- Use prepared statements for database queries
- Follow the principle of least privilege

## Enterprise Features

Enterprise features are developed in the `internal/enterprise/` directory. To contribute to enterprise features:

1. Sign the CLA (Contributor License Agreement)
2. Request access to the enterprise repository
3. Follow the same contribution guidelines

## Mainnet Integration

Mainnet integration requires additional steps:

1. Set up a Stellar testnet account
2. Configure the staking contract address
3. Ensure you have sufficient test tokens

## Community

- [GitHub Discussions](https://github.com/castle-palette/galaxy-node-pool/discussions)
- [Community Forum](https://community.hybridconnect.cloud)
- [Discord](https://discord.gg/hybridconnect)

---

**AI-ID**: CP-GAL-NODEPOOL-001  
**Documentation**: [https://docs.hybridconnect.cloud/galaxy-node-pool/contributing](https://docs.hybridconnect.cloud/galaxy-node-pool/contributing)  
**Maintained by**: [Castle Palette Cloud A.I.](https://hybridconnect.cloud)
