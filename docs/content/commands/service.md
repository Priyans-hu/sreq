---
title: service
description: Manage service configurations
order: 4
---

# sreq service

Manage service configurations.

## Synopsis

```bash
sreq service <command> [flags]
```

## Subcommands

| Command | Description |
|---------|-------------|
| `sreq service list` | List all configured services |
| `sreq service add <name>` | Add a new service |
| `sreq service remove <name>` | Remove a service |

## sreq service list

List all configured services.

```bash
sreq service list
```

Output:

```
Configured services:

  auth-service
    consul_key: auth
    aws_prefix: auth-svc

  billing-service
    consul_key: billing
    aws_prefix: billing
```

## sreq service add

Add a new service configuration.

### Simple Mode

Use path templates defined in `config.yaml`:

```bash
sreq service add auth-service --consul-key auth --aws-prefix auth-svc
```

| Flag | Description |
|------|-------------|
| `--consul-key` | Key prefix for Consul paths |
| `--aws-prefix` | Prefix for AWS Secrets Manager paths |

This creates:

```yaml
services:
  auth-service:
    consul_key: auth
    aws_prefix: auth-svc
```

With path templates like `services/{service}/{env}/base_url`, sreq will look for:
- Consul: `services/auth/dev/base_url`
- AWS: `auth-svc/dev/credentials#password`

### Advanced Mode

Specify explicit paths for each credential:

```bash
sreq service add invoice-api \
  --path base_url=consul:billing/invoice/url \
  --path username=consul:billing/invoice/user \
  --path password=aws:billing-secrets/invoice#pass
```

| Flag | Description |
|------|-------------|
| `--path` | Explicit path mapping as `key=provider:path` (repeatable) |

This creates:

```yaml
services:
  invoice-api:
    paths:
      base_url: "consul:billing/invoice/url"
      username: "consul:billing/invoice/user"
      password: "aws:billing-secrets/invoice#pass"
```

### Path Format

For advanced mode, paths use the format:

```
provider:path[#json_key]
```

| Part | Description | Example |
|------|-------------|---------|
| `provider` | Provider name | `consul`, `aws` |
| `path` | Path in the provider | `services/auth/url` |
| `#json_key` | Optional JSON key extraction | `#password` |

Examples:

```
consul:services/auth/base_url           # Consul KV value
aws:auth/dev/creds                      # AWS secret (full value)
aws:auth/dev/creds#password             # AWS secret JSON key
```

## sreq service remove

Remove a service configuration.

```bash
sreq service remove auth-service
```

Output:

```
Removed service: auth-service
```

## Configuration Location

Services are stored in `~/.sreq/config.yaml` under the `services` key:

```yaml
services:
  auth-service:
    consul_key: auth
    aws_prefix: auth-svc

  billing-service:
    consul_key: billing
    aws_prefix: billing
```

Alternatively, services can be defined in a separate `~/.sreq/services.yaml` file.

## See Also

- [run](/commands/run) — Make requests to services
- [Configuration](/configuration) — Full configuration reference
