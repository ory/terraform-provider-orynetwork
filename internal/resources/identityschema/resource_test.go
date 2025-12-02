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

package identityschema_test

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
	// Identity schema tests can leave schemas behind (deletion not supported)
	if os.Getenv("ORY_SCHEMA_TESTS_ENABLED") != "true" {
		t.Skip("ORY_SCHEMA_TESTS_ENABLED must be 'true' (schemas cannot be deleted from Ory)")
	}
}

func TestAccIdentitySchemaResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentitySchemaResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_identity_schema.test", "id"),
					resource.TestCheckResourceAttr("ory_identity_schema.test", "schema_id", "tf-test-schema"),
				),
			},
		},
	})
}

func testAccIdentitySchemaResourceConfig() string {
	return `
provider "ory" {}

resource "ory_identity_schema" "test" {
  schema_id = "tf-test-schema"
  schema    = jsonencode({
    "$id": "https://example.com/tf-test-schema.json",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Test Schema",
    "type": "object",
    "properties": {
      "traits": {
        "type": "object",
        "properties": {
          "email": {
            "type": "string",
            "format": "email",
            "title": "Email",
            "ory.sh/kratos": {
              "credentials": {
                "password": {"identifier": true}
              },
              "verification": {"via": "email"},
              "recovery": {"via": "email"}
            }
          }
        },
        "required": ["email"]
      }
    }
  })
}
`
}
