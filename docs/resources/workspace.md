---
page_title: "ory_workspace Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory Network workspace.
---

# ory_workspace (Resource)

Manages an Ory Network workspace. Workspaces are containers for projects and provide billing isolation.

~> **Note**: Workspace deletion is not supported by the Ory API. Destroying this resource will remove it from Terraform state but not delete the workspace.

## Example Usage

```terraform
resource "ory_workspace" "team" {
  name = "Engineering Team"
}
```

## Schema

### Required

- `name` (String) - The name of the workspace.

### Read-Only

- `id` (String) - The workspace ID.

## Import

Workspaces can be imported using their ID:

```bash
terraform import ory_workspace.team <workspace-id>
```
