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
			{
				Config: testAccEntryCredentialSSHKeyResourceConfig(
					"tf_test_ssh_key", "tf_test_ssh_key", "test description", "tf_test_folder",
					"testuser", "testpassword", "testpassphrase",
					"-----BEGIN OPENSSH PRIVATE KEY-----\\nfake-key-data\\n-----END OPENSSH PRIVATE KEY-----",
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB",
				),
				Check: testAccEntryCredentialSSHKeyResourceCheck(
					"tf_test_ssh_key", "test description", "tf_test_folder",
					"testuser", "testpassword", "testpassphrase",
					"-----BEGIN OPENSSH PRIVATE KEY-----\nfake-key-data\n-----END OPENSSH PRIVATE KEY-----",
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB",
				),
			},
			{
				Config: testAccEntryCredentialSSHKeyResourceConfig(
					"tf_test_ssh_key", "tf_test_ssh_key_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updatedpassword", "updatedpassphrase",
					"-----BEGIN OPENSSH PRIVATE KEY-----\\nupdated-key-data\\n-----END OPENSSH PRIVATE KEY-----",
					"ssh-rsa UPDATEDAAAB3NzaC1yc2EAAAADAQABAAAB",
				),
				Check: testAccEntryCredentialSSHKeyResourceCheck(
					"tf_test_ssh_key_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updatedpassword", "updatedpassphrase",
					"-----BEGIN OPENSSH PRIVATE KEY-----\nupdated-key-data\n-----END OPENSSH PRIVATE KEY-----",
					"ssh-rsa UPDATEDAAAB3NzaC1yc2EAAAADAQABAAAB",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_ssh_key.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_ssh_key.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialSSHKeyResourceCheck(name, description, folder, username, password, passphrase, privateKeyData, publicKey string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_ssh_key.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_ssh_key.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_ssh_key.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_ssh_key.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "username", username),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "password", password),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "passphrase", passphrase),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "private_key_data", privateKeyData),
		resource.TestCheckResourceAttr("dvls_entry_credential_ssh_key.test", "public_key", publicKey),
	)
}

func testAccEntryCredentialSSHKeyResourceConfig(vaultName, name, description, folder, username, password, passphrase, privateKeyData, publicKey string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_ssh_key",
		vaultName, name, description, folder,
		fmt.Sprintf(`  username = %q
  password = %q
  passphrase = %q
  private_key_data = "%s"
  public_key = %q`, username, password, passphrase, privateKeyData, publicKey),
	)
}
