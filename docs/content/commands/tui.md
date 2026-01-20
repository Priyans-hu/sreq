---
title: tui
description: Interactive terminal UI for sreq
order: 7
---

# sreq tui

Launch the interactive terminal UI.

## Synopsis

```bash
sreq tui
```

## Description

The TUI (Terminal User Interface) provides an interactive way to build and execute requests. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), it offers a modern terminal experience.

## Features

- **Service browser** — Browse and select configured services
- **Request builder** — Build requests with method, path, headers, and body
- **History viewer** — Browse and replay past requests
- **Environment switcher** — Quick environment switching

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch between panels |
| `Enter` | Select / Execute |
| `Esc` | Back / Cancel |
| `q` | Quit |
| `?` | Help |
| `/` | Search |
| `r` | Replay selected request |
| `c` | Copy as curl |

## Panels

### Services Panel

Browse configured services:

```
┌─ Services ─────────────────┐
│ > auth-service             │
│   billing-service          │
│   user-service             │
│   notification-service     │
└────────────────────────────┘
```

### Request Builder

Build and execute requests:

```
┌─ Request ──────────────────┐
│ Method: [GET ▼]            │
│ Path:   /api/v1/users      │
│ Env:    [dev ▼]            │
│                            │
│ Headers:                   │
│   Content-Type: app/json   │
│                            │
│ Body:                      │
│   {"name": "test"}         │
│                            │
│ [Execute]  [Clear]         │
└────────────────────────────┘
```

### Response Panel

View response details:

```
┌─ Response ─────────────────┐
│ Status: 200 OK (145ms)     │
│                            │
│ {                          │
│   "users": [               │
│     {"id": 1, "name": ...} │
│   ]                        │
│ }                          │
└────────────────────────────┘
```

### History Panel

Browse request history:

```
┌─ History ──────────────────┐
│ #1 GET  /api/users    200  │
│ #2 POST /api/users    201  │
│ #3 GET  /api/health   200  │
└────────────────────────────┘
```

## Requirements

- Terminal with 256-color support
- Minimum 80x24 terminal size

## See Also

- [run](/sreq/commands/run) — CLI request execution
- [history](/sreq/commands/history) — CLI history management
