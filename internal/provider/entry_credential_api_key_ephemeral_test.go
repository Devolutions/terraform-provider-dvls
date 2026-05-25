package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialApiKeyEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_api_key_eph_byname", testAccEphFolder),
			{Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byname", "tf_test_api_key_eph_byname", "")},
			{
				Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byname", "tf_test_api_key_eph_byname", "name"),
				Check:  testAccEntryCredentialApiKeyEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialApiKeyEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_api_key_eph_byid", testAccEphFolder),
			{
				Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byid", "tf_test_api_key_eph_byid", "id"),
				Check:  testAccEntryCredentialApiKeyEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialApiKeyEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.api_id", "test-api-id"),
		resource.TestCheckResourceAttr("echo.test", "data.api_key", "test-api-key-secret"),
		resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccEphDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccEphFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckResourceAttr("echo.test", "data.tags.0", testAccEphTags[0]),
		resource.TestCheckResourceAttr("echo.test", "data.tags.1", testAccEphTags[1]),
	)
}

func testAccEntryCredentialApiKeyEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_api_key",
		vaultName,
		entryName,
		`  api_id = "test-api-id"
  api_key = "test-api-key-secret"
  tenant_id = "test-tenant-id"`,
		lookupField,
	)
}
