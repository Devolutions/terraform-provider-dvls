package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialUsernamePasswordDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialUsernamePasswordDataSourceConfig("tf_test_username_password_by_name", "tf_test_username_password_by_name", "name"),
				Check:  testAccEntryCredentialUsernamePasswordDataSourceCheck(),
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
			{
				Config: testAccEntryCredentialUsernamePasswordDataSourceConfig("tf_test_username_password_by_id", "tf_test_username_password_by_id", "id"),
				Check:  testAccEntryCredentialUsernamePasswordDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialUsernamePasswordDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck(
		"dvls_entry_credential_username_password",
		"username", "domain", "password",
	)
}

func testAccEntryCredentialUsernamePasswordDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_username_password",
		vaultName, name,
		`  username = "testuser"
  domain = "testdomain"
  password = "testpassword"`,
		lookupField,
	)
}
