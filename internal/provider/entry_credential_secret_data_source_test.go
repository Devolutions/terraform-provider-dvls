package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSecretDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSecretDataSourceConfig_byName("tf_test_secret_by_name", "tf_test_secret_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "id", "dvls_entry_credential_secret.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "vault_id", "dvls_entry_credential_secret.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "name", "dvls_entry_credential_secret.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "description", "dvls_entry_credential_secret.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "folder", "dvls_entry_credential_secret.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "tags.#", "dvls_entry_credential_secret.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "secret", "dvls_entry_credential_secret.test", "secret"),
				),
			},
		},
	})
}

func TestAccEntryCredentialSecretDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSecretDataSourceConfig_byId("tf_test_secret_by_id", "tf_test_secret_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "id", "dvls_entry_credential_secret.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "vault_id", "dvls_entry_credential_secret.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "name", "dvls_entry_credential_secret.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "description", "dvls_entry_credential_secret.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "folder", "dvls_entry_credential_secret.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "tags.#", "dvls_entry_credential_secret.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_secret.test", "secret", "dvls_entry_credential_secret.test", "secret"),
				),
			},
		},
	})
}

func testAccEntryCredentialSecretDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_secret" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["tf-test", "acceptance"]
  secret      = "my-secret-value-123"
}

data "dvls_entry_credential_secret" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_secret.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialSecretDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_secret" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["tf-test", "acceptance"]
  secret      = "my-secret-value-123"
}

data "dvls_entry_credential_secret" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_secret.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
