---
title: init
description: Initialize sreq configuration
order: 2
---

# sreq init

Initialize sreq configuration with a starter template.

## Synopsis

```bash
sreq init [flags]
```

## Description

Creates the configuration directory (`~/.sreq/`) and generates a starter `config.yaml` file with example provider configurations.

If configuration already exists, `init` will not overwrite it unless `--force` is used.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Overwrite existing configuration | `false` |

## Examples

### Initialize Configuration

```bash
sreq init
```

Output:

```
Created configuration directory: ~/.sreq/
Created config file: ~/.sreq/config.yaml

Next steps:
  1. Run 'sreq auth' to configure provider credentials
  2. Run 'sreq service add <name>' to add a service
  3. Run 'sreq run GET /api -s <service>' to make requests
```

### Force Reinitialize

```bash
sreq init --force
```

## Generated Configuration

The generated `config.yaml` includes:

```yaml
# sreq configuration
# Documentation: https://github.com/Priyans-hu/sreq

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

environments:
  - dev
  - staging
  - prod

default_env: dev

services: {}
```

## What Gets Created

| Path | Description |
|------|-------------|
| `~/.sreq/` | Configuration directory |
| `~/.sreq/config.yaml` | Main configuration file |
| `~/.sreq/cache/` | Credential cache directory |
| `~/.sreq/history.json` | Request history file |

## See Also

- [auth](/commands/auth) — Configure provider authentication
- [Configuration](/configuration) — Full configuration reference
