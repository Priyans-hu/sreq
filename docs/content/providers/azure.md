# Azure Key Vault

> **Coming Soon** — This provider is planned but not yet implemented.

Azure Key Vault provider will enable fetching secrets from Microsoft Azure's Key Vault service.

## Planned Configuration

```yaml
providers:
  azure:
    vault_url: https://my-vault.vault.azure.net
    paths:
      password: "{service}-{env}-password"
      api_key: "{service}-api-key"
```

## Expected Features

- Azure Identity SDK authentication (managed identity, CLI, environment)
- Secret versioning support
- Certificate and key retrieval
- Soft-delete and purge protection awareness

---

**Want this sooner?** Contributions are welcome!

- [Open an issue](https://github.com/Priyans-hu/sreq/issues) to discuss the design
- [View the roadmap](https://github.com/Priyans-hu/sreq/blob/main/docs/ROADMAP.md) for planned features
- [Contributing guide](https://github.com/Priyans-hu/sreq/blob/main/CONTRIBUTING.md)
