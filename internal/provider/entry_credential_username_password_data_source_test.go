package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialUsernamePasswordDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_username_password_by_name", "tf_test_folder"),
			{
				Config: testAccEntryCredentialUsernamePasswordDataSourceConfig_byName("tf_test_username_password_by_name", "tf_test_username_password_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "id", "dvls_entry_credential_username_password.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "vault_id", "dvls_entry_credential_username_password.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "name", "dvls_entry_credential_username_password.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "description", "dvls_entry_credential_username_password.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "folder", "dvls_entry_credential_username_password.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "tags.#", "dvls_entry_credential_username_password.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "username", "dvls_entry_credential_username_password.test", "username"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "domain", "dvls_entry_credential_username_password.test", "domain"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "password", "dvls_entry_credential_username_password.test", "password"),
				),
			},
		},
	})
}

func TestAccEntryCredentialUsernamePasswordDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_username_password_by_id", "tf_test_folder"),
			{
				Config: testAccEntryCredentialUsernamePasswordDataSourceConfig_byId("tf_test_username_password_by_id", "tf_test_username_password_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "id", "dvls_entry_credential_username_password.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "vault_id", "dvls_entry_credential_username_password.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "name", "dvls_entry_credential_username_password.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "description", "dvls_entry_credential_username_password.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "folder", "dvls_entry_credential_username_password.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "tags.#", "dvls_entry_credential_username_password.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "username", "dvls_entry_credential_username_password.test", "username"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "domain", "dvls_entry_credential_username_password.test", "domain"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_username_password.test", "password", "dvls_entry_credential_username_password.test", "password"),
				),
			},
		},
	})
}

func testAccEntryCredentialUsernamePasswordDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_username_password" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["acceptance", "tf-test"]
  username    = "testuser"
  domain      = "testdomain"
  password    = "testpassword123"
}

data "dvls_entry_credential_username_password" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_username_password.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialUsernamePasswordDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_username_password" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["acceptance", "tf-test"]
  username    = "testuser"
  domain      = "testdomain"
  password    = "testpassword123"
}

data "dvls_entry_credential_username_password" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_username_password.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
