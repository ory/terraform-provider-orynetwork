package relationship_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ory/terraform-provider-orynetwork/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ory": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ORY_WORKSPACE_API_KEY"); v == "" {
		t.Skip("ORY_WORKSPACE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ORY_PROJECT_ID"); v == "" {
		t.Skip("ORY_PROJECT_ID must be set for acceptance tests")
	}
	// Relationship tests require project API key and slug for the project API
	if v := os.Getenv("ORY_PROJECT_API_KEY"); v == "" {
		t.Skip("ORY_PROJECT_API_KEY must be set for relationship tests")
	}
	if v := os.Getenv("ORY_PROJECT_SLUG"); v == "" {
		t.Skip("ORY_PROJECT_SLUG must be set for relationship tests")
	}
	// Relationship tests require Ory Keto/Permissions to be enabled and configured
	if os.Getenv("ORY_KETO_TESTS_ENABLED") != "true" {
		t.Skip("ORY_KETO_TESTS_ENABLED must be 'true' (requires Ory Permissions/Keto to be configured)")
	}
}

func TestAccRelationshipResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRelationshipResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_relationship.test", "id"),
					resource.TestCheckResourceAttr("ory_relationship.test", "namespace", "documents"),
					resource.TestCheckResourceAttr("ory_relationship.test", "object", "doc-123"),
					resource.TestCheckResourceAttr("ory_relationship.test", "relation", "viewer"),
					resource.TestCheckResourceAttr("ory_relationship.test", "subject_id", "user-456"),
				),
			},
		},
	})
}

func testAccRelationshipResourceConfig() string {
	return `
provider "ory" {}

resource "ory_relationship" "test" {
  namespace  = "documents"
  object     = "doc-123"
  relation   = "viewer"
  subject_id = "user-456"
}
`
}
