---
title: Configuration
description: Complete guide to configuring sreq
order: 4
---

# Configuration

sreq uses YAML configuration files stored in `~/.sreq/`.

## Configuration Files

| File | Purpose |
|------|---------|
| `~/.sreq/config.yaml` | Main configuration (providers, defaults) |
| `~/.sreq/services.yaml` | Service definitions (optional, can be in config.yaml) |

## Full Configuration Example

```yaml
# ~/.sreq/config.yaml

# Provider configurations
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_TOKEN}
    datacenter: dc1
    # Environment-specific addresses (optional)
    addresses:
      dev: consul-dev.example.com:8500
      staging: consul-staging.example.com:8500
      prod: consul-prod.example.com:8500
    # Path templates for credential resolution
    paths:
      base_url: "services/{service}/{env}/base_url"
      username: "services/{service}/{env}/username"

  aws_secrets:
    region: us-east-1
    profile: default
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"

  env:
    paths:
      base_url: "{SERVICE}_BASE_URL"
      api_key: "{SERVICE}_API_KEY"

  dotenv:
    file: ".env"
    paths:
      base_url: "{SERVICE}_{ENV}_URL"
      password: "{SERVICE}_PASSWORD"

# Available environments
environments:
  - dev
  - staging
  - prod

# Default environment when -e is not specified
default_env: dev

# Contexts (presets for common configurations)
contexts:
  work:
    env: staging
    region: us-west-2
    project: backend
  personal:
    env: dev
    region: us-east-1

# Default context to use
default_context: work

# Service configurations
services:
  auth-service:
    consul_key: auth
    aws_prefix: auth-svc

  billing-service:
    consul_key: billing
    aws_prefix: billing

  # Advanced: explicit path mappings
  legacy-api:
    paths:
      base_url: "consul:legacy/api_url"
      username: "consul:legacy/user"
      password: "aws:legacy-creds/{env}#pass"
```

## Path Templates

Path templates use placeholders that are replaced at runtime:

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `{service}` | Service name | `auth-service` |
| `{env}` | Environment | `dev`, `staging`, `prod` |
| `{region}` | AWS region | `us-east-1` |
| `{project}` | Project name | `backend` |
| `{app}` | App name | `api` |

### Template Examples

```yaml
paths:
  # Simple template
  base_url: "services/{service}/base_url"
  # → services/auth-service/base_url

  # With environment
  password: "{service}/{env}/password"
  # → auth-service/dev/password

  # AWS JSON key extraction (use #)
  api_key: "{service}/{env}/credentials#api_key"
  # → Fetches auth-service/dev/credentials, extracts .api_key from JSON
```

## Service Configuration

### Simple Mode

Use `consul_key` and `aws_prefix` with path templates:

```yaml
services:
  auth-service:
    consul_key: auth        # Uses paths.base_url template with {service}=auth
    aws_prefix: auth-svc    # Uses paths.password template with {service}=auth-svc
```

### Advanced Mode

Specify exact paths per credential:

```yaml
services:
  custom-service:
    paths:
      base_url: "consul:custom/service/url"
      username: "consul:custom/service/user"
      password: "aws:custom-secrets/main#password"
      api_key: "aws:custom-secrets/main#api_key"
```

Path format: `provider:path` or `provider:path#json_key`

## Contexts

Contexts are presets for common flag combinations:

```yaml
contexts:
  production:
    env: prod
    region: us-east-1
    project: main

  development:
    env: dev
    region: us-west-2
    project: sandbox

default_context: development
```

Use contexts:

```bash
# Uses default_context settings
sreq run GET /api -s auth

# Override with specific context
sreq run GET /api -s auth -c production

# Override context values with flags
sreq run GET /api -s auth -c production -e staging
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SREQ_CONFIG` | Config file path | `~/.sreq/config.yaml` |
| `CONSUL_HTTP_TOKEN` | Consul ACL token | — |
| `CONSUL_HTTP_ADDR` | Consul address | — |
| `AWS_PROFILE` | AWS profile | `default` |
| `AWS_REGION` | AWS region | — |
| `SREQ_NO_CACHE` | Disable caching | `0` |
| `SREQ_NO_HISTORY` | Disable history | `0` |

## Environment Variable Substitution

Use `${VAR}` in config to reference environment variables:

```yaml
providers:
  consul:
    token: ${CONSUL_TOKEN}
    address: ${CONSUL_HTTP_ADDR}

  aws_secrets:
    region: ${AWS_REGION}
    profile: ${AWS_PROFILE}
```

## Validation

Test your configuration:

```bash
sreq config show     # Display current config
sreq config test     # Validate provider connections
```

## Configuration Precedence

From lowest to highest priority:

1. Default values
2. `config.yaml` settings
3. Context values (from `-c` or `default_context`)
4. Environment variables
5. Command-line flags

## Next Steps

- [Commands](/commands) — Full command reference
- [Providers](/providers) — Provider-specific setup
