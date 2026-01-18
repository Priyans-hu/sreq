# sreq - AI Assistant Instructions

## Project Overview
sreq is a service-aware API client CLI that auto-fetches credentials from Consul, AWS Secrets Manager, and other providers.

## Local Tracking Files

**IMPORTANT:** Always maintain these local files while working on this project.

### Files to Update (gitignored, local only)

| File | Purpose | When to Update |
|------|---------|----------------|
| `.todos.md` | Active tasks and progress | Before starting, after completing, when planning |
| `.brainstorm.md` | Ideas, designs, decisions | When discussing new features or approaches |
| `.notes.md` | Session notes, blockers | During work sessions |

### Update Format

Always include timestamps in ISO format:

```markdown
## 2026-01-19T10:30:00

### Completed
- [x] Task description

### In Progress
- [ ] Task description

### Notes
- Decision made: xyz
```

## When Working on Tasks

### Before Starting
1. Read `.todos.md` to understand current state
2. Add new task with timestamp if not exists
3. Mark task as "In Progress"

### After Completing
1. Mark task as completed with timestamp
2. Add any follow-up tasks discovered
3. Update `.brainstorm.md` if new ideas emerged
4. Update `docs/ROADMAP.md` if milestone completed (this one is committed)

### When Planning
1. Add ideas to `.brainstorm.md` with timestamp
2. Break down into tasks in `.todos.md`
3. Note any decisions or tradeoffs

## Project Structure

```
sreq/
├── cmd/sreq/           # CLI entry (Cobra)
├── internal/
│   ├── client/         # HTTP client
│   ├── config/         # Config loading
│   └── providers/      # Secret providers (Consul, AWS, etc.)
├── pkg/types/          # Shared types
├── docs/               # Documentation (committed)
│   └── ROADMAP.md      # Public roadmap
└── [local files]       # Not committed
    ├── .todos.md
    ├── .brainstorm.md
    └── .notes.md
```

## Code Style

- Go 1.21+
- Use `gofmt` for formatting
- Follow Effective Go guidelines
- Keep functions small and focused
- Error wrapping with `fmt.Errorf("context: %w", err)`

## Build & Test

```bash
go build -o sreq ./cmd/sreq    # Build
go test ./...                   # Test
./sreq --help                   # Run
```

## Commit Convention

```
<type>(<scope>): <description>

Types: feat, fix, docs, style, refactor, test, chore
Example: feat(consul): implement KV provider
```
