---
page_title: "ory_project_api_key Resource - Ory"
subcategory: ""
description: |-
  Manages an Ory project API key.
---

# ory_project_api_key (Resource)

Manages an Ory project API key. Project API keys are used to authenticate API requests to a specific project.

~> **Important**: The API key `value` is only returned when the key is created. Make sure to capture it immediately as it cannot be retrieved later.

## Example Usage

### Basic API Key

```terraform
resource "ory_project_api_key" "backend" {
  name = "Backend Service"
}

output "api_key" {
  value     = ory_project_api_key.backend.value
  sensitive = true
}
```

### API Key with Expiration

```terraform
resource "ory_project_api_key" "temporary" {
  name       = "Temporary Key"
  expires_at = "2024-12-31T23:59:59Z"
}
```

### Multiple API Keys for Different Services

```terraform
resource "ory_project_api_key" "web_app" {
  name = "Web Application"
}

resource "ory_project_api_key" "mobile_app" {
  name = "Mobile Application"
}

resource "ory_project_api_key" "admin_service" {
  name = "Admin Service"
}
```

## Schema

### Required

- `name` (String) - Human-readable name for the API key.

### Optional

- `project_id` (String) - The project ID. If not set, uses the provider's project_id.
- `expires_at` (String) - Expiration timestamp in RFC3339 format (e.g., `2024-12-31T23:59:59Z`). If not set, the key does not expire.

### Read-Only

- `id` (String) - The API key ID.
- `value` (String, Sensitive) - The API key value. Only available on creation.

## Import

Project API keys can be imported using the format `project_id:key_id`:

```bash
terraform import ory_project_api_key.backend <project-id>:<key-id>
```

~> **Note**: The API key `value` cannot be imported and will be empty after import.
