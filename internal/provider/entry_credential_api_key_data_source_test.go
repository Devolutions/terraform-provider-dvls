package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialApiKeyDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialApiKeyDataSourceConfig_byName("tf_test_api_key_by_name", "tf_test_api_key_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "id", "dvls_entry_credential_api_key.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "vault_id", "dvls_entry_credential_api_key.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "name", "dvls_entry_credential_api_key.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "description", "dvls_entry_credential_api_key.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "folder", "dvls_entry_credential_api_key.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "tags.#", "dvls_entry_credential_api_key.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "api_id", "dvls_entry_credential_api_key.test", "api_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "api_key", "dvls_entry_credential_api_key.test", "api_key"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "tenant_id", "dvls_entry_credential_api_key.test", "tenant_id"),
				),
			},
		},
	})
}

func TestAccEntryCredentialApiKeyDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialApiKeyDataSourceConfig_byId("tf_test_api_key_by_id", "tf_test_api_key_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "id", "dvls_entry_credential_api_key.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "vault_id", "dvls_entry_credential_api_key.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "name", "dvls_entry_credential_api_key.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "description", "dvls_entry_credential_api_key.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "folder", "dvls_entry_credential_api_key.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "tags.#", "dvls_entry_credential_api_key.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "api_id", "dvls_entry_credential_api_key.test", "api_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "api_key", "dvls_entry_credential_api_key.test", "api_key"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_api_key.test", "tenant_id", "dvls_entry_credential_api_key.test", "tenant_id"),
				),
			},
		},
	})
}

func testAccEntryCredentialApiKeyDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_api_key" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["tf-test", "acceptance"]
  api_id      = "test-api-id"
  api_key     = "test-api-key-secret"
  tenant_id   = "test-tenant-id"
}

data "dvls_entry_credential_api_key" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_api_key.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialApiKeyDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_api_key" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for data source"
  folder      = "tf_test_folder"
  tags        = ["tf-test", "acceptance"]
  api_id      = "test-api-id"
  api_key     = "test-api-key-secret"
  tenant_id   = "test-tenant-id"
}

data "dvls_entry_credential_api_key" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_api_key.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
