---
title: history
description: View and manage request history
order: 5
---

# sreq history

View, replay, and manage request history.

## Synopsis

```bash
sreq history [id] [flags]
```

## Description

sreq automatically saves all requests to a local history file. Use the `history` command to view past requests, replay them, or export them as curl/HTTPie commands.

## Arguments

| Argument | Description |
|----------|-------------|
| `id` | Optional request ID to view details or perform actions |

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--all` | Show all history entries | `false` |
| `--service` | Filter by service name | — |
| `--env` | Filter by environment | — |
| `--method` | Filter by HTTP method | — |
| `--curl` | Export as curl command | `false` |
| `--httpie` | Export as HTTPie command | `false` |
| `--replay` | Replay the request | `false` |
| `--clear` | Clear history | `false` |
| `--before` | Clear entries older than duration | — |

## Examples

### List Recent Requests

```bash
sreq history
```

Output:

```
ID    METHOD  PATH                            SERVICE          ENV     STATUS    TIME
------------------------------------------------------------------------------------------
1     GET     /api/v1/users                   auth-service     dev     200       145ms
2     POST    /api/v1/users                   auth-service     dev     201       230ms
3     GET     /api/v1/invoices                billing-service  prod    200       89ms
4     DELETE  /api/v1/users/123               auth-service     dev     404       67ms

Showing 4 of 4 entries. Use --all to see all.
```

### List All History

```bash
sreq history --all
```

### Filter by Service

```bash
sreq history --service auth-service
```

### Filter by Environment

```bash
sreq history --env prod
```

### Filter by Method

```bash
sreq history --method POST
```

### View Request Details

```bash
sreq history 2
```

Output:

```
Request #2
========================================
Time:     2026-01-19T10:30:00Z
Service:  auth-service
Env:      dev
Method:   POST
Path:     /api/v1/users
Base URL: https://auth.example.com
Status:   201
Duration: 230ms

Request Headers:
  Content-Type: application/json

Request Body:
  {"name": "John Doe", "email": "john@example.com"}

Response: 201 Created (156 bytes)

Export: sreq history 2 --curl
Replay: sreq history 2 --replay
```

### Export as curl

```bash
sreq history 2 --curl
```

Output:

```bash
curl -X POST 'https://auth.example.com/api/v1/users' \
  -H 'Content-Type: application/json' \
  -u 'api_user:password' \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### Export as HTTPie

```bash
sreq history 2 --httpie
```

Output:

```bash
http POST 'https://auth.example.com/api/v1/users' \
  Content-Type:application/json \
  -a api_user:password \
  name="John Doe" email="john@example.com"
```

### Replay a Request

```bash
sreq history 2 --replay
```

Output:

```
Replaying request #2: POST /api/v1/users

{"id": "user_789", "name": "John Doe", ...}
```

### Clear All History

```bash
sreq history --clear
```

### Clear Old Entries

```bash
# Clear entries older than 7 days
sreq history --clear --before 7d

# Clear entries older than 24 hours
sreq history --clear --before 24h
```

## Duration Format

For `--before` flag:

| Format | Meaning |
|--------|---------|
| `7d` | 7 days |
| `24h` | 24 hours |
| `30m` | 30 minutes |
| `1w` | 1 week |

## Disabling History

Set the environment variable to disable history:

```bash
export SREQ_NO_HISTORY=1
```

## History Storage

History is stored in `~/.sreq/history.json`.

## See Also

- [run](/commands/run) — Make requests
- [cache](/commands/cache) — Manage credential cache
