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
			{
				Config: testAccEntryCredentialAzureServicePrincipalResourceConfig(
					"tf_test_azure_sp", "tf_test_azure_sp", "test description", "tf_test_folder",
					"test-client-id", "test-client-secret", "test-tenant-id",
				),
				Check: testAccEntryCredentialAzureServicePrincipalResourceCheck(
					"tf_test_azure_sp", "test description", "tf_test_folder",
					"test-client-id", "test-client-secret", "test-tenant-id",
				),
			},
			{
				Config: testAccEntryCredentialAzureServicePrincipalResourceConfig(
					"tf_test_azure_sp", "tf_test_azure_sp_updated", "updated description", "tf_test_folder_updated",
					"updated-client-id", "updated-client-secret", "updated-tenant-id",
				),
				Check: testAccEntryCredentialAzureServicePrincipalResourceCheck(
					"tf_test_azure_sp_updated", "updated description", "tf_test_folder_updated",
					"updated-client-id", "updated-client-secret", "updated-tenant-id",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_azure_service_principal.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_azure_service_principal.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalResourceCheck(name, description, folder, clientId, clientSecret, tenantId string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_azure_service_principal.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_azure_service_principal.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_azure_service_principal.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_azure_service_principal.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_id", clientId),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "client_secret", clientSecret),
		resource.TestCheckResourceAttr("dvls_entry_credential_azure_service_principal.test", "tenant_id", tenantId),
	)
}

func testAccEntryCredentialAzureServicePrincipalResourceConfig(vaultName, name, description, folder, clientId, clientSecret, tenantId string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_azure_service_principal",
		vaultName, name, description, folder,
		fmt.Sprintf(`  client_id = %q
  client_secret = %q
  tenant_id = %q`, clientId, clientSecret, tenantId),
	)
}
