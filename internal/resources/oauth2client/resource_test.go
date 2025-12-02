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

package oauth2client_test

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
}

func TestAccOAuth2ClientResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccOAuth2ClientResourceConfig("Test API Client"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_oauth2_client.test", "id"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "client_name", "Test API Client"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "scope", "api:read"),
					// client_secret is only returned on create
					resource.TestCheckResourceAttrSet("ory_oauth2_client.test", "client_secret"),
				),
			},
			// ImportState
			{
				ResourceName:      "ory_oauth2_client.test",
				ImportState:       true,
				ImportStateVerify: true,
				// client_secret is only returned on create, not on read
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
			// Update
			{
				Config: testAccOAuth2ClientResourceConfigUpdated("Test API Client Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_oauth2_client.test", "id"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "client_name", "Test API Client Updated"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "scope", "api:read api:write"),
				),
			},
		},
	})
}

func TestAccOAuth2ClientResource_withAudience(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOAuth2ClientResourceConfigWithAudience("Test Client with Audience"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_oauth2_client.test", "id"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "client_name", "Test Client with Audience"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "audience.#", "2"),
				),
			},
		},
	})
}

func TestAccOAuth2ClientResource_withRedirectURIs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOAuth2ClientResourceConfigWithRedirectURIs("Test Client with Redirects"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ory_oauth2_client.test", "id"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "client_name", "Test Client with Redirects"),
					resource.TestCheckResourceAttr("ory_oauth2_client.test", "redirect_uris.#", "2"),
				),
			},
		},
	})
}

func testAccOAuth2ClientResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_oauth2_client" "test" {
  client_name = %[1]q

  grant_types    = ["client_credentials"]
  response_types = ["token"]
  scope          = "api:read"
}
`, name)
}

func testAccOAuth2ClientResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_oauth2_client" "test" {
  client_name = %[1]q

  grant_types    = ["client_credentials"]
  response_types = ["token"]
  scope          = "api:read api:write"
}
`, name)
}

func testAccOAuth2ClientResourceConfigWithAudience(name string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_oauth2_client" "test" {
  client_name = %[1]q

  grant_types    = ["client_credentials"]
  response_types = ["token"]
  scope          = "api:read"
  audience       = ["https://api.example.com", "https://api2.example.com"]
}
`, name)
}

func testAccOAuth2ClientResourceConfigWithRedirectURIs(name string) string {
	return fmt.Sprintf(`
provider "ory" {}

resource "ory_oauth2_client" "test" {
  client_name = %[1]q

  grant_types    = ["authorization_code", "refresh_token"]
  response_types = ["code"]
  scope          = "openid profile email"
  redirect_uris  = ["https://app.example.com/callback", "http://localhost:3000/callback"]
}
`, name)
}
