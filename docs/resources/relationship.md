---
page_title: "ory_relationship Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory Permissions (Keto) relationship tuple.
---

# ory_relationship (Resource)

Manages an Ory Permissions (Keto) relationship tuple. Relationships define fine-grained permissions using the ReBAC (Relationship-Based Access Control) model.

~> **Note**: This resource requires Ory Permissions (Keto) to be enabled on your project.

## Example Usage

### Direct User Permission

```terraform
# User 'user-123' is a viewer of document 'doc-456'
resource "ory_relationship" "user_view_doc" {
  namespace  = "documents"
  object     = "doc-456"
  relation   = "viewer"
  subject_id = "user-123"
}
```

### Group Membership

```terraform
# User 'user-123' is a member of 'engineering' group
resource "ory_relationship" "user_in_group" {
  namespace  = "groups"
  object     = "engineering"
  relation   = "member"
  subject_id = "user-123"
}
```

### Subject Set (Group Permission)

```terraform
# Members of 'engineering' group are editors of 'project-abc'
resource "ory_relationship" "group_edit_project" {
  namespace            = "projects"
  object               = "project-abc"
  relation             = "editor"
  subject_set_namespace = "groups"
  subject_set_object    = "engineering"
  subject_set_relation  = "member"
}
```

### Hierarchical Permissions

```terraform
# Parent folder grants access to child documents
resource "ory_relationship" "folder_contains_doc" {
  namespace  = "documents"
  object     = "doc-789"
  relation   = "parent"
  subject_id = "folder-123"
}

# Viewers of parent folder can view child documents
resource "ory_relationship" "folder_viewer_inheritance" {
  namespace            = "documents"
  object               = "doc-789"
  relation             = "viewer"
  subject_set_namespace = "folders"
  subject_set_object    = "folder-123"
  subject_set_relation  = "viewer"
}
```

### Complete RBAC Example

```terraform
# Define namespace for your resources
locals {
  users = {
    alice = "user-alice"
    bob   = "user-bob"
  }

  docs = {
    readme  = "readme.md"
    secrets = "secrets.md"
  }
}

# Alice is owner of readme
resource "ory_relationship" "alice_owns_readme" {
  namespace  = "documents"
  object     = local.docs.readme
  relation   = "owner"
  subject_id = local.users.alice
}

# Bob can view readme
resource "ory_relationship" "bob_views_readme" {
  namespace  = "documents"
  object     = local.docs.readme
  relation   = "viewer"
  subject_id = local.users.bob
}

# Only Alice can access secrets
resource "ory_relationship" "alice_owns_secrets" {
  namespace  = "documents"
  object     = local.docs.secrets
  relation   = "owner"
  subject_id = local.users.alice
}
```

## Schema

### Required

- `namespace` (String) - The namespace of the relationship.
- `object` (String) - The object (resource) identifier.
- `relation` (String) - The relation (permission type).

### Optional (Subject - One Required)

Either `subject_id` OR all three `subject_set_*` attributes must be provided:

- `subject_id` (String) - Direct subject identifier (e.g., user ID).
- `subject_set_namespace` (String) - Subject set namespace.
- `subject_set_object` (String) - Subject set object.
- `subject_set_relation` (String) - Subject set relation.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.

### Read-Only

- `id` (String) - Resource identifier.

## Import

Relationships can be imported using the format:

For direct subjects:
```bash
terraform import ory_relationship.user_view_doc <namespace>:<object>:<relation>:<subject_id>
```

For subject sets:
```bash
terraform import ory_relationship.group_edit_project <namespace>:<object>:<relation>:<subject_namespace>:<subject_object>:<subject_relation>
```

## Relationship Model

Ory Permissions uses a tuple-based model:

```
namespace:object#relation@subject
```

For example:
- `documents:doc-123#viewer@user-456` - User 456 can view document 123
- `documents:doc-123#viewer@groups:engineering#member` - Engineering group members can view document 123

### Subject Types

| Type | Example | Description |
|------|---------|-------------|
| Direct | `user-123` | Direct reference to a subject |
| Subject Set | `groups:engineering#member` | All subjects with the specified relation |

### Common Patterns

| Pattern | Description |
|---------|-------------|
| User → Resource | Direct permission grant |
| User → Group → Resource | Group-based permissions |
| Resource → Resource | Hierarchical/nested permissions |
