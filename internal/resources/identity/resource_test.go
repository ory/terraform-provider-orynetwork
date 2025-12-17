package identity_test

import (
	"fmt"
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
	if v := os.Getenv("ORY_PROJECT_API_KEY"); v == "" {
		t.Skip("ORY_PROJECT_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ORY_PROJECT_ID"); v == "" {
		t.Skip("ORY_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ORY_PROJECT_SLUG"); v == "" {
		t.Skip("ORY_PROJECT_SLUG must be set for acceptance tests")
	}
}

func TestAccIdentityResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIdentityResourceConfig("test-basic@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_identity.test", "id"),
					resource.TestCheckResourceAttr("ory_identity.test", "schema_id", "preset://email"),
					resource.TestCheckResourceAttr("ory_identity.test", "state", "active"),
				),
			},
			// ImportState
			{
				ResourceName:      "ory_identity.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Password is write-only, not returned on read
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update
			{
				Config: testAccIdentityResourceConfig("test-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_identity.test", "id"),
					resource.TestCheckResourceAttr("ory_identity.test", "state", "active"),
				),
			},
		},
	})
}

func TestAccIdentityResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityResourceConfigWithMetadata("test-metadata@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_identity.test", "id"),
					resource.TestCheckResourceAttr("ory_identity.test", "state", "active"),
					resource.TestCheckResourceAttrSet("ory_identity.test", "metadata_public"),
				),
			},
		},
	})
}

func TestAccIdentityResource_inactive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityResourceConfigInactive("test-inactive@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_identity.test", "id"),
					resource.TestCheckResourceAttr("ory_identity.test", "state", "inactive"),
				),
			},
		},
	})
}

func testAccIdentityResourceConfig(email string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_identity" "test" {
  schema_id = "preset://email"

  traits = jsonencode({
    email = %[1]q
  })

  state = "active"
}
`, email)
}

func testAccIdentityResourceConfigWithMetadata(email string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_identity" "test" {
  schema_id = "preset://email"

  traits = jsonencode({
    email = %[1]q
  })

  state = "active"

  metadata_public = jsonencode({
    role        = "admin"
    created_by  = "terraform"
  })
}
`, email)
}

func testAccIdentityResourceConfigInactive(email string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_identity" "test" {
  schema_id = "preset://email"

  traits = jsonencode({
    email = %[1]q
  })

  state = "inactive"
}
`, email)
}
