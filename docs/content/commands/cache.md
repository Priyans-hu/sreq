---
title: cache & sync
description: Manage credential cache for offline use
order: 6
---

# sreq cache & sync

Manage the local credential cache for faster requests and offline use.

## Overview

sreq caches credentials locally after fetching them from providers. This enables:

- **Faster requests** — Skip provider calls on subsequent requests
- **Offline mode** — Work without network access to providers
- **Reduced load** — Fewer calls to Consul/AWS

## Commands

| Command | Description |
|---------|-------------|
| `sreq cache status` | Show cache status and entries |
| `sreq cache clear [env]` | Clear cached credentials |
| `sreq sync <env>` | Sync credentials to cache |

## sreq cache status

Show current cache status:

```bash
sreq cache status
```

Output:

```
Cache Status
============
Enabled:    true
Directory:  /Users/you/.sreq/cache
TTL:        1h0m0s
Entries:    3
Total Size: 2048 bytes

Cached Entries:
  auth-service/dev - cached 10:30:00, expires 11:30:00
  auth-service/prod - cached 10:31:00, expires 11:31:00
  billing-service/dev - cached 10:32:00, expires 11:32:00 (EXPIRED)
```

## sreq cache clear

Clear cached credentials:

```bash
# Clear all cached credentials
sreq cache clear

# Clear specific environment only
sreq cache clear dev
```

## sreq sync {#sync}

Pre-fetch and cache credentials for all services in an environment:

```bash
sreq sync dev
```

Output:

```
Syncing dev environment...
  auth-service: synced
  billing-service: synced
  user-service: synced

Synced 3 credentials successfully
```

### Sync Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--all` | Sync all environments | `false` |
| `--force` | Force refresh even if cache is valid | `false` |

### Sync All Environments

```bash
sreq sync --all
```

### Force Refresh

```bash
sreq sync dev --force
```

## Offline Mode

After syncing, use `--offline` to make requests without provider access:

```bash
# Pre-cache credentials
sreq sync dev

# Later, work offline
sreq run GET /api/v1/users -s auth-service --offline
```

If credentials aren't cached, offline mode will fail with a helpful message:

```
Error: no cached credentials found for auth-service/dev (run 'sreq sync dev' first)
```

## Cache Security

Cached credentials are protected with:

- **AES-256-GCM encryption** — Industry-standard encryption
- **User-specific key** — Generated on first use, stored in `~/.sreq/.cache_key`
- **File permissions** — `600` (owner read/write only)
- **Configurable TTL** — Default 1 hour

## Disabling Cache

### Per-Request

```bash
sreq run GET /api -s auth --no-cache
```

### Globally

```bash
export SREQ_NO_CACHE=1
```

### In CI/CD

Cache is automatically disabled when `CI=true` is set (common in GitHub Actions, GitLab CI, etc.).

## Cache Configuration

Configure cache behavior in `config.yaml`:

```yaml
cache:
  ttl: 2h          # Cache TTL (default: 1h)
  dir: ~/.sreq/cache  # Cache directory
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SREQ_NO_CACHE` | Disable caching | `0` |
| `CI` | Auto-disable cache in CI | — |

## See Also

- [run](/commands/run) — Make requests with `--offline` flag
- [Configuration](/configuration) — Cache configuration options
