# Configure the Ory provider
terraform {
  required_providers {
    ory = {
      source = "MaterializeLabs/ory"
    }
  }
}

provider "ory" {
  # Workspace API key for project/organization management
  workspace_api_key = var.ory_workspace_api_key # or set ORY_WORKSPACE_API_KEY env var

  # Project API key for identity/OAuth2 operations
  project_api_key = var.ory_project_api_key # or set ORY_PROJECT_API_KEY env var

  # Project identifiers
  project_id   = var.ory_project_id   # or set ORY_PROJECT_ID env var
  project_slug = var.ory_project_slug # or set ORY_PROJECT_SLUG env var
}

variable "ory_workspace_api_key" {
  type        = string
  sensitive   = true
  description = "Ory Workspace API Key (ory_wak_...)"
}

variable "ory_project_api_key" {
  type        = string
  sensitive   = true
  description = "Ory Project API Key (ory_pat_...)"
}

variable "ory_project_id" {
  type        = string
  description = "Ory Project ID"
}

variable "ory_project_slug" {
  type        = string
  description = "Ory Project Slug (e.g., vibrant-moore-abc123)"
}
