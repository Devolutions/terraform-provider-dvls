package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialUsernamePasswordResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_username_password", "tf_test_folder", "tf_test_folder_updated"),
			// Create
			{
				Config: testAccEntryCredentialUsernamePasswordResourceConfig(
					"tf_test_username_password", "tf_test_username_password", "test description", "tf_test_folder",
					"testuser", "testdomain", "testpassword123",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_username_password.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_username_password.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "name", "tf_test_username_password"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "tags.0", "acceptance"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "tags.1", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "username", "testuser"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "domain", "testdomain"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "password", "testpassword123"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialUsernamePasswordResourceConfig(
					"tf_test_username_password", "tf_test_username_password_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updateddomain", "updatedpassword456",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "name", "tf_test_username_password_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "username", "updateduser"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "domain", "updateddomain"),
					resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "password", "updatedpassword456"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_username_password.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryCredentialImportStateIdFunc("dvls_entry_credential_username_password.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialUsernamePasswordResourceConfig(vaultName, name, description, folder, username, domain, password string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_username_password" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = %[4]q
  folder      = %[5]q
  tags        = ["acceptance", "tf-test"]
  username    = %[6]q
  domain      = %[7]q
  password    = %[8]q
}
`, testAccProviderConfig(), vaultName, name, description, folder, username, domain, password)
}
