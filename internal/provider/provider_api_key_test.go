package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProvider_apiKeyFromEnv(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiKey(t)
			t.Setenv("DVLS_API_KEY", os.Getenv("TEST_DVLS_API_KEY"))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVaultResourceConfig_minimal("tf_test_vault_api_key_env"),
				Check:  resource.TestCheckResourceAttrSet("dvls_vault.test", "id"),
			},
		},
	})
}

func TestAccProvider_apiKeyAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckApiKey(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "dvls" {
  base_uri = %q
  api_key  = %q
}

resource "dvls_vault" "test" {
  name = "tf_test_vault_api_key_attribute"
}
`, os.Getenv("TEST_DVLS_BASE_URI"), os.Getenv("TEST_DVLS_API_KEY")),
				Check: resource.TestCheckResourceAttrSet("dvls_vault.test", "id"),
			},
		},
	})
}

func testAccPreCheckApiKey(t *testing.T) {
	t.Helper()

	testAccPreCheck(t)

	if os.Getenv("TEST_DVLS_API_KEY") == "" {
		t.Skip("TEST_DVLS_API_KEY must be set for api_key acceptance tests")
	}

	t.Setenv("DVLS_APP_ID", "")
	t.Setenv("DVLS_APP_SECRET", "")
}
