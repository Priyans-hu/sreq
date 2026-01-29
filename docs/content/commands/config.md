---
title: config
description: View and validate configuration
order: 8
---

# sreq config

View and validate sreq configuration.

## Synopsis

```bash
sreq config <command>
```

## Subcommands

| Command | Description |
|---------|-------------|
| `sreq config show` | Display current configuration |
| `sreq config path` | Show config file path |
| `sreq config test` | Test provider connections |

## sreq config show

Display the current configuration:

```bash
sreq config show
```

Output:

```yaml
# Current Configuration

providers:
  consul:
    address: consul.example.com:8500
    token: <set via CONSUL_TOKEN>

  aws_secrets:
    region: us-east-1
    profile: default

environments:
  - dev
  - staging
  - prod

default_env: dev

services:
  auth-service:
    consul_key: auth
    aws_prefix: auth-svc
```

## sreq config path

Show the config file path:

```bash
sreq config path
```

Output:

```
/Users/you/.sreq/config.yaml
```

## sreq config test

Test connections to configured providers:

```bash
sreq config test
```

Output:

```
Testing provider connections...

Consul:
  Address: consul.example.com:8500
  Status:  ✓ Connected
  Leader:  10.0.1.5:8300

AWS Secrets Manager:
  Region:  us-east-1
  Profile: default
  Status:  ✓ Credentials valid

All providers connected successfully!
```

### Test Failures

If a provider fails:

```
Testing provider connections...

Consul:
  Address: consul.example.com:8500
  Status:  ✗ Connection failed
  Error:   dial tcp: connection refused

AWS Secrets Manager:
  Region:  us-east-1
  Status:  ✓ Credentials valid

1 of 2 providers failed. Run 'sreq auth' to reconfigure.
```

## Configuration Location

| File | Description |
|------|-------------|
| `~/.sreq/config.yaml` | Main configuration |
| `~/.sreq/services.yaml` | Service definitions (optional) |

Override with environment variable:

```bash
export SREQ_CONFIG=/path/to/custom/config.yaml
```

## See Also

- [init](/commands/init) — Initialize configuration
- [auth](/commands/auth) — Configure authentication
- [Configuration](/configuration) — Full reference
