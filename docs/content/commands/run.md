---
title: run
description: Make HTTP requests with automatic credential resolution
order: 1
---

# sreq run

Make HTTP requests with automatic credential resolution.

## Synopsis

```bash
sreq run <METHOD> <path> [flags]
```

## Description

The `run` command is the core of sreq. It resolves credentials from configured providers and makes authenticated HTTP requests.

## Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `METHOD` | HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS) | Yes |
| `path` | Request path (e.g., `/api/v1/users`) | Yes |

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--service` | `-s` | Service name | — |
| `--env` | `-e` | Environment | `dev` |
| `--data` | `-d` | Request body (or @filename) | — |
| `--header` | `-H` | Add header (repeatable) | — |
| `--output` | `-o` | Output format: `json`, `raw`, `headers` | `json` |
| `--timeout` | | Request timeout | `30s` |
| `--offline` | | Use cached credentials only | `false` |
| `--no-cache` | | Skip cache, fetch fresh credentials | `false` |
| `--verbose` | `-v` | Show detailed output | `false` |
| `--dry-run` | | Preview without executing | `false` |

## Examples

### Basic GET Request

```bash
sreq run GET /api/v1/users -s auth-service -e dev
```

### POST with JSON Body

```bash
sreq run POST /api/v1/users -s auth-service -e dev \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### POST with Body from File

```bash
sreq run POST /api/v1/users -s auth-service -d @payload.json
```

### Custom Headers

```bash
sreq run GET /api/v1/users -s auth-service \
  -H "X-Request-ID: abc123" \
  -H "Accept: application/json"
```

### Verbose Mode

See what sreq is doing under the hood:

```bash
sreq run GET /api/v1/users -s auth-service -v
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
  Password:  ap****rd

{"users": [...]}
```

### Dry Run

Preview what would be sent without executing:

```bash
sreq run POST /api/v1/users -s auth-service --dry-run
```

### Different Output Formats

```bash
# Pretty-printed JSON (default)
sreq run GET /api -s auth -o json

# Raw response body
sreq run GET /api -s auth -o raw

# Include response headers
sreq run GET /api -s auth -o headers
```

### Offline Mode

Use cached credentials without network access:

```bash
# First, cache credentials
sreq sync dev

# Then use offline
sreq run GET /api -s auth --offline
```

### Using Contexts

```bash
# Use a named context
sreq run GET /api -s auth -c production

# Override context values
sreq run GET /api -s auth -c production -e staging
```

## Request Flow

1. Load configuration from `~/.sreq/config.yaml`
2. Resolve context (from `-c` flag or `default_context`)
3. Check credential cache (unless `--no-cache`)
4. If not cached, fetch from providers (Consul, AWS, etc.)
5. Cache credentials for next time
6. Make HTTP request with resolved credentials
7. Save to history (unless `SREQ_NO_HISTORY=1`)
8. Display response

## See Also

- [history](/commands/history) — View and replay requests
- [cache](/commands/cache) — Manage credential cache
- [Configuration](/configuration) — Setup providers and services
