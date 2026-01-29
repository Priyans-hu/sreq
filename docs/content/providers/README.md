---
title: Providers
description: Credential providers supported by sreq
order: 10
---

# Providers

Providers are the sources from which sreq fetches credentials. Each provider connects to a different secret management system.

## Supported Providers

| Provider | Status | Description |
|----------|--------|-------------|
| [Consul KV](/sreq/providers/consul) | Available | HashiCorp Consul Key-Value store |
| [AWS Secrets Manager](/sreq/providers/aws) | Available | AWS Secrets Manager |
| [Environment Variables](/sreq/providers/env) | Available | System environment variables |
| [dotenv](/sreq/providers/dotenv) | Available | Local `.env` files |
| HashiCorp Vault | Planned | HashiCorp Vault KV secrets |
| GCP Secret Manager | Planned | Google Cloud Secret Manager |
| Azure Key Vault | Planned | Azure Key Vault |

## How Providers Work

When you make a request, sreq:

1. Reads the service configuration
2. Determines which providers to query
3. Fetches credentials from each provider
4. Combines them into a complete credential set
5. Makes the authenticated request

## Provider Priority

When multiple providers can supply the same credential, sreq uses the first successful result in this order:

1. Explicit path mappings in service config
2. Provider path templates
3. Default values

## Credential Types

Providers can supply these credential types:

| Credential | Description | Common Sources |
|------------|-------------|----------------|
| `base_url` | Service base URL | Consul |
| `username` | Authentication username | Consul |
| `password` | Authentication password | AWS, Vault |
| `api_key` | API key | AWS, Vault |
| `token` | Bearer token | AWS, Vault |

## Configuration

Providers are configured in `~/.sreq/config.yaml`:

```yaml
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_TOKEN}
    paths:
      base_url: "services/{service}/{env}/base_url"
      username: "services/{service}/{env}/username"

  aws_secrets:
    region: us-east-1
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"
```

## Path Templates

Path templates use placeholders replaced at runtime:

| Placeholder | Description |
|-------------|-------------|
| `{service}` | Service name or consul_key/aws_prefix |
| `{env}` | Current environment |
| `{region}` | AWS region |
| `{project}` | Project name |
| `{app}` | Application name |

## Adding a Provider

To add a new provider to sreq, implement the `Provider` interface:

```go
type Provider interface {
    Name() string
    Fetch(ctx context.Context, path string) (string, error)
}
```

See the [Contributing Guide](https://github.com/Priyans-hu/sreq/blob/main/CONTRIBUTING.md) for details.

## Next Steps

- [Consul Provider](/sreq/providers/consul) — Setup and configuration
- [AWS Provider](/sreq/providers/aws) — Setup and configuration
- [Environment Variables](/sreq/providers/env) — Read from system env vars
- [dotenv Provider](/sreq/providers/dotenv) — Read from .env files
