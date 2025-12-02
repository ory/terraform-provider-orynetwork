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

package projectapikey_test

import (
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
		t.Skip("ORY_WORKSPACE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ORY_PROJECT_ID"); v == "" {
		t.Skip("ORY_PROJECT_ID must be set for acceptance tests")
	}
}

func TestAccProjectAPIKeyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAPIKeyResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_project_api_key.test", "id"),
					resource.TestCheckResourceAttr("ory_project_api_key.test", "name", "tf-test-key"),
					resource.TestCheckResourceAttrSet("ory_project_api_key.test", "value"),
					resource.TestCheckResourceAttrSet("ory_project_api_key.test", "owner_id"),
				),
			},
		},
	})
}

func testAccProjectAPIKeyResourceConfig() string {
	return `
provider "ory" {}

resource "ory_project_api_key" "test" {
  name = "tf-test-key"
}
`
}
