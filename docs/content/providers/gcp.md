# GCP Secret Manager

> **Coming Soon** — This provider is planned but not yet implemented.

GCP Secret Manager provider will enable fetching secrets from Google Cloud's Secret Manager service.

## Planned Configuration

```yaml
providers:
  gcp:
    project: my-gcp-project
    paths:
      password: "projects/{project}/secrets/{service}-{env}-password/versions/latest"
      api_key: "projects/{project}/secrets/{service}-api-key/versions/latest"
```

## Expected Features

- Secret version management
- IAM-based access control via Application Default Credentials
- Automatic secret rotation support
- Regional and global secret support

---

**Want this sooner?** Contributions are welcome!

- [Open an issue](https://github.com/Priyans-hu/sreq/issues) to discuss the design
- [View the roadmap](https://github.com/Priyans-hu/sreq/blob/main/docs/ROADMAP.md) for planned features
- [Contributing guide](https://github.com/Priyans-hu/sreq/blob/main/CONTRIBUTING.md)
