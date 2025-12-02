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

package project_test

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
	if v := os.Getenv("ORY_WORKSPACE_API_KEY"); v == "" {
		t.Skip("ORY_WORKSPACE_API_KEY must be set for project acceptance tests")
	}
	// Project creation/deletion is expensive and may have quotas
	// Only run if explicitly enabled
	if os.Getenv("ORY_PROJECT_TESTS_ENABLED") != "true" {
		t.Skip("ORY_PROJECT_TESTS_ENABLED must be 'true' to run project tests (creates/deletes real projects)")
	}
}

// TestAccProjectResource_basic tests the full CRUD lifecycle of a project.
// WARNING: This test creates and deletes a real Ory project.
// Only run this test if you have quota available and understand the implications.
func TestAccProjectResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccProjectResourceConfig("tf-test-project", "dev"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project.test", "id"),
					resource.TestCheckResourceAttr("ory_project.test", "name", "tf-test-project"),
					resource.TestCheckResourceAttr("ory_project.test", "environment", "dev"),
					resource.TestCheckResourceAttrSet("ory_project.test", "slug"),
					resource.TestCheckResourceAttr("ory_project.test", "state", "running"),
				),
			},
			// ImportState
			{
				ResourceName:      "ory_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccProjectResource_prodEnvironment tests creating a production project.
func TestAccProjectResource_prodEnvironment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig("tf-test-prod-project", "prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project.test", "id"),
					resource.TestCheckResourceAttr("ory_project.test", "environment", "prod"),
				),
			},
		},
	})
}

func testAccProjectResourceConfig(name, environment string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_project" "test" {
  name        = %[1]q
  environment = %[2]q
}
`, name, environment)
}
