# Create an organization for multi-tenancy
resource "ory_organization" "acme" {
  label   = "Acme Corporation"
  domains = ["acme.com", "acme.io"]
}

output "organization_id" {
  value = ory_organization.acme.id
}
