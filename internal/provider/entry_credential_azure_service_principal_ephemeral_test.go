package provider

import (
	"fmt"
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
			testAccVaultWithFoldersStep("tf_test_azsp_eph_byname", "tf_test_folder"),
			{Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byname", "tf_test_azsp_eph_byname", "")},
			{
				Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byname", "tf_test_azsp_eph_byname", `
ephemeral "dvls_entry_credential_azure_service_principal" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_azure_service_principal.test.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.client_id", "test-client-id"),
					resource.TestCheckResourceAttr("echo.test", "data.client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.0", "acceptance"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.1", "tf-test"),
				),
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
			testAccVaultWithFoldersStep("tf_test_azsp_eph_byid", "tf_test_folder"),
			{
				Config: testAccEntryCredentialAzureServicePrincipalEphemeralConfig("tf_test_azsp_eph_byid", "tf_test_azsp_eph_byid", `
ephemeral "dvls_entry_credential_azure_service_principal" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_azure_service_principal.test.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.client_id", "test-client-id"),
					resource.TestCheckResourceAttr("echo.test", "data.client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.0", "acceptance"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.1", "tf-test"),
				),
			},
		},
	})
}

func testAccEntryCredentialAzureServicePrincipalEphemeralConfig(vaultName, entryName, ephemeralBlock string) string {
	echoConfig := ""
	if ephemeralBlock != "" {
		echoConfig = testAccEphemeralEchoConfig("ephemeral.dvls_entry_credential_azure_service_principal.test")
	}

	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_azure_service_principal" "test" {
  vault_id      = dvls_vault.test.id
  name          = %[3]q
  description   = "test entry for ephemeral resource"
  folder        = "tf_test_folder"
  tags          = ["acceptance", "tf-test"]
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
  tenant_id     = "test-tenant-id"
}

%s

%s
`, testAccProviderConfig(), vaultName, entryName, ephemeralBlock, echoConfig)
}
