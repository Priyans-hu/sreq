# Roadmap

Detailed planning and feature specifications for sreq.

## Current Status

**Latest Release:** v0.1.0 (2026-01-19)

- [x] Project setup
- [x] Core CLI structure (Cobra)
- [x] Consul provider (with env-specific addresses)
- [x] AWS Secrets Manager provider
- [x] HTTP client with auth
- [x] Environment switching
- [x] Config file management
- [x] Interactive auth setup (`sreq auth`)
- [x] Error messages with suggestions
- [x] Request history with replay and export
- [x] TUI mode (Bubble Tea)
- [x] Encrypted credential caching
- [x] Automated releases (GitHub Actions + GoReleaser)

## Planned Features

### Phase 1: Core Functionality ✅

- [x] Consul KV provider
- [x] AWS Secrets Manager provider
- [x] HTTP client with Basic Auth and API key support
- [x] Environment switching (dev/staging/prod)
- [x] YAML configuration management
- [x] Context support (presets for project/env/region/app)

### Phase 2: Request History & TUI ✅

#### Request History

Save and replay previous requests:

```bash
sreq history                    # List recent requests (last 20)
sreq history --all              # List all
sreq history --service auth     # Filter by service
sreq history --env prod         # Filter by env
sreq history 5                  # Show details of request #5
sreq history 5 --replay         # Replay request #5
sreq history 5 --curl           # Export as curl command
sreq history --clear            # Clear all history
```

#### TUI Mode (Bubble Tea)

Interactive terminal UI for building and executing requests.

```bash
sreq tui
```

Features:
- Dashboard with service list and recent requests
- Request builder with dropdowns for service/env/method
- History browser with replay support
- Keyboard shortcuts for common actions

### Phase 3: Credential Caching ✅

Cache credentials locally for faster requests and offline use.

**Commands:**
```bash
sreq sync dev              # Sync credentials for an environment
sreq sync --all            # Sync all environments
sreq sync --force          # Force refresh ignoring TTL
sreq cache status          # Show cache status
sreq cache clear           # Clear all cache
```

**Security measures:**
- AES-256-GCM encryption with user-specific key
- File permissions `600` (owner read/write only)
- Configurable TTL (default: 1 hour)
- Environment variable to disable caching: `SREQ_NO_CACHE=1`
- Auto-disable in CI/CD detection (`CI=true`)

### Phase 4: Additional Providers

#### Priority Providers

| Provider | Status | Notes |
|----------|--------|-------|
| Environment Variables | Planned | Essential for local dev/CI |
| HashiCorp Vault | Planned | Enterprise demand, KV v2 API |
| dotenv files | Planned | Parse local `.env` files |
| GCP Secret Manager | Planned | Complete cloud trifecta |
| Azure Key Vault | Planned | Complete cloud trifecta |

#### Environment Variable Provider

```yaml
providers:
  env:
    paths:
      api_key: "{SERVICE}_API_KEY"
      password: "{SERVICE}_{ENV}_PASSWORD"
```

#### HashiCorp Vault Provider

```yaml
providers:
  vault:
    address: https://vault.example.com:8200
    token: ${VAULT_TOKEN}
    paths:
      password: "secret/data/{service}/{env}#password"
```

#### dotenv Provider

```yaml
providers:
  dotenv:
    file: .env.local
    paths:
      api_key: "{SERVICE}_API_KEY"
```

### Phase 5: Integrations & Export

#### Export Formats

| Format | Status | Notes |
|--------|--------|-------|
| curl | Planned | `sreq history 5 --curl` |
| HTTPie | Planned | `sreq history 5 --httpie` |
| Bruno | Planned | `sreq export bruno ./collection` |
| Postman | Future | Collection v2.1 format |

#### CI/CD Integration

| Platform | Status | Notes |
|----------|--------|-------|
| GitHub Actions | Planned | Action to run sreq in workflows |
| GitLab CI | Future | Template for .gitlab-ci.yml |

#### Package Managers

| Method | Status | Command |
|--------|--------|---------|
| Binary download | ✅ Available | [Releases page](https://github.com/Priyans-hu/sreq/releases) |
| Go install | ✅ Available | `go install github.com/Priyans-hu/sreq/cmd/sreq@latest` |
| Homebrew | 🔜 Coming | `brew tap Priyans-hu/tap && brew install sreq` |
| Scoop | Planned | `scoop install sreq` |

### Phase 6: Advanced Features

- **Auto-Discovery**: `sreq discover consul` - List keys in providers
- **Import**: `sreq import consul billing_service --as billing` - Auto-generate service config
- **Collections**: Group related requests together
- **Response filtering**: jq-style filtering (`--filter '.data[0].name'`)
- **Metrics/timing**: `--timing` flag for DNS/TLS/TTFB breakdown
- **Retry with backoff**: Auto-retry failed requests
- **Proxy support**: HTTP/SOCKS proxy configuration

## Future Provider Ideas

These are potential providers we may support based on demand:

| Provider | Complexity | Notes |
|----------|------------|-------|
| AWS Parameter Store | Low | Already have AWS SDK |
| Doppler | Low | REST API, popular for startups |
| Infisical | Low | Open-source, REST API |
| 1Password | Medium | CLI (`op`) or Connect API |
| Kubernetes Secrets | Medium | client-go library |
| etcd | Low | Similar to Consul KV |
| SOPS | Low | Encrypted YAML/JSON files |
| CyberArk Conjur | High | Enterprise, complex setup |

## Ideas / Backlog

- Request templates with variables
- Team config sharing via git
- VS Code extension
- JetBrains plugin
- fzf integration for fuzzy selection
- direnv integration for per-directory contexts
