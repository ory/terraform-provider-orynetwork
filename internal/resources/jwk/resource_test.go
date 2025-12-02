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

package jwk_test

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
	// JWK tests require project API key and slug for the project API
	if v := os.Getenv("ORY_PROJECT_API_KEY"); v == "" {
		t.Skip("ORY_PROJECT_API_KEY must be set for JWK tests")
	}
	if v := os.Getenv("ORY_PROJECT_SLUG"); v == "" {
		t.Skip("ORY_PROJECT_SLUG must be set for JWK tests")
	}
}

func TestAccJWKResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJWKResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_json_web_key_set.test", "id"),
					resource.TestCheckResourceAttr("ory_json_web_key_set.test", "set_id", "tf-test-jwks"),
					resource.TestCheckResourceAttr("ory_json_web_key_set.test", "key_id", "tf-test-key"),
					resource.TestCheckResourceAttr("ory_json_web_key_set.test", "algorithm", "RS256"),
					resource.TestCheckResourceAttr("ory_json_web_key_set.test", "use", "sig"),
					resource.TestCheckResourceAttrSet("ory_json_web_key_set.test", "keys"),
				),
			},
		},
	})
}

func testAccJWKResourceConfig() string {
	return `
provider "ory" {}

resource "ory_json_web_key_set" "test" {
  set_id    = "tf-test-jwks"
  key_id    = "tf-test-key"
  algorithm = "RS256"
  use       = "sig"
}
`
}
