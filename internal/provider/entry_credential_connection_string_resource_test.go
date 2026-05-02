package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialConnectionStringResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccEntryCredentialConnectionStringResourceConfig(
					"tf_test_connection_string", "tf_test_connection_string", "test description", "tf_test_folder",
					"Server=localhost;Database=testdb;User=sa;Password=test123",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_connection_string.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_connection_string.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "name", "tf_test_connection_string"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_connection_string.test", "tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_connection_string.test", "tags.*", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "connection_string", "Server=localhost;Database=testdb;User=sa;Password=test123"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialConnectionStringResourceConfig(
					"tf_test_connection_string", "tf_test_connection_string_updated", "updated description", "tf_test_folder_updated",
					"Server=remotehost;Database=proddb;User=admin;Password=updated456",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "name", "tf_test_connection_string_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "connection_string", "Server=remotehost;Database=proddb;User=admin;Password=updated456"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_connection_string.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_connection_string.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialConnectionStringResourceConfig(vaultName, name, description, folder, connectionString string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder"
}

resource "dvls_entry_folder" "updated" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder_updated"
}

resource "dvls_entry_credential_connection_string" "test" {
  vault_id          = dvls_vault.test.id
  name              = %[3]q
  description       = %[4]q
  folder            = %[5]q
  tags              = ["acceptance", "tf-test"]
  connection_string = %[6]q

  depends_on = [dvls_entry_folder.default, dvls_entry_folder.updated]
}
`, testAccProviderConfig(), vaultName, name, description, folder, connectionString)
}
