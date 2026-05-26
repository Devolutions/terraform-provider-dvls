package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialAzureServicePrincipalEphemeralResource_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byname", "tf_test_azsp_eph_byname", "")},
			{
				Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byname", "tf_test_azsp_eph_byname", "name"),
				Check:  testAccEntryCredentialAzureServicePrincipalEphemeralCheck(),
			},
		},
	})
}

func TestAccEntryCredentialAzureServicePrincipalEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byid", "tf_test_azsp_eph_byid", "id"),
				Check:  testAccEntryCredentialAzureServicePrincipalEphemeralCheck(),
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalEphemeralCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("echo.test", "data.client_id", "test-client-id"),
		resource.TestCheckResourceAttr("echo.test", "data.client_secret", "test-client-secret"),
		resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
		resource.TestCheckResourceAttr("echo.test", "data.description", testAccTestDescription),
		resource.TestCheckResourceAttr("echo.test", "data.folder", testAccTestFolder),
		resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", testAccTestTags[1]),
	)
}

func testAccEntryCredentialAzureServicePrincipalEphemeralConfig(vaultName, entryName, lookupField string) string {
	return testAccEntryCredentialEphemeralConfig(
		"dvls_entry_credential_azure_service_principal",
		vaultName,
		entryName,
		`  client_id = "test-client-id"
  client_secret = "test-client-secret"
  tenant_id = "test-tenant-id"`,
		lookupField,
	)
}
