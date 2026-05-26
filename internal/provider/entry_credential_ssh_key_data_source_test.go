package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSSHKeyDataSource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyDataSourceConfig("tf_test_ssh_key_by_name", "tf_test_ssh_key_by_name", "name"),
				Check:  testAccEntryCredentialSSHKeyDataSourceCheck(),
			},
		},
	})
}

func TestAccEntryCredentialSSHKeyDataSource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyDataSourceConfig("tf_test_ssh_key_by_id", "tf_test_ssh_key_by_id", "id"),
				Check:  testAccEntryCredentialSSHKeyDataSourceCheck(),
			},
		},
	})
}

func testAccEntryCredentialSSHKeyDataSourceCheck() resource.TestCheckFunc {
	return testAccEntryCredentialDataSourceCheck(
		"dvls_entry_credential_ssh_key",
		"username", "password", "passphrase", "private_key_data", "public_key",
	)
}

func testAccEntryCredentialSSHKeyDataSourceConfig(vaultName, name, lookupField string) string {
	return testAccEntryCredentialDataSourceConfig(
		"dvls_entry_credential_ssh_key",
		vaultName, name,
		`  username = "testuser"
  password = "testpassword"
  passphrase = "testpassphrase"
  private_key_data = "-----BEGIN OPENSSH PRIVATE KEY-----\nfake-key-data\n-----END OPENSSH PRIVATE KEY-----"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB"`,
		lookupField,
	)
}
