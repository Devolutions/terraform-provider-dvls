package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVaultResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccVaultResourceConfig("tf_test_vault", "test description", "private", "high", "credentials"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_vault.test", "name", "tf_test_vault"),
					resource.TestCheckResourceAttr("dvls_vault.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_vault.test", "visibility", "private"),
					resource.TestCheckResourceAttr("dvls_vault.test", "security_level", "high"),
					resource.TestCheckResourceAttr("dvls_vault.test", "content_type", "credentials"),
				),
			},
			// Update
			{
				Config: testAccVaultResourceConfig("tf_test_vault_updated", "updated description", "public", "high", "credentials"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_vault.test", "name", "tf_test_vault_updated"),
					resource.TestCheckResourceAttr("dvls_vault.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_vault.test", "visibility", "public"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_vault.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVaultResource_minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVaultResourceConfig_minimal("tf_test_vault_minimal"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_vault.test", "name", "tf_test_vault_minimal"),
					resource.TestCheckResourceAttr("dvls_vault.test", "visibility", "default"),
					resource.TestCheckResourceAttr("dvls_vault.test", "security_level", "standard"),
					resource.TestCheckResourceAttr("dvls_vault.test", "content_type", "everything"),
				),
			},
		},
	})
}

func testAccVaultResourceConfig(name, description, visibility, securityLevel, contentType string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name           = %[2]q
  description    = %[3]q
  visibility     = %[4]q
  security_level = %[5]q
  content_type   = %[6]q
}
`, testAccProviderConfig(), name, description, visibility, securityLevel, contentType)
}

func testAccVaultResourceConfig_minimal(name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}
`, testAccProviderConfig(), name)
}
