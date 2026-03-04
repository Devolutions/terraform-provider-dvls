package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialAzureServicePrincipalDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialAzureServicePrincipalDataSourceConfig_byName("tf_test_azure_sp_by_name", "tf_test_azure_sp_by_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "id", "dvls_entry_credential_azure_service_principal.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "vault_id", "dvls_entry_credential_azure_service_principal.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "name", "dvls_entry_credential_azure_service_principal.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "description", "dvls_entry_credential_azure_service_principal.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "folder", "dvls_entry_credential_azure_service_principal.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "tags.#", "dvls_entry_credential_azure_service_principal.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "client_id", "dvls_entry_credential_azure_service_principal.test", "client_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "client_secret", "dvls_entry_credential_azure_service_principal.test", "client_secret"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "tenant_id", "dvls_entry_credential_azure_service_principal.test", "tenant_id"),
				),
			},
		},
	})
}

func TestAccEntryCredentialAzureServicePrincipalDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialAzureServicePrincipalDataSourceConfig_byId("tf_test_azure_sp_by_id", "tf_test_azure_sp_by_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "id", "dvls_entry_credential_azure_service_principal.test", "id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "vault_id", "dvls_entry_credential_azure_service_principal.test", "vault_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "name", "dvls_entry_credential_azure_service_principal.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "description", "dvls_entry_credential_azure_service_principal.test", "description"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "folder", "dvls_entry_credential_azure_service_principal.test", "folder"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "tags.#", "dvls_entry_credential_azure_service_principal.test", "tags.#"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "client_id", "dvls_entry_credential_azure_service_principal.test", "client_id"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "client_secret", "dvls_entry_credential_azure_service_principal.test", "client_secret"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_credential_azure_service_principal.test", "tenant_id", "dvls_entry_credential_azure_service_principal.test", "tenant_id"),
				),
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalDataSourceConfig_byName(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_azure_service_principal" "test" {
  vault_id      = dvls_vault.test.id
  name          = %[3]q
  description   = "test entry for data source"
  folder        = "tf_test_folder"
  tags          = ["tf-test", "acceptance"]
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
  tenant_id     = "test-tenant-id"
}

data "dvls_entry_credential_azure_service_principal" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_azure_service_principal.test.name
}
`, testAccProviderConfig(), vaultName, name)
}

func testAccEntryCredentialAzureServicePrincipalDataSourceConfig_byId(vaultName, name string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_azure_service_principal" "test" {
  vault_id      = dvls_vault.test.id
  name          = %[3]q
  description   = "test entry for data source"
  folder        = "tf_test_folder"
  tags          = ["tf-test", "acceptance"]
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
  tenant_id     = "test-tenant-id"
}

data "dvls_entry_credential_azure_service_principal" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_azure_service_principal.test.id
}
`, testAccProviderConfig(), vaultName, name)
}
