package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialConnectionStringDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialConnectionStringDataSourceConfig("tf_test_connection_string_by_name", "tf_test_connection_string_by_name", "name"),
				Check:  testAccEntryCredentialConnectionStringDataSourceCheck(),
			},
		},
	})
}

func TestAccEntryCredentialConnectionStringDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialConnectionStringDataSourceConfig("tf_test_connection_string_by_id", "tf_test_connection_string_by_id", "id"),
				Check:  testAccEntryCredentialConnectionStringDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialConnectionStringDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck("dvls_entry_credential_connection_string", "connection_string")
}

func testAccEntryCredentialConnectionStringDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_connection_string",
		vaultName, name,
		`  connection_string = "Server=localhost;Database=test;Trusted_Connection=True;"`,
		lookupField,
	)
}
