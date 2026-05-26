package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSSHKeyEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{Config: testAccEntryCredentialSSHKeyEphemeralConfig("tf_test_sshkey_eph_byname", "tf_test_sshkey_eph_byname", "")},
			{
				Config: testAccEntryCredentialSSHKeyEphemeralConfig("tf_test_sshkey_eph_byname", "tf_test_sshkey_eph_byname", "name"),
				Check:  testAccEntryCredentialSSHKeyEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialSSHKeyEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyEphemeralConfig("tf_test_sshkey_eph_byid", "tf_test_sshkey_eph_byid", "id"),
				Check:  testAccEntryCredentialSSHKeyEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialSSHKeyEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.username", "testuser"),
		resource.TestCheckResourceAttr("echo.test", "data.password", "testpassword"),
		resource.TestCheckResourceAttr("echo.test", "data.passphrase", "testpassphrase"),
		resource.TestCheckResourceAttrSet("echo.test", "data.private_key_data"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccTestDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccTestFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[1]),
	)
}

func testAccEntryCredentialSSHKeyEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_ssh_key",
		vaultName,
		entryName,
		`  username = "testuser"
  password = "testpassword"
  passphrase = "testpassphrase"
  private_key_data = "-----BEGIN OPENSSH PRIVATE KEY-----\nfake-key-data\n-----END OPENSSH PRIVATE KEY-----"`,
		lookupField,
	)
}
