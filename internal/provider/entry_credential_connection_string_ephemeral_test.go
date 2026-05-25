package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialConnectionStringEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_connstr_eph_byname", testAccEphFolder),
			{Config: testAccEntryCredentialConnectionStringEphemeralConfig("tf_test_connstr_eph_byname", "tf_test_connstr_eph_byname", "")},
			{
				Config: testAccEntryCredentialConnectionStringEphemeralConfig("tf_test_connstr_eph_byname", "tf_test_connstr_eph_byname", "name"),
				Check:  testAccEntryCredentialConnectionStringEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialConnectionStringEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_connstr_eph_byid", testAccEphFolder),
			{
				Config: testAccEntryCredentialConnectionStringEphemeralConfig("tf_test_connstr_eph_byid", "tf_test_connstr_eph_byid", "id"),
				Check:  testAccEntryCredentialConnectionStringEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialConnectionStringEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.connection_string", "Server=localhost;Database=test;Trusted_Connection=True;"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccEphDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccEphFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckResourceAttr("echo.test", "data.tags.0", testAccEphTags[0]),
		resource.TestCheckResourceAttr("echo.test", "data.tags.1", testAccEphTags[1]),
	)
}

func testAccEntryCredentialConnectionStringEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_connection_string",
		vaultName,
		entryName,
		`  connection_string = "Server=localhost;Database=test;Trusted_Connection=True;"`,
		lookupField,
	)
}
