# Read the current project (from provider config)
data "ory_project" "current" {}

output "project_name" {
  value = data.ory_project.current.name
}

output "project_slug" {
  value = data.ory_project.current.slug
}

output "project_state" {
  value = data.ory_project.current.state
}

# Read a specific project by ID
data "ory_project" "other" {
  id = "other-project-uuid"
}
