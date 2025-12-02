---
page_title: "ory_oauth2_client Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory OAuth2/OIDC client.
---

# ory_oauth2_client (Resource)

Manages an Ory OAuth2/OIDC client. OAuth2 clients are used for machine-to-machine authentication, single sign-on, and API access.

~> **Important**: The `client_secret` is only returned when the client is created. Make sure to capture it immediately as it cannot be retrieved later.

## Example Usage

### Machine-to-Machine (Client Credentials)

```terraform
resource "ory_oauth2_client" "api_service" {
  client_name   = "API Service"
  grant_types   = ["client_credentials"]
  token_endpoint_auth_method = "client_secret_post"
  scope         = "read write admin"
}

output "client_id" {
  value = ory_oauth2_client.api_service.client_id
}

output "client_secret" {
  value     = ory_oauth2_client.api_service.client_secret
  sensitive = true
}
```

### Web Application (Authorization Code)

```terraform
resource "ory_oauth2_client" "web_app" {
  client_name   = "My Web Application"
  grant_types   = ["authorization_code", "refresh_token"]
  response_types = ["code"]
  redirect_uris = [
    "https://app.example.com/callback",
    "https://app.example.com/auth/callback"
  ]
  post_logout_redirect_uris = ["https://app.example.com/logout"]
  token_endpoint_auth_method = "client_secret_basic"
  scope = "openid profile email offline_access"
}
```

### Single Page Application (PKCE)

```terraform
resource "ory_oauth2_client" "spa" {
  client_name   = "Single Page App"
  grant_types   = ["authorization_code", "refresh_token"]
  response_types = ["code"]
  redirect_uris = [
    "https://spa.example.com/callback",
    "http://localhost:3000/callback"  # For development
  ]
  token_endpoint_auth_method = "none"  # Public client
  scope = "openid profile email"
}
```

### Mobile Application

```terraform
resource "ory_oauth2_client" "mobile" {
  client_name   = "Mobile App"
  grant_types   = ["authorization_code", "refresh_token"]
  response_types = ["code"]
  redirect_uris = [
    "com.example.app://callback",
    "https://app.example.com/.well-known/apple-app-site-association"
  ]
  token_endpoint_auth_method = "none"
  scope = "openid profile email offline_access"
}
```

## Schema

### Required

- `client_name` (String) - Human-readable name for the client.
- `grant_types` (List of String) - OAuth2 grant types. Common values: `authorization_code`, `client_credentials`, `refresh_token`, `implicit`.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `response_types` (List of String) - OAuth2 response types. Values: `code`, `token`, `id_token`.
- `redirect_uris` (List of String) - Allowed redirect URIs after authentication.
- `post_logout_redirect_uris` (List of String) - Allowed redirect URIs after logout.
- `token_endpoint_auth_method` (String) - Authentication method for token endpoint. Values: `client_secret_basic`, `client_secret_post`, `private_key_jwt`, `none`.
- `scope` (String) - Space-separated list of allowed scopes.
- `audience` (List of String) - Allowed audiences for tokens.
- `owner` (String) - Owner identifier for the client.
- `metadata` (String) - JSON-encoded custom metadata.

### Read-Only

- `id` (String) - The client ID (same as `client_id`).
- `client_id` (String) - The OAuth2 client ID.
- `client_secret` (String, Sensitive) - The OAuth2 client secret. Only available on create.

## Import

OAuth2 clients can be imported using the format `project_id:client_id`:

```bash
terraform import ory_oauth2_client.api_service <project-id>:<client-id>
```

~> **Note**: The `client_secret` cannot be imported and will be empty after import. You may need to regenerate the client if you need the secret.
