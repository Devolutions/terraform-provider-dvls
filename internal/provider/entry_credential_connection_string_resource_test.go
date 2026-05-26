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
			{
				Config: testAccEntryCredentialConnectionStringResourceConfig(
					"tf_test_connection_string", "tf_test_connection_string", "test description", "tf_test_folder",
					"Server=localhost;Database=test;Trusted_Connection=True;",
				),
				Check: testAccEntryCredentialConnectionStringResourceCheck(
					"tf_test_connection_string", "test description", "tf_test_folder",
					"Server=localhost;Database=test;Trusted_Connection=True;",
				),
			},
			{
				Config: testAccEntryCredentialConnectionStringResourceConfig(
					"tf_test_connection_string", "tf_test_connection_string_updated", "updated description", "tf_test_folder_updated",
					"Server=remote;Database=prod;Trusted_Connection=True;",
				),
				Check: testAccEntryCredentialConnectionStringResourceCheck(
					"tf_test_connection_string_updated", "updated description", "tf_test_folder_updated",
					"Server=remote;Database=prod;Trusted_Connection=True;",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_connection_string.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_connection_string.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialConnectionStringResourceCheck(name, description, folder, connectionString string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_connection_string.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_connection_string.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_connection_string.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_connection_string.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_connection_string.test", "connection_string", connectionString),
	)
}

func testAccEntryCredentialConnectionStringResourceConfig(vaultName, name, description, folder, connectionString string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_connection_string",
		vaultName, name, description, folder,
		fmt.Sprintf(`  connection_string = %q`, connectionString),
	)
}
