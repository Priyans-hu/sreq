---
title: auth
description: Configure authentication for providers
order: 3
---

# sreq auth

Interactively configure authentication for providers.

## Synopsis

```bash
sreq auth [provider] [flags]
```

## Description

The `auth` command provides an interactive setup wizard for configuring provider credentials. It guides you through setting up Consul and AWS authentication.

## Subcommands

| Command | Description |
|---------|-------------|
| `sreq auth` | Interactive setup for all providers |
| `sreq auth consul` | Configure Consul only |
| `sreq auth aws` | Configure AWS only |

## Examples

### Full Interactive Setup

```bash
sreq auth
```

Output:

```
sreq Authentication Setup
==========================

1. Consul Configuration
------------------------
Consul address [consul.example.com:8500]: consul.mycompany.com:8500
Consul token: <hidden>
Test connection? [Y/n]: y
✓ Connected to Consul successfully

2. AWS Configuration
---------------------
AWS region [us-east-1]: us-west-2
AWS profile [default]: myprofile
Test connection? [Y/n]: y
✓ AWS credentials valid

Authentication setup complete!
Run 'sreq config test' to verify your configuration.
```

### Configure Consul Only

```bash
sreq auth consul
```

### Configure AWS Only

```bash
sreq auth aws
```

## Provider Configuration

### Consul

The wizard prompts for:

| Field | Description | Environment Variable |
|-------|-------------|---------------------|
| Address | Consul server address | `CONSUL_HTTP_ADDR` |
| Token | ACL token for authentication | `CONSUL_TOKEN` |
| Datacenter | Optional datacenter | — |

### AWS

The wizard prompts for:

| Field | Description | Environment Variable |
|-------|-------------|---------------------|
| Region | AWS region | `AWS_REGION` |
| Profile | AWS credentials profile | `AWS_PROFILE` |

## Configuration Output

The wizard updates `~/.sreq/config.yaml`:

```yaml
providers:
  consul:
    address: consul.mycompany.com:8500
    token: ${CONSUL_TOKEN}  # References environment variable
    datacenter: dc1

  aws_secrets:
    region: us-west-2
    profile: myprofile
```

## Security Notes

- Tokens can be stored directly or reference environment variables with `${VAR}`
- Environment variable references are recommended for sensitive values
- Config file permissions are set to `600` (owner read/write only)

## See Also

- [Configuration](/sreq/configuration) — Full configuration reference
- [Consul Provider](/sreq/providers/consul) — Consul setup details
- [AWS Provider](/sreq/providers/aws) — AWS setup details
