---
page_title: "ory_project Data Source - Ory"
subcategory: ""
description: |-
  Fetches information about an Ory project.
---

# ory_project (Data Source)

Fetches information about an Ory project. Use this data source to read project details without managing the project lifecycle.

## Example Usage

### Read Current Project

```terraform
# Read the project configured in the provider
data "ory_project" "current" {}

output "project_name" {
  value = data.ory_project.current.name
}

output "project_slug" {
  value = data.ory_project.current.slug
}
```

### Read Specific Project

```terraform
data "ory_project" "other" {
  id = "other-project-uuid"
}
```

## Schema

### Optional

- `id` (String) - Project ID to look up. If not specified, uses the provider's project_id.

### Read-Only

- `name` (String) - The project name.
- `slug` (String) - The project slug.
- `state` (String) - The project state.
- `workspace_id` (String) - The workspace ID the project belongs to.
