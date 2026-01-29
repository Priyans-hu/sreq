---
title: Consul Provider
description: HashiCorp Consul KV provider configuration
order: 11
---

# Consul Provider

The Consul provider fetches credentials from HashiCorp Consul's Key-Value store.

## Overview

Consul KV is commonly used to store:

- Service base URLs
- Usernames and non-sensitive configuration
- Feature flags and settings

## Configuration

### Basic Setup

```yaml
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_TOKEN}
    paths:
      base_url: "services/{service}/{env}/base_url"
      username: "services/{service}/{env}/username"
```

### Full Configuration

```yaml
providers:
  consul:
    # Consul server address
    address: consul.example.com:8500

    # ACL token for authentication
    token: ${CONSUL_TOKEN}

    # Optional datacenter
    datacenter: dc1

    # Optional TLS configuration
    scheme: https
    tls_skip_verify: false

    # Environment-specific addresses
    addresses:
      dev: consul-dev.example.com:8500
      staging: consul-staging.example.com:8500
      prod: consul-prod.example.com:8500

    # Path templates for credential resolution
    paths:
      base_url: "services/{service}/{env}/base_url"
      username: "services/{service}/{env}/username"
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `address` | Consul server address | `localhost:8500` |
| `token` | ACL token | — |
| `datacenter` | Datacenter name | — |
| `scheme` | HTTP or HTTPS | `http` |
| `tls_skip_verify` | Skip TLS verification | `false` |
| `addresses` | Environment-specific addresses | — |
| `paths` | Path templates | — |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CONSUL_HTTP_ADDR` | Consul address (fallback) |
| `CONSUL_TOKEN` | ACL token |
| `CONSUL_DATACENTER` | Datacenter |

## Path Templates

Templates support these placeholders:

| Placeholder | Description | Example Value |
|-------------|-------------|---------------|
| `{service}` | Service's `consul_key` or name | `auth` |
| `{env}` | Current environment | `dev` |

### Example Resolution

Configuration:

```yaml
providers:
  consul:
    paths:
      base_url: "services/{service}/{env}/base_url"

services:
  auth-service:
    consul_key: auth
```

Request:

```bash
sreq run GET /api -s auth-service -e dev
```

Consul path queried: `services/auth/dev/base_url`

## Environment-Specific Addresses

For organizations with separate Consul clusters per environment:

```yaml
providers:
  consul:
    # Default address (used if no environment match)
    address: consul.example.com:8500

    # Environment-specific addresses
    addresses:
      dev: consul-dev.internal:8500
      staging: consul-staging.internal:8500
      prod: consul-prod.internal:8500
```

When you run `sreq run GET /api -s auth -e prod`, sreq connects to `consul-prod.internal:8500`.

## Authentication

### ACL Token

Consul ACL tokens can be provided via:

1. **Config file** (not recommended for sensitive tokens):

   ```yaml
   providers:
     consul:
       token: my-token
   ```

2. **Environment variable reference** (recommended):

   ```yaml
   providers:
     consul:
       token: ${CONSUL_TOKEN}
   ```

3. **Environment variable directly**:

   ```bash
   export CONSUL_TOKEN=my-token
   ```

### Token Permissions

The token needs `read` permission on the KV paths:

```hcl
key_prefix "services/" {
  policy = "read"
}
```

## TLS Configuration

For HTTPS connections:

```yaml
providers:
  consul:
    address: consul.example.com:8501
    scheme: https
    tls_skip_verify: false  # Set true only for self-signed certs
```

## Testing Connection

Verify Consul connectivity:

```bash
sreq config test
```

Output:

```
Consul:
  Address: consul.example.com:8500
  Status:  ✓ Connected
  Leader:  10.0.1.5:8300
```

## Troubleshooting

### Connection Refused

```
Error: failed to connect to Consul: dial tcp: connection refused
```

**Solutions:**

- Verify Consul address is correct
- Check network connectivity
- Ensure Consul is running

### Permission Denied

```
Error: permission denied for key "services/auth/dev/base_url"
```

**Solutions:**

- Verify ACL token has read permissions
- Check token is not expired
- Verify key path exists

### Key Not Found

```
Error: key not found: services/auth/dev/base_url
```

**Solutions:**

- Verify the key exists in Consul
- Check path template is correct
- Verify service's `consul_key` matches

## See Also

- [AWS Provider](/providers/aws) — AWS Secrets Manager setup
- [Configuration](/configuration) — Full configuration reference
