package provider

import (
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
				Config: testAccEntryCredentialSecretDataSourceConfig("tf_test_secret_by_name", "tf_test_secret_by_name", "name"),
				Check:  testAccEntryCredentialSecretDataSourceCheck(),
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
				Config: testAccEntryCredentialSecretDataSourceConfig("tf_test_secret_by_id", "tf_test_secret_by_id", "id"),
				Check:  testAccEntryCredentialSecretDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialSecretDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck("dvls_entry_credential_secret", "secret")
}

func testAccEntryCredentialSecretDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_secret",
		vaultName, name,
		`  secret = "my-secret-value-123"`,
		lookupField,
	)
}
