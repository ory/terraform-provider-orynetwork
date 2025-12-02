---
page_title: "ory_organization Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory organization for multi-tenancy.
---

# ory_organization (Resource)

Manages an Ory organization. Organizations enable multi-tenancy for B2B applications, allowing you to group users by their company or team.

~> **Note**: Organizations require an Ory Network plan that supports multi-tenancy.

## Example Usage

### Basic Organization

```terraform
resource "ory_organization" "acme" {
  label = "Acme Corporation"
}
```

### Organization with Email Domains

```terraform
resource "ory_organization" "acme" {
  label   = "Acme Corporation"
  domains = ["acme.com", "acme.io"]
}

output "organization_id" {
  value = ory_organization.acme.id
}
```

### Migrating with Existing ID

When migrating from another identity provider, you can preserve the existing organization ID:

```terraform
resource "ory_organization" "migrated" {
  organization_id = "existing-uuid-from-old-system"
  label          = "Migrated Organization"
}
```

## Schema

### Required

- `label` (String) - Display name for the organization.

### Optional

- `organization_id` (String) - Specify a custom UUID for the organization. Useful for migrations.
- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `domains` (List of String) - Email domains associated with this organization for automatic user assignment.

### Read-Only

- `id` (String) - The organization ID.

## Import

Organizations can be imported using the format `project_id:organization_id`:

```bash
terraform import ory_organization.acme <project-id>:<organization-id>
```
