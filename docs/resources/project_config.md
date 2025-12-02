---
page_title: "ory_project_config Resource - Ory"
subcategory: ""
description: |-
  Manages Ory project configuration settings.
---

# ory_project_config (Resource)

Manages Ory project configuration settings including CORS, sessions, authentication methods, password policies, MFA, email delivery, and branding.

## Example Usage

### Basic Configuration

```terraform
resource "ory_project_config" "main" {
  cors_enabled        = true
  cors_origins        = ["https://app.example.com", "https://admin.example.com"]
  password_min_length = 10
  session_lifespan    = "720h"  # 30 days
}
```

### Full Security Configuration

```terraform
resource "ory_project_config" "secure" {
  # CORS
  cors_enabled = true
  cors_origins = ["https://app.example.com"]

  # Sessions
  session_lifespan             = "168h"  # 7 days
  session_cookie_same_site     = "Strict"
  session_cookie_persistent    = true

  # Password Policy
  password_min_length                   = 12
  password_identifier_similarity        = true
  password_haveibeenpwned_enabled       = true
  password_max_breaches                 = 0

  # Authentication Methods
  enable_password = true
  enable_code     = true
  enable_passkey  = true

  # MFA
  enable_totp               = true
  totp_issuer               = "MyApp"
  enable_webauthn           = true
  webauthn_rp_display_name  = "MyApp"
  webauthn_rp_id            = "app.example.com"
  webauthn_rp_origins       = ["https://app.example.com"]
  webauthn_passwordless     = true
  required_aal              = "aal2"  # Require MFA
}
```

### Email/SMTP Configuration

```terraform
resource "ory_project_config" "email" {
  smtp_connection_uri = "smtps://user:pass@smtp.example.com:465"
  smtp_from_address   = "noreply@example.com"
  smtp_from_name      = "MyApp"
  smtp_headers = {
    "X-Custom-Header" = "value"
  }
}
```

### Account Experience Branding

```terraform
resource "ory_project_config" "branding" {
  account_experience_name           = "MyApp"
  account_experience_logo_url       = "https://cdn.example.com/logo.png"
  account_experience_favicon_url    = "https://cdn.example.com/favicon.ico"
  account_experience_default_locale = "en"
  account_experience_stylesheet     = <<-CSS
    :root {
      --primary-color: #0066cc;
      --font-family: 'Inter', sans-serif;
    }
  CSS
}
```

### OAuth2 Token Configuration

```terraform
resource "ory_project_config" "oauth2" {
  oauth2_access_token_ttl   = "1h"
  oauth2_refresh_token_ttl  = "720h"  # 30 days
  oauth2_id_token_ttl       = "1h"
  oauth2_auth_code_ttl      = "10m"
}
```

## Schema

### Optional

#### CORS Settings
- `cors_enabled` (Boolean) - Enable CORS. Default: `false`.
- `cors_origins` (List of String) - Allowed CORS origins.

#### Session Settings
- `session_lifespan` (String) - Session duration (e.g., `720h`).
- `session_cookie_same_site` (String) - SameSite cookie attribute: `Strict`, `Lax`, `None`.
- `session_cookie_persistent` (Boolean) - Use persistent cookies.

#### Password Policy
- `password_min_length` (Number) - Minimum password length.
- `password_identifier_similarity` (Boolean) - Check password similarity to identifier.
- `password_haveibeenpwned_enabled` (Boolean) - Check passwords against HaveIBeenPwned.
- `password_max_breaches` (Number) - Maximum allowed breaches in HaveIBeenPwned.

#### Authentication Methods
- `enable_password` (Boolean) - Enable password authentication.
- `enable_oidc` (Boolean) - Enable OIDC/social authentication.
- `enable_code` (Boolean) - Enable one-time code authentication.
- `enable_passkey` (Boolean) - Enable passkey authentication.
- `enable_totp` (Boolean) - Enable TOTP (authenticator apps).
- `enable_webauthn` (Boolean) - Enable WebAuthn (security keys).
- `enable_lookup_secret` (Boolean) - Enable backup codes.

#### MFA Settings
- `mfa_enforcement` (String) - MFA enforcement policy: `disabled`, `optional`, `required`.
- `totp_issuer` (String) - TOTP issuer name shown in authenticator apps.
- `webauthn_rp_display_name` (String) - WebAuthn relying party display name.
- `webauthn_rp_id` (String) - WebAuthn relying party ID (domain).
- `webauthn_rp_origins` (List of String) - WebAuthn allowed origins.
- `webauthn_passwordless` (Boolean) - Enable passwordless WebAuthn.
- `required_aal` (String) - Required authentication assurance level: `aal1`, `aal2`.
- `session_whoami_required_aal` (String) - Required AAL for session whoami endpoint.

#### Email/SMTP Settings
- `smtp_connection_uri` (String, Sensitive) - SMTP connection URI (e.g., `smtps://user:pass@host:465`).
- `smtp_from_address` (String) - From email address.
- `smtp_from_name` (String) - From display name.
- `smtp_headers` (Map of String) - Custom SMTP headers.

#### Account Experience
- `account_experience_name` (String) - Application name shown in Account Experience.
- `account_experience_logo_url` (String) - Logo URL.
- `account_experience_favicon_url` (String) - Favicon URL.
- `account_experience_default_locale` (String) - Default locale.
- `account_experience_stylesheet` (String) - Custom CSS stylesheet.

#### OAuth2 Token TTLs
- `oauth2_access_token_ttl` (String) - Access token lifetime.
- `oauth2_refresh_token_ttl` (String) - Refresh token lifetime.
- `oauth2_id_token_ttl` (String) - ID token lifetime.
- `oauth2_auth_code_ttl` (String) - Authorization code lifetime.

### Read-Only

- `id` (String) - The project ID.

## Import

Project configuration can be imported using the project ID:

```bash
terraform import ory_project_config.main <project-id>
```
