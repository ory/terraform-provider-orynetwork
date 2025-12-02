---
page_title: "Ory Provider"
subcategory: ""
description: |-
  Terraform provider for managing Ory Network resources.
---

# Ory Provider

The Ory provider enables Terraform to manage [Ory Network](https://www.ory.sh/) resources for identity management, authentication, and authorization.

> **Note**: This provider is for **Ory Network** (the managed SaaS offering) only. It does not support self-hosted Ory deployments.

## Example Usage

```terraform
terraform {
  required_providers {
    ory = {
      source  = "jasonhernandez/orynetwork"
      version = "~> 0.1"
    }
  }
}

# Configure the provider using environment variables
provider "ory" {}

# Or configure explicitly
provider "ory" {
  workspace_api_key = var.ory_workspace_key
  project_api_key   = var.ory_project_key
  project_id        = var.ory_project_id
  project_slug      = var.ory_project_slug
}
```

## Authentication

Ory Network uses two types of API keys:

| Key Type | Prefix | Purpose |
|----------|--------|---------|
| Workspace API Key | `ory_wak_...` | Projects, organizations, workspace management |
| Project API Key | `ory_pat_...` | Identities, OAuth2 clients, relationships |

### Creating API Keys

1. **Workspace API Key**: Go to [Ory Console](https://console.ory.sh/) → Settings → API Keys
2. **Project API Key**: Go to your Project → Settings → API Keys

### Environment Variables

The recommended approach is to use environment variables:

```bash
export ORY_WORKSPACE_API_KEY="ory_wak_..."
export ORY_PROJECT_API_KEY="ory_pat_..."
export ORY_PROJECT_ID="your-project-uuid"
export ORY_PROJECT_SLUG="your-project-slug"
```

## Schema

### Optional

- `workspace_api_key` (String, Sensitive) - Ory Workspace API Key (`ory_wak_...`). Used for organization and project management. Can also be set via `ORY_WORKSPACE_API_KEY` environment variable.
- `project_api_key` (String, Sensitive) - Ory Project API Key (`ory_pat_...`). Used for identity and OAuth2 operations. Can also be set via `ORY_PROJECT_API_KEY` environment variable.
- `project_id` (String) - Ory Project ID. Can also be set via `ORY_PROJECT_ID` environment variable.
- `project_slug` (String) - Ory Project Slug (e.g., `vibrant-moore-abc123`). Required for identity and OAuth2 operations. Can also be set via `ORY_PROJECT_SLUG` environment variable.
- `workspace_id` (String) - Ory Workspace ID. Can also be set via `ORY_WORKSPACE_ID` environment variable.
- `console_api_url` (String) - Override the console API URL (default: `https://api.console.ory.sh`). Mainly for testing.

~> **Note**: At least one of `workspace_api_key` or `project_api_key` must be configured.
