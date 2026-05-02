package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialConnectionStringDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialConnectionStringDataSourceConfig_byName("tf_test_connection_string_by_name", "tf_test_connection_string_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "id", "dvls_entry_credential_connection_string.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "vault_id", "dvls_entry_credential_connection_string.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "name", "dvls_entry_credential_connection_string.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "description", "dvls_entry_credential_connection_string.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "folder", "dvls_entry_credential_connection_string.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "tags.#", "dvls_entry_credential_connection_string.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "connection_string", "dvls_entry_credential_connection_string.test", "connection_string"),
				),
			},
		},
	})
}

func TestAccEntryCredentialConnectionStringDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialConnectionStringDataSourceConfig_byId("tf_test_connection_string_by_id", "tf_test_connection_string_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "id", "dvls_entry_credential_connection_string.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "vault_id", "dvls_entry_credential_connection_string.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "name", "dvls_entry_credential_connection_string.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "description", "dvls_entry_credential_connection_string.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "folder", "dvls_entry_credential_connection_string.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "tags.#", "dvls_entry_credential_connection_string.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_connection_string.test", "connection_string", "dvls_entry_credential_connection_string.test", "connection_string"),
				),
			},
		},
	})
}

func testAccEntryCredentialConnectionStringDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder"
}

resource "dvls_entry_credential_connection_string" "test" {
  vault_id          = dvls_vault.test.id
  name              = %[3]q
  description       = "test entry for data source"
  folder            = "tf_test_folder"
  tags              = ["acceptance", "tf-test"]
  connection_string = "Server=localhost;Database=testdb;User=sa;Password=test123"

  depends_on = [dvls_entry_folder.default]
}

data "dvls_entry_credential_connection_string" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_connection_string.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialConnectionStringDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder"
}

resource "dvls_entry_credential_connection_string" "test" {
  vault_id          = dvls_vault.test.id
  name              = %[3]q
  description       = "test entry for data source"
  folder            = "tf_test_folder"
  tags              = ["acceptance", "tf-test"]
  connection_string = "Server=localhost;Database=testdb;User=sa;Password=test123"

  depends_on = [dvls_entry_folder.default]
}

data "dvls_entry_credential_connection_string" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_connection_string.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
