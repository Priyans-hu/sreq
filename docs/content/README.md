---
title: sreq
description: Service-aware API client with automatic credential resolution
order: 1
---

# sreq

**Service-aware API client with automatic credential resolution.**

sreq eliminates the overhead of manually fetching credentials from multiple sources when testing APIs. Just specify the service name and environment — sreq handles the rest.

## The Problem

You want to test an API endpoint on your auth service:

```
POST /api/v1/users on auth-service
```

**Current workflow:**

1. Open Consul → find auth-service base URL for dev
2. Open AWS Secrets Manager → find auth-service credentials
3. Copy-paste into Postman/curl
4. Repeat for staging... repeat for prod...

This context-switching kills productivity and introduces errors.

## The Solution

```bash
sreq run POST /api/v1/users -s auth-service -e dev -d '{"name":"test"}'
```

**sreq automatically:**

- Fetches base URL from Consul
- Fetches credentials from AWS Secrets Manager
- Makes the authenticated request
- Caches credentials for faster subsequent requests

## Key Features

### Service-Aware

Pass the service name, sreq resolves everything else. No more hunting for URLs and credentials.

```bash
sreq run GET /api/v1/health -s billing-service -e prod
```

### Multi-Provider Support

Pull credentials from multiple sources in a single request:

- **Consul KV** — Base URLs, usernames, configuration
- **AWS Secrets Manager** — Passwords, API keys, sensitive data
- **HashiCorp Vault** — Coming soon

### Environment Switching

Seamlessly switch between dev, staging, and production:

```bash
sreq run GET /status -s auth -e dev      # Development
sreq run GET /status -s auth -e staging  # Staging
sreq run GET /status -s auth -e prod     # Production
```

### Offline Mode

Cache credentials locally for faster requests and offline use:

```bash
sreq sync dev           # Cache dev credentials
sreq run GET /api -s auth --offline  # Works without network
```

### Request History

Track, replay, and export previous requests:

```bash
sreq history                # View recent requests
sreq history 5 --replay     # Replay request #5
sreq history 5 --curl       # Export as curl command
```

### Interactive TUI

Build and execute requests with a terminal UI:

```bash
sreq tui
```

## Quick Example

**1. Initialize configuration:**

```bash
sreq init
```

**2. Configure providers:**

```bash
sreq auth
```

**3. Add a service:**

```bash
sreq service add auth-service --consul-key auth --aws-prefix auth-svc
```

**4. Make requests:**

```bash
sreq run GET /api/v1/users -s auth-service -e dev
```

## Why "sreq"?

**s**ervice + **req**uest = **sreq**

A CLI tool that makes service-aware HTTP requests with automatic credential resolution.

## Next Steps

- [Installation](/sreq/installation) — Get sreq installed
- [Getting Started](/sreq/getting-started) — Your first request in 5 minutes
- [Configuration](/sreq/configuration) — Deep dive into config options
