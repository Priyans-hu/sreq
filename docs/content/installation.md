---
title: Installation
description: How to install sreq on your system
order: 2
---

# Installation

sreq is available for macOS, Linux, and Windows on both AMD64 and ARM64 architectures.

## Quick Install (curl)

```bash
curl -fsSL https://raw.githubusercontent.com/Priyans-hu/sreq/main/install.sh | bash
```

Detects your OS and architecture, downloads the correct binary, and installs to `/usr/local/bin`.

## Homebrew (macOS/Linux)

```bash
brew install Priyans-hu/tap/sreq
```

## Go Install

If you have Go 1.21+ installed:

```bash
go install github.com/Priyans-hu/sreq/cmd/sreq@latest
```

Make sure `$GOPATH/bin` is in your PATH.

## Binary Download

Download the latest release for your platform from [GitHub Releases](https://github.com/Priyans-hu/sreq/releases).

### macOS (Apple Silicon)

```bash
curl -L https://github.com/Priyans-hu/sreq/releases/latest/download/sreq_darwin_arm64.tar.gz | tar xz
sudo mv sreq /usr/local/bin/
```

### macOS (Intel)

```bash
curl -L https://github.com/Priyans-hu/sreq/releases/latest/download/sreq_darwin_amd64.tar.gz | tar xz
sudo mv sreq /usr/local/bin/
```

### Linux (AMD64)

```bash
curl -L https://github.com/Priyans-hu/sreq/releases/latest/download/sreq_linux_amd64.tar.gz | tar xz
sudo mv sreq /usr/local/bin/
```

### Linux (ARM64)

```bash
curl -L https://github.com/Priyans-hu/sreq/releases/latest/download/sreq_linux_arm64.tar.gz | tar xz
sudo mv sreq /usr/local/bin/
```

### Windows

Download `sreq_windows_amd64.zip` from [releases](https://github.com/Priyans-hu/sreq/releases), extract, and add to your PATH.

## Verify Installation

```bash
sreq version
```

You should see output like:

```
sreq version 0.1.0
```

## Updating

### Homebrew

```bash
brew upgrade sreq
```

### Self-Update

sreq can update itself:

```bash
sreq upgrade           # Update to latest version
sreq upgrade --check   # Check for updates without installing
```

## Shell Completions

Generate shell completions for your shell:

### Bash

```bash
sreq completion bash > /etc/bash_completion.d/sreq
```

### Zsh

```bash
sreq completion zsh > "${fpath[1]}/_sreq"
```

### Fish

```bash
sreq completion fish > ~/.config/fish/completions/sreq.fish
```

## Prerequisites

sreq requires credentials for the providers you want to use:

| Provider | Requirement |
|----------|-------------|
| Consul | `CONSUL_HTTP_TOKEN` or token in config |
| AWS Secrets Manager | AWS credentials (profile, env vars, or IAM role) |
| Environment Variables | Variables set in shell |
| Dotenv | `.env` file in project directory |

See [Configuration](/configuration) for detailed setup instructions.

## Next Steps

- [Getting Started](/getting-started) — Make your first request
- [Configuration](/configuration) — Set up providers and services
