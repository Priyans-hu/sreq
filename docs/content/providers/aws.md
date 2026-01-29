---
title: AWS Provider
description: AWS Secrets Manager provider configuration
order: 12
---

# AWS Secrets Manager Provider

The AWS provider fetches credentials from AWS Secrets Manager.

## Overview

AWS Secrets Manager is commonly used to store:

- Passwords and API keys
- Database credentials
- Sensitive configuration

## Configuration

### Basic Setup

```yaml
providers:
  aws_secrets:
    region: us-east-1
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"
```

### Full Configuration

```yaml
providers:
  aws_secrets:
    # AWS region
    region: us-east-1

    # AWS credentials profile (optional)
    profile: default

    # Optional endpoint override (for LocalStack, etc.)
    endpoint: ""

    # Path templates for credential resolution
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"
      token: "{service}/{env}/credentials#token"
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `region` | AWS region | `us-east-1` |
| `profile` | AWS credentials profile | `default` |
| `endpoint` | Custom endpoint (for testing) | — |
| `paths` | Path templates | — |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AWS_REGION` | AWS region |
| `AWS_PROFILE` | Credentials profile |
| `AWS_ACCESS_KEY_ID` | Access key (if not using profile) |
| `AWS_SECRET_ACCESS_KEY` | Secret key (if not using profile) |
| `AWS_SESSION_TOKEN` | Session token (for temporary credentials) |

## Authentication

AWS credentials can be provided via (in order of precedence):

1. **Environment variables**:

   ```bash
   export AWS_ACCESS_KEY_ID=AKIA...
   export AWS_SECRET_ACCESS_KEY=...
   ```

2. **AWS profile** (in `~/.aws/credentials`):

   ```yaml
   providers:
     aws_secrets:
       profile: myprofile
   ```

3. **IAM role** (on EC2/ECS/Lambda):
   No configuration needed — uses instance role automatically.

4. **SSO profile**:

   ```bash
   aws sso login --profile my-sso-profile
   ```

   ```yaml
   providers:
     aws_secrets:
       profile: my-sso-profile
   ```

## Path Templates

Templates support these placeholders:

| Placeholder | Description | Example Value |
|-------------|-------------|---------------|
| `{service}` | Service's `aws_prefix` or name | `auth-svc` |
| `{env}` | Current environment | `dev` |
| `{region}` | AWS region | `us-east-1` |

### JSON Key Extraction

Use `#` to extract a specific key from a JSON secret:

```yaml
paths:
  password: "{service}/{env}/credentials#password"
```

If the secret `auth-svc/dev/credentials` contains:

```json
{
  "password": "secret123",
  "api_key": "key456"
}
```

Then `#password` extracts `"secret123"`.

### Example Resolution

Configuration:

```yaml
providers:
  aws_secrets:
    region: us-east-1
    paths:
      password: "{service}/{env}/credentials#password"

services:
  auth-service:
    aws_prefix: auth-svc
```

Request:

```bash
sreq run GET /api -s auth-service -e dev
```

Secret queried: `auth-svc/dev/credentials`, key `password` extracted.

## IAM Permissions

The IAM user/role needs these permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": [
        "arn:aws:secretsmanager:us-east-1:123456789:secret:auth-svc/*",
        "arn:aws:secretsmanager:us-east-1:123456789:secret:billing-svc/*"
      ]
    }
  ]
}
```

For broader access (not recommended for production):

```json
{
  "Effect": "Allow",
  "Action": "secretsmanager:GetSecretValue",
  "Resource": "*"
}
```

## Testing Connection

Verify AWS connectivity:

```bash
sreq config test
```

Output:

```
AWS Secrets Manager:
  Region:  us-east-1
  Profile: default
  Status:  ✓ Credentials valid
```

## Local Development

### Using LocalStack

For local development with [LocalStack](https://localstack.cloud/):

```yaml
providers:
  aws_secrets:
    region: us-east-1
    endpoint: http://localhost:4566
```

### Using AWS Profile

```bash
# Configure AWS CLI profile
aws configure --profile dev

# Reference in sreq config
providers:
  aws_secrets:
    profile: dev
```

## Troubleshooting

### Credentials Not Found

```
Error: NoCredentialProviders: no valid providers in chain
```

**Solutions:**

- Set `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
- Configure AWS profile: `aws configure`
- On EC2: ensure instance has IAM role attached

### Access Denied

```
Error: AccessDeniedException: User is not authorized to perform secretsmanager:GetSecretValue
```

**Solutions:**

- Add `secretsmanager:GetSecretValue` permission to IAM policy
- Verify resource ARN matches your secrets
- Check if secret has resource policy blocking access

### Secret Not Found

```
Error: ResourceNotFoundException: Secrets Manager can't find the specified secret
```

**Solutions:**

- Verify secret name is correct
- Check region matches where secret is stored
- Verify path template produces correct secret name

### Invalid JSON Key

```
Error: key "password" not found in secret JSON
```

**Solutions:**

- Verify the JSON key exists in the secret value
- Check for typos in the `#key` suffix
- Ensure secret value is valid JSON

## See Also

- [Consul Provider](/providers/consul) — Consul KV setup
- [Configuration](/configuration) — Full configuration reference
