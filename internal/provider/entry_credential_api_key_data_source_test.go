package provider

import (
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
				Config: testAccEntryCredentialApiKeyDataSourceConfig("tf_test_api_key_by_name", "tf_test_api_key_by_name", "name"),
				Check:  testAccEntryCredentialApiKeyDataSourceCheck(),
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
				Config: testAccEntryCredentialApiKeyDataSourceConfig("tf_test_api_key_by_id", "tf_test_api_key_by_id", "id"),
				Check:  testAccEntryCredentialApiKeyDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialApiKeyDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck("dvls_entry_credential_api_key", "api_id", "api_key", "tenant_id")
}

func testAccEntryCredentialApiKeyDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_api_key",
		vaultName, name,
		`  api_id = "test-api-id"
  api_key = "test-api-key-secret"
  tenant_id = "test-tenant-id"`,
		lookupField,
	)
}
