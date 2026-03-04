package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSSHKeyDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyDataSourceConfig_byName("tf_test_ssh_key_by_name", "tf_test_ssh_key_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "id", "dvls_entry_credential_ssh_key.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "vault_id", "dvls_entry_credential_ssh_key.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "name", "dvls_entry_credential_ssh_key.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "description", "dvls_entry_credential_ssh_key.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "folder", "dvls_entry_credential_ssh_key.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "tags.#", "dvls_entry_credential_ssh_key.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "username", "dvls_entry_credential_ssh_key.test", "username"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "password", "dvls_entry_credential_ssh_key.test", "password"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "passphrase", "dvls_entry_credential_ssh_key.test", "passphrase"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "private_key_data", "dvls_entry_credential_ssh_key.test", "private_key_data"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "public_key", "dvls_entry_credential_ssh_key.test", "public_key"),
				),
			},
		},
	})
}

func TestAccEntryCredentialSSHKeyDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyDataSourceConfig_byId("tf_test_ssh_key_by_id", "tf_test_ssh_key_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "id", "dvls_entry_credential_ssh_key.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "vault_id", "dvls_entry_credential_ssh_key.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "name", "dvls_entry_credential_ssh_key.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "description", "dvls_entry_credential_ssh_key.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "folder", "dvls_entry_credential_ssh_key.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "tags.#", "dvls_entry_credential_ssh_key.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "username", "dvls_entry_credential_ssh_key.test", "username"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "password", "dvls_entry_credential_ssh_key.test", "password"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "passphrase", "dvls_entry_credential_ssh_key.test", "passphrase"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "private_key_data", "dvls_entry_credential_ssh_key.test", "private_key_data"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_ssh_key.test", "public_key", "dvls_entry_credential_ssh_key.test", "public_key"),
				),
			},
		},
	})
}

func testAccEntryCredentialSSHKeyDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_ssh_key" "test" {
  vault_id         = dvls_vault.test.id
  name             = %[3]q
  description      = "test entry for data source"
  folder           = "tf_test_folder"
  tags             = ["tf-test", "acceptance"]
  username         = "testuser"
  password         = "testpassword"
  passphrase       = "testpassphrase"
  private_key_data = "test-private-key-data"
  public_key       = "test-public-key"
}

data "dvls_entry_credential_ssh_key" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_ssh_key.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialSSHKeyDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_ssh_key" "test" {
  vault_id         = dvls_vault.test.id
  name             = %[3]q
  description      = "test entry for data source"
  folder           = "tf_test_folder"
  tags             = ["tf-test", "acceptance"]
  username         = "testuser"
  password         = "testpassword"
  passphrase       = "testpassphrase"
  private_key_data = "test-private-key-data"
  public_key       = "test-public-key"
}

data "dvls_entry_credential_ssh_key" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_ssh_key.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
