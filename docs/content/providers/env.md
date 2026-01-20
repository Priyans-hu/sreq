# Environment Variables Provider

The `env` provider reads credentials directly from environment variables. This is useful for local development, CI/CD pipelines, and containerized environments.

## Configuration

```yaml
providers:
  env:
    prefix: "SREQ_"  # Optional prefix for all env vars
    paths:
      api_key: "{SERVICE}_API_KEY"
      base_url: "{SERVICE}_{ENV}_URL"
```

### Options

| Option | Type | Description |
|--------|------|-------------|
| `prefix` | string | Optional prefix added to all environment variable names |
| `paths` | map | Path templates for credential resolution |

## Path Templates

Templates support placeholders that are replaced with context values:

| Placeholder | Description |
|-------------|-------------|
| `{service}` | Service name |
| `{env}` | Environment (dev, staging, prod) |
| `{region}` | Region name |
| `{project}` | Project name |

Values are automatically:
- Converted to uppercase
- Hyphens replaced with underscores
- Dots replaced with underscores

### Example

Template: `{SERVICE}_{ENV}_API_KEY`
With: `service=auth-service`, `env=prod`
Result: `AUTH_SERVICE_PROD_API_KEY`

## Usage

### Simple Mode

```yaml
providers:
  env:
    prefix: "APP_"

services:
  billing:
    consul_key: billing  # Uses env provider with prefix
```

### Advanced Mode

```yaml
services:
  billing:
    paths:
      api_key: "env:BILLING_API_KEY"
      base_url: "env:{SERVICE}_{ENV}_URL"
```

## Setting Environment Variables

```bash
# Direct variable
export BILLING_API_KEY="sk-123456"

# With prefix
export APP_BILLING_API_KEY="sk-123456"

# With template placeholders
export BILLING_PROD_URL="https://api.billing.example.com"
```

## Use Cases

### Local Development

```bash
export AUTH_API_KEY="dev-key-123"
export AUTH_DEV_URL="http://localhost:8080"

sreq run auth GET /health -e dev
```

### CI/CD Pipelines

```yaml
# GitHub Actions
env:
  AUTH_API_KEY: ${{ secrets.AUTH_API_KEY }}
  AUTH_PROD_URL: ${{ vars.AUTH_PROD_URL }}

steps:
  - run: sreq run auth GET /api/test -e prod
```

### Docker/Kubernetes

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - AUTH_API_KEY=${AUTH_API_KEY}
      - AUTH_PROD_URL=https://auth.internal
```

```yaml
# Kubernetes Secret
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
data:
  AUTH_API_KEY: c2stMTIzNDU2  # base64 encoded
```

## Security Notes

- Environment variables can be visible in process listings
- Use secrets management in production (AWS Secrets Manager, Vault, etc.)
- Never commit `.env` files with real credentials
- Consider using the `dotenv` provider for local development instead
