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
5. Update `docs/SETUP.md` if relevant (see Documentation Maintenance below)

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

## Documentation Maintenance

**IMPORTANT:** Keep these docs updated when making relevant changes.

### `docs/SETUP.md` - Update when:
| Change Type | What to Update |
|-------------|----------------|
| New provider (Azure, Vault, etc.) | Add to Optional prerequisites, provider config examples, troubleshooting |
| New command | Add to Quick Reference section |
| New environment variable | Add to Environment Variables table |
| Config format change | Update configuration examples |
| New installation method | Update Installation section |
| Release/binary changes | Update download URLs if naming changes |

### `docs/ROADMAP.md` - Update when:
- Feature milestone completed
- New features planned
- Priorities change

### `README.md` - Update when:
- Major new features added
- Installation method changes
- Quick start examples need updating

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

## Git Workflow

### Branch Naming

Always use prefixed branch names:

| Prefix | Purpose | Example |
|--------|---------|---------|
| `feat/` | New feature | `feat/consul-provider` |
| `fix/` | Bug fix | `fix/config-parsing` |
| `refactor/` | Code refactoring | `refactor/provider-interface` |
| `docs/` | Documentation | `docs/api-examples` |
| `test/` | Adding tests | `test/consul-integration` |
| `chore/` | Maintenance | `chore/update-deps` |

```bash
# Create branch
git checkout -b feat/consul-provider

# Always branch from latest main
git checkout main && git pull && git checkout -b feat/new-feature
```

### Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
```

**Types:**
| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no code change |
| `refactor` | Code change, no feature/fix |
| `test` | Adding/updating tests |
| `chore` | Maintenance, deps, CI |

**Scopes:** `consul`, `aws`, `vault`, `config`, `client`, `cli`, `cache`

**Examples:**
```bash
feat(consul): implement KV provider with path templates
fix(aws): handle pagination in secret listing
docs(readme): add installation instructions
refactor(providers): extract common path resolution logic
test(consul): add integration tests for KV fetch
chore(deps): update aws-sdk-go to v2.5.0
```

### PR Workflow

1. Create feature branch from `main`
2. Make commits following convention
3. Push and create PR to `main`
4. After merge, delete feature branch

```bash
# Typical flow
git checkout main && git pull
git checkout -b feat/consul-provider
# ... make changes ...
git add . && git commit -m "feat(consul): implement KV provider"
git push -u origin feat/consul-provider
# Create PR on GitHub
```
