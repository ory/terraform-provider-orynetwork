---
page_title: "ory_project Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory Network project.
---

# ory_project (Resource)

Manages an Ory Network project. Projects are the top-level container for all Ory resources.

## Example Usage

```terraform
resource "ory_project" "main" {
  name         = "my-production-app"
  workspace_id = var.ory_workspace_id
}
```

## Schema

### Required

- `name` (String) - The name of the project.

### Optional

- `workspace_id` (String) - The workspace ID. If not set, uses the provider's workspace_id.

### Read-Only

- `id` (String) - The project ID.
- `slug` (String) - The project slug (e.g., `vibrant-moore-abc123`).
- `state` (String) - The project state.

## Import

Projects can be imported using their ID:

```bash
terraform import ory_project.main <project-id>
```
