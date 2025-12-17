package workspace_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ory/terraform-provider-orynetwork/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ory": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ORY_WORKSPACE_API_KEY"); v == "" {
		t.Skip("ORY_WORKSPACE_API_KEY must be set for workspace acceptance tests")
	}
	// Workspace creation may have quotas, and deletion is NOT supported by Ory API.
	// Created workspaces will remain in your Ory account!
	if os.Getenv("ORY_WORKSPACE_TESTS_ENABLED") != "true" {
		t.Skip("ORY_WORKSPACE_TESTS_ENABLED must be 'true' to run workspace tests (creates real workspaces that CANNOT be deleted)")
	}
}

func testAccPreCheckReadOnly(t *testing.T) {
	if v := os.Getenv("ORY_WORKSPACE_API_KEY"); v == "" {
		t.Skip("ORY_WORKSPACE_API_KEY must be set for workspace acceptance tests")
	}
	if v := os.Getenv("ORY_WORKSPACE_ID"); v == "" {
		t.Skip("ORY_WORKSPACE_ID must be set for read-only workspace tests")
	}
}

// TestAccWorkspaceResource_basic tests the CRUD lifecycle of a workspace.
// WARNING: Ory API does NOT support workspace deletion!
// Workspaces created by this test will remain in your Ory account.
// Only run this test if you understand the implications.
func TestAccWorkspaceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccWorkspaceResourceConfig("tf-test-workspace"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_workspace.test", "id"),
					resource.TestCheckResourceAttr("ory_workspace.test", "name", "tf-test-workspace"),
					resource.TestCheckResourceAttrSet("ory_workspace.test", "created_at"),
					resource.TestCheckResourceAttrSet("ory_workspace.test", "updated_at"),
				),
			},
			// ImportState
			{
				ResourceName:      "ory_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccWorkspaceResourceConfig("tf-test-workspace-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_workspace.test", "id"),
					resource.TestCheckResourceAttr("ory_workspace.test", "name", "tf-test-workspace-updated"),
				),
			},
		},
	})
}

// TestAccWorkspaceDataSource_existing tests reading an existing workspace.
// This test uses an existing workspace ID and doesn't create new resources.
func TestAccWorkspaceResource_import(t *testing.T) {
	workspaceID := os.Getenv("ORY_WORKSPACE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckReadOnly(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import existing workspace
			{
				Config:        testAccWorkspaceResourceConfigEmpty(),
				ImportState:   true,
				ImportStateId: workspaceID,
				ResourceName:  "ory_workspace.test",
				// After import, we can verify basic attributes
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					state := states[0]
					if state.ID == "" {
						return fmt.Errorf("expected non-empty ID")
					}
					if state.Attributes["name"] == "" {
						return fmt.Errorf("expected non-empty name")
					}
					return nil
				},
			},
		},
	})
}

func testAccWorkspaceResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_workspace" "test" {
  name = %[1]q
}
`, name)
}

func testAccWorkspaceResourceConfigEmpty() string {
	return `
provider "ory" {}

resource "ory_workspace" "test" {
  name = "placeholder"
}
`
}
