---
page_title: "ory_identity Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory identity (user).
---

# ory_identity (Resource)

Manages an Ory identity. Identities represent users in your application and contain traits (profile data) and credentials.

## Example Usage

### Basic Identity

```terraform
resource "ory_identity" "user" {
  schema_id = "preset://email"
  traits = jsonencode({
    email = "user@example.com"
  })
}
```

### Identity with Password

```terraform
resource "ory_identity" "user" {
  schema_id = "preset://email"
  traits = jsonencode({
    email = "user@example.com"
  })
  password = var.user_password  # Should be marked sensitive
  state    = "active"
}
```

### Identity with Custom Schema

```terraform
resource "ory_identity" "customer" {
  schema_id = ory_identity_schema.customer.id
  traits = jsonencode({
    email = "customer@example.com"
    name = {
      first = "John"
      last  = "Doe"
    }
    company = "Acme Inc"
  })
  metadata_public = jsonencode({
    tier = "enterprise"
  })
  metadata_admin = jsonencode({
    internal_id = "cust-12345"
  })
}
```

### Identity in Organization

```terraform
resource "ory_identity" "org_user" {
  schema_id       = "preset://email"
  organization_id = ory_organization.acme.id
  traits = jsonencode({
    email = "user@acme.com"
  })
}
```

## Schema

### Required

- `schema_id` (String) - The identity schema ID. Use `preset://email` for the default email schema.
- `traits` (String) - JSON-encoded identity traits. Must conform to the identity schema.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `password` (String, Sensitive) - The user's password. Only used for password-based authentication.
- `state` (String) - Identity state: `active` or `inactive`. Default: `active`.
- `organization_id` (String) - Organization ID for multi-tenant setups.
- `metadata_public` (String) - JSON-encoded public metadata (visible to the user).
- `metadata_admin` (String) - JSON-encoded admin metadata (only visible to admins).

### Read-Only

- `id` (String) - The identity ID.

## Import

Identities can be imported using the format `project_id:identity_id`:

```bash
terraform import ory_identity.user <project-id>:<identity-id>
```

~> **Note**: The password cannot be imported and will need to be re-set after import.
