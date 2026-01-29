---
title: upgrade
description: Self-update sreq to the latest version
order: 10
---

# sreq upgrade

Update sreq to the latest version from GitHub Releases.

## Synopsis

```bash
sreq upgrade [flags]
```

## Description

The `upgrade` command checks for a newer version of sreq on GitHub and downloads it in-place. It detects your OS and architecture automatically.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Install even if already on latest version | `false` |
| `--check` | Check for updates without installing | `false` |

## Examples

### Update to Latest

```bash
sreq upgrade
```

Output:

```
Current version: v0.1.0
Latest version:  v0.2.0
Downloading sreq v0.2.0 for darwin/arm64...
Updated successfully!
```

### Check for Updates

```bash
sreq upgrade --check
```

Output:

```
Current version: v0.1.0
Latest version:  v0.2.0
Update available! Run 'sreq upgrade' to install.
```

### Force Reinstall

```bash
sreq upgrade --force
```

## How It Works

1. Queries the GitHub Releases API for the latest version
2. Compares with the currently installed version
3. Downloads the correct binary for your OS/architecture
4. Replaces the existing binary in-place

## Alternative Update Methods

### Homebrew

```bash
brew upgrade sreq
```

### Go Install

```bash
go install github.com/Priyans-hu/sreq/cmd/sreq@latest
```

## See Also

- [Installation](/sreq/installation) — All installation methods
- [version](/sreq/commands/version) — Show current version
