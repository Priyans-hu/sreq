---
title: env
description: Manage environments
order: 9
---

# sreq env

Manage environment configurations.

## Synopsis

```bash
sreq env <command>
```

## Subcommands

| Command | Description |
|---------|-------------|
| `sreq env list` | List available environments |
| `sreq env switch <env>` | Switch default environment |
| `sreq env current` | Show current default environment |

## sreq env list

List all configured environments:

```bash
sreq env list
```

Output:

```
Available environments:

  dev      (default)
  staging
  prod
```

## sreq env switch

Change the default environment:

```bash
sreq env switch staging
```

Output:

```
Default environment changed: dev → staging
```

Now requests without `-e` flag will use staging:

```bash
# Uses staging environment
sreq run GET /api -s auth-service
```

## Environment Configuration

Environments are defined in `config.yaml`:

```yaml
environments:
  - dev
  - staging
  - prod
  - qa        # Add custom environments

default_env: dev
```

## sreq env current

Show the current default environment:

```bash
sreq env current
```

Output:

```
Current default environment: dev
```

## Per-Request Override

Always override with `-e` flag:

```bash
sreq run GET /api -s auth -e prod
```

## Environment-Specific Provider Addresses

Configure different provider addresses per environment:

```yaml
providers:
  consul:
    # Default address
    address: consul.example.com:8500

    # Environment-specific addresses
    addresses:
      dev: consul-dev.example.com:8500
      staging: consul-staging.example.com:8500
      prod: consul-prod.example.com:8500
```

## See Also

- [run](/sreq/commands/run) — Make requests with `-e` flag
- [Configuration](/sreq/configuration) — Full configuration reference
