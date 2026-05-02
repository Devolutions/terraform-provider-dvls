package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSSHKeyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccEntryCredentialSSHKeyResourceConfig(
					"tf_test_ssh_key", "tf_test_ssh_key", "test description", "tf_test_folder",
					"testuser", "testpassword", "testpassphrase", "test-private-key-data", "test-public-key",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_ssh_key.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_ssh_key.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "name", "tf_test_ssh_key"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_ssh_key.test", "tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_ssh_key.test", "tags.*", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "username", "testuser"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "password", "testpassword"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "passphrase", "testpassphrase"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "private_key_data", "test-private-key-data"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "public_key", "test-public-key"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialSSHKeyResourceConfig(
					"tf_test_ssh_key", "tf_test_ssh_key_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updatedpassword", "updatedpassphrase", "updated-private-key-data", "updated-public-key",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "name", "tf_test_ssh_key_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "username", "updateduser"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "password", "updatedpassword"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "passphrase", "updatedpassphrase"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "private_key_data", "updated-private-key-data"),
					resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "public_key", "updated-public-key"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_ssh_key.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_ssh_key.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialSSHKeyResourceConfig(vaultName, name, description, folder, username, password, passphrase, privateKeyData, publicKey string) string {
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

resource "dvls_entry_credential_ssh_key" "test" {
  vault_id         = dvls_vault.test.id
  name             = %[3]q
  description      = %[4]q
  folder           = %[5]q
  tags             = ["acceptance", "tf-test"]
  username         = %[6]q
  password         = %[7]q
  passphrase       = %[8]q
  private_key_data = %[9]q
  public_key       = %[10]q

  depends_on = [dvls_entry_folder.default, dvls_entry_folder.updated]
}
`, testAccProviderConfig(), vaultName, name, description, folder, username, password, passphrase, privateKeyData, publicKey)
}
