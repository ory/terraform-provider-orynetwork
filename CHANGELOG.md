# Changelog

All notable changes to the Ory Terraform Provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2024-11-29

### Added

#### Resources
- `ory_project` - Manage Ory Network projects
- `ory_workspace` - Manage Ory workspaces
- `ory_organization` - Manage organizations for multi-tenancy (B2B)
- `ory_identity` - Create and manage user identities
- `ory_identity_schema` - Define custom identity schemas
- `ory_oauth2_client` - Manage OAuth2/OIDC client applications
- `ory_project_config` - Configure project settings including:
  - CORS configuration
  - Session settings (lifespan, cookie options)
  - Password policy (min length, haveibeenpwned checks)
  - Authentication methods (password, OIDC, code)
  - MFA settings (TOTP, WebAuthn, Passkeys)
  - SMTP/email configuration
  - Account Experience branding
  - OAuth2 token TTLs
- `ory_action` - Configure webhooks for identity flows (login, registration, recovery, settings, verification)
- `ory_social_provider` - Configure social sign-in providers:
  - Google
  - GitHub
  - Microsoft
  - Apple
  - Generic OIDC
- `ory_email_template` - Customize email templates:
  - Recovery (code and magic link)
  - Verification (code and magic link)
  - Login (code)
  - Registration (code)
- `ory_project_api_key` - Manage project API keys
- `ory_json_web_key_set` - Manage JSON Web Key Sets for token signing
- `ory_relationship` - Manage Ory Permissions (Keto) relationship tuples

#### Provider
- Dual API key support (Workspace API Key + Project API Key)
- Environment variable configuration
- Console API URL override for testing

### Known Limitations
- `ory_identity_schema`: Content is immutable; changes require resource replacement
- `ory_identity_schema`: Delete not supported by Ory API (resource removed from state only)
- `ory_workspace`: Delete not supported by Ory API
- `ory_oauth2_client`: `client_secret` only returned on initial creation
- `ory_email_template`: Delete resets to Ory defaults rather than removing
- `ory_relationship`: Requires Ory Permissions (Keto) to be enabled on the project

[Unreleased]: https://github.com/jasonhernandez/terraform-provider-orynetwork/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/jasonhernandez/terraform-provider-orynetwork/releases/tag/v0.1.0
