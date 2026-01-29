---
title: Commands
description: Complete command reference for sreq
order: 5
---

# Commands

Complete reference for all sreq commands.

## Command Overview

| Command | Description |
|---------|-------------|
| [`sreq run`](/commands/run) | Make HTTP requests |
| [`sreq init`](/commands/init) | Initialize configuration |
| [`sreq auth`](/commands/auth) | Configure provider authentication |
| [`sreq service`](/commands/service) | Manage service configurations |
| [`sreq env`](/commands/env) | Manage environments |
| [`sreq config`](/commands/config) | View and validate configuration |
| [`sreq history`](/commands/history) | View and manage request history |
| [`sreq cache`](/commands/cache) | Manage credential cache |
| [`sreq sync`](/commands/cache#sync) | Sync credentials to cache |
| [`sreq tui`](/commands/tui) | Interactive terminal UI |
| `sreq version` | Show version |
| `sreq upgrade` | Self-update to latest version |
| `sreq completion` | Generate shell completions |

## Global Flags

These flags work with most commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--service` | `-s` | Service name |
| `--env` | `-e` | Environment (dev/staging/prod) |
| `--context` | `-c` | Use a named context |
| `--region` | `-r` | Override region |
| `--project` | `-p` | Override project |
| `--app` | `-a` | Override app |
| `--verbose` | `-v` | Verbose output |
| `--dry-run` | | Preview without executing |
| `--help` | `-h` | Help for command |

## Quick Examples

```bash
# Make a GET request
sreq run GET /api/v1/users -s auth-service -e dev

# Make a POST request with body
sreq run POST /api/v1/users -s auth-service -d '{"name":"test"}'

# List services
sreq service list

# View request history
sreq history

# Replay a previous request
sreq history 5 --replay

# Cache credentials for offline use
sreq sync dev

# Launch interactive TUI
sreq tui
```
