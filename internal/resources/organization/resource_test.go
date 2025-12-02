// Copyright 2025 Materialize Inc. and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package organization_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/jasonhernandez/terraform-provider-orynetwork/internal/provider"
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
	// Organization operations require workspace API key
	if v := os.Getenv("ORY_WORKSPACE_API_KEY"); v == "" {
		t.Skip("ORY_WORKSPACE_API_KEY must be set for organization acceptance tests")
	}
}

func testAccPreCheckB2B(t *testing.T) {
	testAccPreCheck(t)
	// Organizations require B2B features which are only available on paid plans
	// Set ORY_B2B_ENABLED=true if your plan supports organizations
	if os.Getenv("ORY_B2B_ENABLED") != "true" {
		t.Skip("ORY_B2B_ENABLED must be set to 'true' for organization tests (requires paid Ory plan)")
	}
}

func TestAccOrganizationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckB2B(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOrganizationResourceConfig("Test Organization"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_organization.test", "id"),
					resource.TestCheckResourceAttr("ory_organization.test", "label", "Test Organization"),
					resource.TestCheckResourceAttrSet("ory_organization.test", "created_at"),
				),
			},
			// ImportState
			{
				ResourceName:      "ory_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccOrganizationResourceConfig("Test Organization Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_organization.test", "id"),
					resource.TestCheckResourceAttr("ory_organization.test", "label", "Test Organization Updated"),
				),
			},
		},
	})
}

func TestAccOrganizationResource_withDomains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckB2B(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfigWithDomains("Org with Domains", []string{"example.com", "test.example.com"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_organization.test", "id"),
					resource.TestCheckResourceAttr("ory_organization.test", "label", "Org with Domains"),
					resource.TestCheckResourceAttr("ory_organization.test", "domains.#", "2"),
				),
			},
			// Update domains
			{
				Config: testAccOrganizationResourceConfigWithDomains("Org with Domains", []string{"example.com", "test.example.com", "new.example.com"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ory_organization.test", "domains.#", "3"),
				),
			},
		},
	})
}

func testAccOrganizationResourceConfig(label string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_organization" "test" {
  label = %[1]q
}
`, label)
}

func testAccOrganizationResourceConfigWithDomains(label string, domains []string) string {
	// Build domain list for HCL
	domainList := ""
	for i, d := range domains {
		if i > 0 {
			domainList += ", "
		}
		domainList += fmt.Sprintf("%q", d)
	}

	return fmt.Sprintf(`
provider "ory" {}

resource "ory_organization" "test" {
  label   = %[1]q
  domains = [%[2]s]
}
`, label, domainList)
}
