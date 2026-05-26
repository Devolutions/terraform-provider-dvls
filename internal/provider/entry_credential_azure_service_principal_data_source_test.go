package provider

import (
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
				Config: testAccEntryCredentialAzureServicePrincipalDataSourceConfig("tf_test_azure_sp_by_name", "tf_test_azure_sp_by_name", "name"),
				Check:  testAccEntryCredentialAzureServicePrincipalDataSourceCheck(),
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
				Config: testAccEntryCredentialAzureServicePrincipalDataSourceConfig("tf_test_azure_sp_by_id", "tf_test_azure_sp_by_id", "id"),
				Check:  testAccEntryCredentialAzureServicePrincipalDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck(
		"dvls_entry_credential_azure_service_principal",
		"client_id", "client_secret", "tenant_id",
	)
}

func testAccEntryCredentialAzureServicePrincipalDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_azure_service_principal",
		vaultName, name,
		`  client_id = "test-client-id"
  client_secret = "test-client-secret"
  tenant_id = "test-tenant-id"`,
		lookupField,
	)
}
