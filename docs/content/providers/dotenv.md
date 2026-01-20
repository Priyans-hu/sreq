# dotenv Provider

The `dotenv` provider reads credentials from `.env` files. This is ideal for local development when you want to keep credentials in files rather than setting environment variables manually.

## Configuration

```yaml
providers:
  dotenv:
    file: ".env"           # Single file (backward compatible)
    files:                 # Multiple files (later overrides earlier)
      - ".env"
      - ".env.local"
      - ".env.{env}"
    paths:
      api_key: "{SERVICE}_API_KEY"
```

### Options

| Option | Type | Description |
|--------|------|-------------|
| `file` | string | Single .env file path |
| `files` | list | Multiple .env files to load (in order) |
| `paths` | map | Path templates for credential resolution |

## File Format

Standard `.env` file format is supported:

```bash
# Comments start with #
API_KEY=your-api-key

# Quoted values (spaces preserved)
DATABASE_URL="postgres://localhost:5432/mydb"
SECRET='value with spaces'

# Export prefix (optional)
export AUTH_TOKEN=token123

# Escape sequences in double quotes
MULTILINE="line1\nline2"
TABBED="col1\tcol2"
```

## Multiple Files

When multiple files are specified, they are loaded in order. Later files override values from earlier files:

```yaml
providers:
  dotenv:
    files:
      - ".env"           # Base config
      - ".env.local"     # Local overrides (gitignored)
      - ".env.prod"      # Environment-specific
```

**Load order example:**

```bash
# .env
API_KEY=default-key
DATABASE_URL=postgres://localhost/db

# .env.local
API_KEY=my-local-key

# Result: API_KEY=my-local-key, DATABASE_URL=postgres://localhost/db
```

## Path Templates

Templates work the same as the `env` provider:

| Placeholder | Description |
|-------------|-------------|
| `{service}` | Service name |
| `{env}` | Environment |
| `{region}` | Region name |
| `{project}` | Project name |

Values are converted to uppercase with hyphens/dots replaced by underscores.

## Usage

### Basic Setup

Create a `.env` file in your project root:

```bash
# .env
BILLING_API_KEY=sk-test-123
BILLING_DEV_URL=https://api-dev.billing.example.com
BILLING_PROD_URL=https://api.billing.example.com
```

Configure sreq:

```yaml
providers:
  dotenv:
    file: ".env"
    paths:
      api_key: "{SERVICE}_API_KEY"
      base_url: "{SERVICE}_{ENV}_URL"

services:
  billing:
    consul_key: billing
```

Run requests:

```bash
sreq run billing GET /invoices -e dev
```

### Advanced Mode

Reference dotenv directly in service paths:

```yaml
services:
  billing:
    paths:
      api_key: "dotenv:BILLING_API_KEY"
      base_url: "dotenv:{SERVICE}_{ENV}_URL"
```

### Home Directory Files

Support for `~/` prefix:

```yaml
providers:
  dotenv:
    files:
      - "~/.sreq/.env"      # Global credentials
      - ".env"              # Project-specific
      - ".env.local"        # Local overrides
```

## Best Practices

### Git Ignore

Add to `.gitignore`:

```gitignore
# Environment files with secrets
.env.local
.env.*.local
.env.production

# Keep example file
!.env.example
```

### Example File

Create `.env.example` with placeholder values:

```bash
# .env.example - Copy to .env and fill in values
BILLING_API_KEY=your-api-key-here
BILLING_DEV_URL=https://api-dev.example.com
BILLING_PROD_URL=https://api.example.com
```

### Environment-Specific Files

Structure for multiple environments:

```
.env              # Shared defaults
.env.local        # Local overrides (gitignored)
.env.development  # Dev environment
.env.staging      # Staging environment
.env.production   # Production (gitignored)
```

## Reloading

The dotenv provider caches loaded values. To reload after file changes:

```bash
# Clear cache and reload
sreq cache clear
```

Or programmatically via the `Reload()` method.

## Security Notes

- Never commit `.env` files with real credentials
- Use `.env.example` for documentation
- Add sensitive `.env` files to `.gitignore`
- For production, prefer proper secrets management (AWS Secrets Manager, Vault)
- File permissions should be restricted (`chmod 600 .env`)
