package projectconfig_test

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

func TestAccProjectConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfigResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project_config.test", "id"),
					resource.TestCheckResourceAttr("ory_project_config.test", "cors_enabled", "true"),
					resource.TestCheckResourceAttr("ory_project_config.test", "password_min_length", "10"),
				),
			},
		},
	})
}

func testAccProjectConfigResourceConfig() string {
	return `
provider "ory" {}

resource "ory_project_config" "test" {
  cors_enabled        = true
  cors_origins        = ["https://example.com"]
  password_min_length = 10
}
`
}

func TestAccProjectConfigResource_mfaPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfigResourceMFAConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project_config.test", "id"),
					resource.TestCheckResourceAttr("ory_project_config.test", "enable_totp", "true"),
					resource.TestCheckResourceAttr("ory_project_config.test", "totp_issuer", "TerraformTest"),
				),
			},
		},
	})
}

func testAccProjectConfigResourceMFAConfig() string {
	return `
provider "ory" {}

resource "ory_project_config" "test" {
  enable_totp  = true
  totp_issuer  = "TerraformTest"
}
`
}

func TestAccProjectConfigResource_accountExperience(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfigResourceAccountExperienceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project_config.test", "id"),
					resource.TestCheckResourceAttr("ory_project_config.test", "account_experience_name", "TF Test App"),
					resource.TestCheckResourceAttr("ory_project_config.test", "account_experience_default_locale", "en"),
				),
			},
		},
	})
}

func testAccProjectConfigResourceAccountExperienceConfig() string {
	return `
provider "ory" {}

resource "ory_project_config" "test" {
  account_experience_name           = "TF Test App"
  account_experience_default_locale = "en"
}
`
}
