# HashiCorp Vault

> **Coming Soon** — This provider is planned but not yet implemented.

HashiCorp Vault provider will enable fetching secrets from Vault's KV v2 secret engine.

## Planned Configuration

```yaml
providers:
  vault:
    address: https://vault.example.com:8200
    token: ${VAULT_TOKEN}
    paths:
      password: "secret/data/{service}/{env}#password"
      api_key: "secret/data/{service}/{env}#api_key"
```

## Expected Features

- KV v2 secret engine support
- Token and AppRole authentication
- Namespace support for Vault Enterprise
- Dynamic secret leasing
- Auto-renewal of leased secrets

---

**Want this sooner?** Contributions are welcome!

- [Open an issue](https://github.com/Priyans-hu/sreq/issues) to discuss the design
- [View the roadmap](https://github.com/Priyans-hu/sreq/blob/main/docs/ROADMAP.md) for planned features
- [Contributing guide](https://github.com/Priyans-hu/sreq/blob/main/CONTRIBUTING.md)
