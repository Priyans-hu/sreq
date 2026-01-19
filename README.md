# sreq

[![CI](https://github.com/Priyans-hu/sreq/actions/workflows/ci.yml/badge.svg)](https://github.com/Priyans-hu/sreq/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Priyans-hu/sreq/graph/badge.svg)](https://codecov.io/gh/Priyans-hu/sreq)
[![Go Report Card](https://goreportcard.com/badge/github.com/Priyans-hu/sreq)](https://goreportcard.com/report/github.com/Priyans-hu/sreq)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Service-aware API client with automatic credential resolution.**

sreq eliminates the overhead of manually fetching credentials from multiple sources when testing APIs. Just specify the service name and environment — sreq handles the rest.

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
# ✓ Makes the request
```

## Features

- **Service-aware** — Pass service name, sreq resolves everything
- **Multi-provider** — Consul, AWS Secrets Manager, HashiCorp Vault (planned)
- **Environment switching** — Seamlessly switch between dev, staging, prod
- **Zero copy-paste** — No more manual credential hunting
- **Git-friendly config** — Share service configs with your team
- **AI/LLM friendly** — Designed for integration with AI assistants

## Installation

### Homebrew (coming soon)

```bash
brew install sreq
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
    token: ${CONSUL_TOKEN}
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

### 2. Add a service

```bash
sreq service add auth-service \
  --consul-key auth \
  --aws-prefix auth-service
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

### 3. Make requests

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
```

## Commands

| Command | Description |
|---------|-------------|
| `sreq init` | Initialize configuration |
| `sreq <METHOD> <path>` | Make HTTP request |
| `sreq service list` | List configured services |
| `sreq service add <name>` | Add a new service |
| `sreq env list` | List environments |
| `sreq env switch <env>` | Switch default environment |
| `sreq config show` | Show current configuration |
| `sreq version` | Show version |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--service` | `-s` | Service name |
| `--env` | `-e` | Environment (dev/staging/prod) |
| `--data` | `-d` | Request body (or @filename) |
| `--header` | `-H` | Add header (repeatable) |
| `--verbose` | `-v` | Show detailed output |
| `--dry-run` | | Show what would be sent without executing |
| `--output` | `-o` | Output format (json/table/raw) |

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `SREQ_CONFIG` | Config file path (default: `~/.sreq/config.yaml`) |
| `CONSUL_TOKEN` | Consul ACL token |
| `AWS_PROFILE` | AWS profile for Secrets Manager |
| `AWS_REGION` | AWS region override |

### Provider Configuration

<details>
<summary><b>Consul</b></summary>

```yaml
providers:
  consul:
    address: consul.example.com:8500
    token: ${CONSUL_TOKEN}
    datacenter: dc1  # optional
    paths:
      base_url: "services/{service}/base_url"
      username: "services/{service}/username"
```

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

Note: Use `#` to access JSON keys within a secret.

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

## Project Structure

```
sreq/
├── cmd/
│   └── sreq/
│       └── main.go          # CLI entrypoint
├── internal/
│   ├── config/              # Configuration loading
│   ├── providers/           # Secret providers
│   │   ├── consul/          # Consul KV provider
│   │   └── aws/             # AWS Secrets Manager provider
│   └── client/              # HTTP client
├── pkg/
│   └── types/               # Shared types
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── llms.txt                 # AI/LLM context file
```

## Roadmap

- [x] Project setup
- [ ] Consul provider
- [ ] AWS Secrets Manager provider
- [ ] Credential caching (offline mode)
- [ ] HashiCorp Vault provider
- [ ] Request history
- [ ] TUI mode

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

# Run
./sreq --help
```

## Why "sreq"?

**s**ervice + **req**uest = **sreq**

A CLI tool that makes service-aware HTTP requests with automatic credential resolution.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

Built by [Priyanshu](https://priyans-hu.netlify.app) · [GitHub](https://github.com/Priyans-hu)
