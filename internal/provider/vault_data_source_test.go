package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVaultDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVaultDataSourceConfig_byName("tf_test_vault_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "name", "dvls_vault.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "description", "dvls_vault.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "visibility", "dvls_vault.test", "visibility"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "security_level", "dvls_vault.test", "security_level"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "content_type", "dvls_vault.test", "content_type"),
				),
			},
		},
	})
}

func TestAccVaultDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVaultDataSourceConfig_byId("tf_test_vault_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "name", "dvls_vault.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "description", "dvls_vault.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "visibility", "dvls_vault.test", "visibility"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "security_level", "dvls_vault.test", "security_level"),
					resource.TestCheckResourceAttrPair("data.dvls_vault.test", "content_type", "dvls_vault.test", "content_type"),
				),
			},
		},
	})
}

func testAccVaultDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name        = %[2]q
  description = "test vault for data source"
}

data "dvls_vault" "test" {
  name = dvls_vault.test.name
}
`, testAccProviderConfig(), name)
}

func testAccVaultDataSourceConfig_byId(name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name        = %[2]q
  description = "test vault for data source"
}

data "dvls_vault" "test" {
  id = dvls_vault.test.id
}
`, testAccProviderConfig(), name)
}
