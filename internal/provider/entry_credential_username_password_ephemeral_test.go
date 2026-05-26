package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialUsernamePasswordEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{Config: testAccEntryCredentialUsernamePasswordEphemeralConfig("tf_test_userpass_eph_byname", "tf_test_userpass_eph_byname", "")},
			{
				Config: testAccEntryCredentialUsernamePasswordEphemeralConfig("tf_test_userpass_eph_byname", "tf_test_userpass_eph_byname", "name"),
				Check:  testAccEntryCredentialUsernamePasswordEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialUsernamePasswordEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialUsernamePasswordEphemeralConfig("tf_test_userpass_eph_byid", "tf_test_userpass_eph_byid", "id"),
				Check:  testAccEntryCredentialUsernamePasswordEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialUsernamePasswordEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.username", "testuser"),
		resource.TestCheckResourceAttr("echo.test", "data.domain", "testdomain"),
		resource.TestCheckResourceAttr("echo.test", "data.password", "testpassword123"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccTestDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccTestFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[1]),
	)
}

func testAccEntryCredentialUsernamePasswordEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_username_password",
		vaultName,
		entryName,
		`  username = "testuser"
  domain = "testdomain"
  password = "testpassword123"`,
		lookupField,
	)
}
