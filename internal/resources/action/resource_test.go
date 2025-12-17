package action_test

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
}

func TestAccActionResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_action.test", "id"),
					resource.TestCheckResourceAttr("ory_action.test", "flow", "registration"),
					resource.TestCheckResourceAttr("ory_action.test", "timing", "after"),
					resource.TestCheckResourceAttr("ory_action.test", "auth_method", "password"),
					resource.TestCheckResourceAttr("ory_action.test", "method", "POST"),
				),
			},
		},
	})
}

func testAccActionResourceConfig() string {
	return `
provider "ory" {}

resource "ory_action" "test" {
  flow        = "registration"
  timing      = "after"
  auth_method = "password"
  url         = "https://webhook.example.com/user-registered"
  method      = "POST"
}
`
}
