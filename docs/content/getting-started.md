---
title: Getting Started
description: Make your first request with sreq in 5 minutes
order: 3
---

# Getting Started

This guide will have you making authenticated API requests in under 5 minutes.

## Prerequisites

- sreq installed ([Installation](/sreq/installation))
- Access to at least one provider (Consul, AWS, env vars, or a `.env` file)

## Step 1: Initialize Configuration

Create the default configuration:

```bash
sreq init
```

This creates `~/.sreq/config.yaml` with a starter template.

## Step 2: Configure Authentication

Run the interactive auth setup:

```bash
sreq auth
```

This guides you through configuring:

1. **Consul** — Address and token
2. **AWS** — Region and profile

Alternatively, configure providers individually:

```bash
sreq auth consul   # Just Consul
sreq auth aws      # Just AWS
```

> **Tip:** For quick local testing, you can skip `sreq auth` and use the `env` or `dotenv` provider instead — just set environment variables or create a `.env` file. See [Providers](/sreq/providers/) for details.

## Step 3: Add a Service

Add your first service:

```bash
sreq service add auth-service \
  --consul-key auth \
  --aws-prefix auth-svc
```

This tells sreq:

- Look for `auth` keys in Consul for base URL and username
- Look for `auth-svc` prefix in AWS Secrets Manager for passwords

## Step 4: Make a Request

Now make your first request:

```bash
sreq run GET /api/v1/health -s auth-service -e dev
```

**What happens:**

1. sreq loads config for `auth-service`
2. Fetches base URL from Consul: `services/auth/dev/base_url`
3. Fetches credentials from AWS: `auth-svc/dev/password`
4. Makes authenticated GET request to `<base_url>/api/v1/health`
5. Returns the response

## Step 5: Explore Features

### Verbose Mode

See what sreq is doing:

```bash
sreq run GET /api/v1/users -s auth-service -e dev --verbose
```

Output:

```
Request Details:
  Service:     auth-service
  Environment: dev
  Method:      GET
  Path:        /api/v1/users

Resolving credentials from providers...
  Base URL:  https://auth.example.com
  Username:  api_user
  Password:  ap****ey

{"users": [...]}
```

### POST with Body

```bash
sreq run POST /api/v1/users -s auth-service -e dev \
  -d '{"name": "John", "email": "john@example.com"}'
```

### POST with File

```bash
sreq run POST /api/v1/users -s auth-service -e dev -d @user.json
```

### Custom Headers

```bash
sreq run GET /api/v1/users -s auth-service -e dev \
  -H "X-Request-ID: 12345" \
  -H "Accept: application/json"
```

### Dry Run

See what would be sent without executing:

```bash
sreq run GET /api/v1/users -s auth-service -e dev --dry-run
```

## Quick Reference

| Command | Description |
|---------|-------------|
| `sreq init` | Initialize configuration |
| `sreq auth` | Configure provider authentication |
| `sreq service add <name>` | Add a new service |
| `sreq service list` | List configured services |
| `sreq run <METHOD> <path>` | Make an HTTP request |
| `sreq history` | View request history |
| `sreq sync <env>` | Cache credentials for offline use |
| `sreq tui` | Launch interactive TUI |

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--service` | `-s` | Service name |
| `--env` | `-e` | Environment (dev/staging/prod) |
| `--data` | `-d` | Request body (or @filename) |
| `--header` | `-H` | Add header (repeatable) |
| `--verbose` | `-v` | Show detailed output |
| `--dry-run` | | Preview without executing |
| `--output` | `-o` | Output format (json/raw/headers) |

## Next Steps

- [Configuration](/sreq/configuration) — Deep dive into config options
- [Commands](/sreq/commands) — Full command reference
- [Providers](/sreq/providers) — Provider-specific configuration
