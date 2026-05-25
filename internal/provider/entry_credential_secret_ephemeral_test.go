package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSecretEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_secret_eph_byname", testAccEphFolder),
			{Config: testAccEntryCredentialSecretEphemeralConfig("tf_test_secret_eph_byname", "tf_test_secret_eph_byname", "")},
			{
				Config: testAccEntryCredentialSecretEphemeralConfig("tf_test_secret_eph_byname", "tf_test_secret_eph_byname", "name"),
				Check:  testAccEntryCredentialSecretEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialSecretEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_secret_eph_byid", testAccEphFolder),
			{
				Config: testAccEntryCredentialSecretEphemeralConfig("tf_test_secret_eph_byid", "tf_test_secret_eph_byid", "id"),
				Check:  testAccEntryCredentialSecretEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialSecretEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.secret", "my-secret-value-123"),
		resource.TestCheckResourceAttrPair("echo.test", "data.id", "dvls_entry_credential_secret.test", "id"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccEphDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccEphFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckResourceAttr("echo.test", "data.tags.0", testAccEphTags[0]),
		resource.TestCheckResourceAttr("echo.test", "data.tags.1", testAccEphTags[1]),
	)
}

func testAccEntryCredentialSecretEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_secret",
		vaultName,
		entryName,
		`  secret = "my-secret-value-123"`,
		lookupField,
	)
}
