package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialAzureServicePrincipalResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccEntryCredentialAzureServicePrincipalResourceConfig(
					"tf_test_azure_sp", "tf_test_azure_sp", "test description", "tf_test_folder",
					"test-client-id", "test-client-secret", "test-tenant-id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_azure_service_principal.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_azure_service_principal.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "name", "tf_test_azure_sp"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_azure_service_principal.test", "tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_credential_azure_service_principal.test", "tags.*", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_id", "test-client-id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "tenant_id", "test-tenant-id"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialAzureServicePrincipalResourceConfig(
					"tf_test_azure_sp", "tf_test_azure_sp_updated", "updated description", "tf_test_folder_updated",
					"updated-client-id", "updated-client-secret", "updated-tenant-id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "name", "tf_test_azure_sp_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_id", "updated-client-id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_secret", "updated-client-secret"),
					resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "tenant_id", "updated-tenant-id"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_azure_service_principal.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_azure_service_principal.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalResourceConfig(vaultName, name, description, folder, clientId, clientSecret, tenantId string) string {
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

resource "dvls_entry_credential_azure_service_principal" "test" {
  vault_id      = dvls_vault.test.id
  name          = %[3]q
  description   = %[4]q
  folder        = %[5]q
  tags          = ["acceptance", "tf-test"]
  client_id     = %[6]q
  client_secret = %[7]q
  tenant_id     = %[8]q

  depends_on = [dvls_entry_folder.default, dvls_entry_folder.updated]
}
`, testAccProviderConfig(), vaultName, name, description, folder, clientId, clientSecret, tenantId)
}
