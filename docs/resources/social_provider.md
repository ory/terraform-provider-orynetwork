---
page_title: "ory_social_provider Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory social sign-in provider (OIDC).
---

# ory_social_provider (Resource)

Manages an Ory social sign-in provider. This resource configures OAuth2/OIDC providers like Google, GitHub, Microsoft, and Apple for social login.

## Example Usage

### Google

```terraform
resource "ory_social_provider" "google" {
  provider_id   = "google"
  client_id     = var.google_client_id
  client_secret = var.google_client_secret
  scopes        = ["email", "profile"]
}
```

### GitHub

```terraform
resource "ory_social_provider" "github" {
  provider_id   = "github"
  client_id     = var.github_client_id
  client_secret = var.github_client_secret
  scopes        = ["user:email", "read:user"]
}
```

### Microsoft (Azure AD)

```terraform
resource "ory_social_provider" "microsoft" {
  provider_id   = "microsoft"
  client_id     = var.azure_client_id
  client_secret = var.azure_client_secret
  tenant        = var.azure_tenant_id  # or "common" for multi-tenant
  scopes        = ["openid", "profile", "email"]
}
```

### Apple

```terraform
resource "ory_social_provider" "apple" {
  provider_id        = "apple"
  client_id          = var.apple_client_id
  apple_team_id      = var.apple_team_id
  apple_private_key  = var.apple_private_key
  apple_private_key_id = var.apple_key_id
  scopes             = ["email", "name"]
}
```

### Generic OIDC Provider

```terraform
resource "ory_social_provider" "corporate_sso" {
  provider_id      = "generic"
  provider_label   = "Corporate SSO"
  client_id        = var.sso_client_id
  client_secret    = var.sso_client_secret
  issuer_url       = "https://sso.example.com"
  scopes           = ["openid", "profile", "email"]
}
```

### Custom Mapper

```terraform
resource "ory_social_provider" "custom" {
  provider_id   = "google"
  client_id     = var.google_client_id
  client_secret = var.google_client_secret
  scopes        = ["email", "profile"]

  mapper = <<-JSONNET
    local claims = std.extVar('claims');
    {
      identity: {
        traits: {
          email: claims.email,
          name: {
            first: claims.given_name,
            last: claims.family_name
          },
          picture: claims.picture
        }
      }
    }
  JSONNET
}
```

## Schema

### Required

- `provider_id` (String) - Provider type. Values: `google`, `github`, `microsoft`, `apple`, `discord`, `slack`, `spotify`, `generic`.
- `client_id` (String) - OAuth2 client ID from the provider.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `client_secret` (String, Sensitive) - OAuth2 client secret. Required for most providers.
- `scopes` (List of String) - OAuth2 scopes to request.
- `provider_label` (String) - Display label for the provider (used in UI).
- `issuer_url` (String) - OIDC issuer URL. Required for `generic` provider.
- `tenant` (String) - Microsoft Azure AD tenant ID. Use `common` for multi-tenant.
- `mapper` (String) - Jsonnet template for mapping provider claims to identity traits.

#### Apple-specific
- `apple_team_id` (String) - Apple Developer Team ID.
- `apple_private_key` (String, Sensitive) - Apple Sign-in private key (PEM format).
- `apple_private_key_id` (String) - Apple private key ID.

### Read-Only

- `id` (String) - Resource identifier in format `project_id:provider_id`.

## Import

Social providers can be imported using the format `project_id:provider_id`:

```bash
terraform import ory_social_provider.google <project-id>:google
```

## Claim Mapping

The `mapper` attribute allows you to customize how provider claims are mapped to identity traits using Jsonnet:

```jsonnet
local claims = std.extVar('claims');
{
  identity: {
    traits: {
      email: claims.email,
      // Map custom fields
      [if "custom_field" in claims then "custom"]: claims.custom_field
    },
    // Optional: Set metadata
    metadata_public: {
      provider: "google",
      picture: claims.picture
    }
  }
}
```

### Available Claims

| Provider | Common Claims |
|----------|--------------|
| Google | `email`, `email_verified`, `name`, `given_name`, `family_name`, `picture`, `locale` |
| GitHub | `email`, `login`, `name`, `avatar_url` |
| Microsoft | `email`, `name`, `given_name`, `family_name`, `preferred_username` |
| Apple | `email`, `email_verified`, `name` (first login only) |
