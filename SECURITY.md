# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Security Best Practices

When using this provider:

1. **Protect API Keys**: Never commit API keys to version control. Use environment variables or a secrets manager.

2. **Use Terraform State Encryption**: Enable encryption for your Terraform state, especially in remote backends.

3. **Restrict State Access**: Limit who can access Terraform state files, as they may contain sensitive values.

4. **Review Plans**: Always review `terraform plan` output before applying changes.

5. **Audit Changes**: Use version control and code review for all Terraform configuration changes.

## Known Security Considerations

- `client_secret` values in `ory_oauth2_client` are stored in Terraform state
- `password` values in `ory_identity` are stored in Terraform state
- SMTP connection URIs may contain credentials
- API keys configured in the provider are passed to the Ory API

Use Terraform's sensitive variable handling and state encryption to protect these values.
