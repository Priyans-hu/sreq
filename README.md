# sreq

[![CI](https://github.com/Priyans-hu/sreq/actions/workflows/ci.yml/badge.svg)](https://github.com/Priyans-hu/sreq/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Priyans-hu/sreq/graph/badge.svg)](https://codecov.io/gh/Priyans-hu/sreq)
[![Go Report Card](https://goreportcard.com/badge/github.com/Priyans-hu/sreq)](https://goreportcard.com/report/github.com/Priyans-hu/sreq)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Service-aware API client with automatic credential resolution.**

sreq eliminates the overhead of manually fetching credentials from multiple sources when testing APIs. Just specify the service name and environment — sreq handles the rest.

> If you find this useful, consider giving it a [⭐ star on GitHub](https://github.com/Priyans-hu/sreq) — it helps others discover the project!

## The Problem

```
You want to test: POST /api/v1/users on auth-service

Current workflow:
1. Open Consul → find auth-service base URL for dev
2. Open AWS Secrets Manager → find auth-service credentials
3. Copy-paste into Postman/curl
4. Repeat for staging... repeat for prod...
```

## The Solution

```bash
sreq POST /api/v1/users -s auth-service -e dev -d '{"name":"test"}'

# sreq automatically:
# ✓ Fetches base URL from Consul
# ✓ Fetches credentials from AWS Secrets Manager
# ✓ Caches credentials locally (AES-256 encrypted)
# ✓ Makes the request and saves it to history
```

## Features

- **Service-aware** — Pass service name, sreq resolves everything
- **Multi-provider** — Consul, AWS Secrets Manager, Environment Variables, Dotenv files
- **Environment switching** — Seamlessly switch between dev, staging, prod
- **Credential caching** — AES-256 encrypted local cache with TTL, offline mode support
- **Request history** — Track, replay, and export requests as curl/HTTPie
- **Interactive TUI** — Terminal UI for browsing services and history
- **Self-update** — `sreq upgrade` fetches the latest version from GitHub
- **Context system** — Save presets for env/region/project/app combinations
- **Health checks** — Test provider connectivity before making requests
- **Zero copy-paste** — No more manual credential hunting
- **Git-friendly config** — Share service configs with your team

## Installation

### Quick Install (curl)

```bash
curl -fsSL https://raw.githubusercontent.com/Priyans-hu/sreq/main/install.sh | bash
```

### Homebrew

```bash
brew install Priyans-hu/tap/sreq
```

### Go

```bash
go install github.com/Priyans-hu/sreq/cmd/sreq@latest
```

### Binary

Download from [Releases](https://github.com/Priyans-hu/sreq/releases).

## Quick Start

### 1. Initialize config

```bash
sreq init
```

This creates `~/.sreq/config.yaml`:

```yaml
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_HTTP_TOKEN}
    paths:
      base_url: "services/{service}/config/base_url"
      username: "services/{service}/config/username"

  aws_secrets:
    region: us-east-1
    paths:
      password: "{service}/{env}/password"
      api_key: "{service}/{env}/api_key"

environments:
  - dev
  - staging
  - prod

default_env: dev
```

### 2. Set up authentication

```bash
# Interactive setup for all providers
sreq auth

# Or configure a specific provider
sreq auth consul
sreq auth aws
```

### 3. Add a service

```bash
sreq service add auth-service
```

Or edit `~/.sreq/services.yaml`:

```yaml
services:
  auth-service:
    consul_key: auth
    aws_prefix: auth-service

  billing-service:
    consul_key: billing
    aws_prefix: billing-svc
```

### 4. Make requests

```bash
# GET request
sreq GET /api/v1/users -s auth-service -e dev

# POST with body
sreq POST /api/v1/users -s auth-service -e prod -d '{"name":"test"}'

# POST with body from file
sreq POST /api/v1/users -s auth-service -e staging -d @body.json

# With headers
sreq GET /api/v1/users -s auth-service -H "X-Custom: value"

# Verbose output (shows resolved URLs/creds)
sreq GET /api/v1/users -s auth-service -v

# Dry run (show what would be sent without executing)
sreq GET /api/v1/users -s auth-service --dry-run

# Use cached credentials only (no provider calls)
sreq GET /api/v1/users -s auth-service --offline
```

## Commands

| Command | Description |
|---------|-------------|
| `sreq init` | Initialize configuration |
| `sreq <METHOD> <path>` | Make HTTP request |
| `sreq service list` | List configured services |
| `sreq service add <name>` | Add a new service |
| `sreq service remove <name>` | Remove a service |
| `sreq env list` | List environments |
| `sreq env switch <env>` | Switch default environment |
| `sreq env current` | Show current default environment |
| `sreq auth` | Interactive authentication setup |
| `sreq auth consul` | Configure Consul authentication |
| `sreq auth aws` | Configure AWS authentication |
| `sreq config show` | Show current configuration |
| `sreq config path` | Show config file path |
| `sreq config test` | Test provider connectivity |
| `sreq cache status` | Show cache status and entries |
| `sreq cache clear [env]` | Clear cached credentials |
| `sreq sync [env]` | Sync credentials to local cache |
| `sreq history` | List request history |
| `sreq history <id>` | View a specific request |
| `sreq history <id> --replay` | Replay a previous request |
| `sreq history <id> --curl` | Export as curl command |
| `sreq history <id> --httpie` | Export as HTTPie command |
| `sreq history --clear` | Clear request history |
| `sreq tui` | Open interactive terminal UI |
| `sreq upgrade` | Update to latest version |
| `sreq version` | Show version |

## Flags

### Global

| Flag | Short | Description |
|------|-------|-------------|
| `--service` | `-s` | Service name |
| `--env` | `-e` | Environment (dev/staging/prod) |
| `--context` | `-c` | Context preset (overrides env/region/project/app) |
| `--region` | `-r` | Region |
| `--project` | `-p` | Project name |
| `--app` | `-a` | App name |
| `--verbose` | `-v` | Show detailed output |

### Request

| Flag | Short | Description |
|------|-------|-------------|
| `--data` | `-d` | Request body (or @filename for file) |
| `--header` | `-H` | Add header (repeatable) |
| `--output` | `-o` | Output format (json/raw/headers) |
| `--timeout` | | Request timeout (default: 30s) |
| `--dry-run` | | Show what would be sent without executing |
| `--offline` | | Use cached credentials only |
| `--no-cache` | | Skip cache, fetch fresh credentials |

### History

| Flag | Description |
|------|-------------|
| `--service` | Filter by service name |
| `--env` | Filter by environment |
| `--method` | Filter by HTTP method |
| `--all` | Show all history entries |
| `--clear` | Clear history |
| `--before` | Clear entries older than duration (7d, 24h) |
| `--curl` | Export as curl command |
| `--httpie` | Export as HTTPie command |
| `--replay` | Replay the request |

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `SREQ_CONFIG` | Config file path (default: `~/.sreq/config.yaml`) |
| `SREQ_NO_CACHE` | Disable credential caching |
| `SREQ_NO_HISTORY` | Disable request history |
| `CONSUL_HTTP_TOKEN` | Consul ACL token |
| `AWS_PROFILE` | AWS profile for Secrets Manager |
| `AWS_REGION` | AWS region override |

### Provider Configuration

<details>
<summary><b>Consul</b></summary>

```yaml
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_HTTP_TOKEN}
    datacenter: dc1  # optional
    env_addresses:   # optional, per-environment addresses
      prod: consul-prod.internal:8500
    paths:
      base_url: "services/{service}/base_url"
      username: "services/{service}/username"
```

Path placeholders: `{service}`, `{env}`, `{region}`, `{project}`, `{app}`

</details>

<details>
<summary><b>AWS Secrets Manager</b></summary>

```yaml
providers:
  aws_secrets:
    region: us-east-1
    profile: default  # optional, uses AWS_PROFILE
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"
```

Use `#` to extract a JSON key from within a secret value.

</details>

<details>
<summary><b>Environment Variables</b></summary>

```yaml
providers:
  env:
    paths:
      base_url: "{SERVICE}_BASE_URL"
      username: "{SERVICE}_USERNAME"
      password: "{SERVICE}_PASSWORD"
      api_key: "{SERVICE}_API_KEY"
```

Reads credentials directly from environment variables.

</details>

<details>
<summary><b>Dotenv Files</b></summary>

```yaml
providers:
  dotenv:
    path: ".env.{env}"  # resolves to .env.dev, .env.prod, etc.
    paths:
      base_url: "BASE_URL"
      username: "USERNAME"
      password: "PASSWORD"
```

Reads credentials from `.env` files with environment-based path resolution.

</details>

<details>
<summary><b>HashiCorp Vault (planned)</b></summary>

```yaml
providers:
  vault:
    address: https://vault.example.com:8200
    token: ${VAULT_TOKEN}
    paths:
      password: "secret/data/{service}/{env}#password"
```

</details>

### Contexts

Contexts are presets that bundle env, region, project, and app into a single name:

```yaml
contexts:
  dev-us:
    env: dev
    region: us-east-1
    project: myproject
    app: myapp
  prod-eu:
    env: prod
    region: eu-west-1
    project: myproject
    app: myapp

default_context: dev-us
```

```bash
# Use a context instead of individual flags
sreq GET /api/v1/users -s auth-service -c prod-eu
```

## Project Structure

```
sreq/
├── cmd/sreq/             # CLI entrypoint and commands
├── internal/
│   ├── cache/            # AES-256 encrypted credential cache
│   ├── client/           # HTTP client
│   ├── config/           # Configuration loading
│   ├── errors/           # Custom error types
│   ├── history/          # Request history tracking
│   ├── providers/        # Credential providers
│   │   ├── consul/       # Consul KV provider
│   │   ├── aws/          # AWS Secrets Manager provider
│   │   ├── env/          # Environment variables provider
│   │   └── dotenv/       # Dotenv file provider
│   ├── resolver/         # Multi-provider resolution
│   └── tui/              # Terminal UI (BubbleTea)
├── pkg/types/            # Shared types
├── docs/                 # Documentation site (Docsify)
├── website/              # Landing page (Next.js)
├── .goreleaser.yml       # Cross-platform release automation
├── install.sh            # Quick install script
└── codecov.yml           # Coverage config
```

## Documentation

Full documentation is available at the [docs site](https://priyans-hu.github.io/sreq/).

## Roadmap

- [x] Project setup & CI/CD
- [x] Consul KV provider
- [x] AWS Secrets Manager provider
- [x] Environment variables provider
- [x] Dotenv file provider
- [x] Credential caching (AES-256, offline mode)
- [x] Request history (replay, curl/HTTPie export)
- [x] Interactive TUI mode
- [x] Interactive auth setup
- [x] Self-update command
- [x] Codecov integration
- [x] Documentation site
- [x] Landing page
- [x] Cross-platform release automation
- [ ] HashiCorp Vault provider
- [ ] TUI clipboard copy support
- [ ] Integration tests

See [docs/ROADMAP.md](docs/ROADMAP.md) for detailed planning.

## Contributing

Contributions are welcome! Please read the [Contributing Guide](CONTRIBUTING.md) first.

```bash
# Clone
git clone https://github.com/Priyans-hu/sreq.git
cd sreq

# Install dependencies
go mod download

# Build
go build -o sreq ./cmd/sreq

# Test
go test ./...

# Run
./sreq --help
```

## Why "sreq"?

**s**ervice + **req**uest = **sreq**

A CLI tool that makes service-aware HTTP requests with automatic credential resolution.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

Built by [Priyanshu](https://github.com/Priyans-hu) · [LinkedIn](https://linkedin.com/in/priyans-hu)
