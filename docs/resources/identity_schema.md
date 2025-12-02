---
page_title: "ory_identity_schema Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory identity schema.
---

# ory_identity_schema (Resource)

Manages an Ory identity schema. Identity schemas define the structure of user profiles (traits) and how credentials are configured.

~> **Important**: Identity schemas are **immutable** in Ory. Changes to the schema content require creating a new schema. The `name` attribute can be updated, but `schema` changes will trigger resource replacement.

~> **Note**: Schema deletion is not supported by the Ory API. Destroying this resource will remove it from Terraform state but not delete the schema from Ory.

## Example Usage

### Basic Email Schema

```terraform
resource "ory_identity_schema" "basic" {
  name = "basic_user_v1"
  schema = jsonencode({
    "$id"     = "https://example.com/basic.schema.json"
    "$schema" = "http://json-schema.org/draft-07/schema#"
    title     = "Basic User"
    type      = "object"
    properties = {
      traits = {
        type = "object"
        properties = {
          email = {
            type   = "string"
            format = "email"
            title  = "Email"
            "ory.sh/kratos" = {
              credentials = {
                password = { identifier = true }
              }
              verification = { via = "email" }
              recovery     = { via = "email" }
            }
          }
        }
        required = ["email"]
      }
    }
  })
}
```

### Full Customer Schema

```terraform
resource "ory_identity_schema" "customer" {
  name = "customer_v1"
  schema = jsonencode({
    "$id"     = "https://example.com/customer.schema.json"
    "$schema" = "http://json-schema.org/draft-07/schema#"
    title     = "Customer"
    type      = "object"
    properties = {
      traits = {
        type = "object"
        properties = {
          email = {
            type   = "string"
            format = "email"
            title  = "Email Address"
            "ory.sh/kratos" = {
              credentials = {
                password = { identifier = true }
                code     = { identifier = true, via = "email" }
              }
              verification = { via = "email" }
              recovery     = { via = "email" }
            }
          }
          name = {
            type = "object"
            properties = {
              first = {
                type  = "string"
                title = "First Name"
              }
              last = {
                type  = "string"
                title = "Last Name"
              }
            }
          }
          phone = {
            type   = "string"
            title  = "Phone Number"
            "ory.sh/kratos" = {
              credentials = {
                code = { identifier = true, via = "sms" }
              }
            }
          }
        }
        required           = ["email"]
        additionalProperties = false
      }
    }
  })
}
```

## Schema

### Required

- `name` (String) - Name of the identity schema.
- `schema` (String) - JSON-encoded JSON Schema defining identity traits and credentials.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.

### Read-Only

- `id` (String) - The schema ID (hash-based, assigned by Ory).

## Import

Identity schemas can be imported using the format `project_id:schema_id`:

```bash
terraform import ory_identity_schema.customer <project-id>:<schema-id>
```

## JSON Schema Reference

### Ory-Specific Extensions

The `ory.sh/kratos` extension configures how traits are used for authentication:

```json
{
  "ory.sh/kratos": {
    "credentials": {
      "password": { "identifier": true },
      "code": { "identifier": true, "via": "email" },
      "webauthn": { "identifier": true },
      "totp": { "account_name": true }
    },
    "verification": { "via": "email" },
    "recovery": { "via": "email" }
  }
}
```

### Credential Types

| Type | Description |
|------|-------------|
| `password` | Password-based authentication |
| `code` | One-time code via email or SMS |
| `webauthn` | Security keys and passkeys |
| `totp` | Time-based one-time passwords |
| `oidc` | Social/OIDC login |
