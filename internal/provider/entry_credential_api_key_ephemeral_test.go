package provider

import (
	"fmt"
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
			{Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byname", "tf_test_api_key_eph_byname", "")},
			{
				Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byname", "tf_test_api_key_eph_byname", `
ephemeral "dvls_entry_credential_api_key" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_api_key.test.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.api_id", "test-api-id"),
					resource.TestCheckResourceAttr("echo.test", "data.api_key", "test-api-key-secret"),
					resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", "tf-test"),
				),
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
			{
				Config: testAccEntryCredentialApiKeyEphemeralConfig("tf_test_api_key_eph_byid", "tf_test_api_key_eph_byid", `
ephemeral "dvls_entry_credential_api_key" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_api_key.test.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.api_id", "test-api-id"),
					resource.TestCheckResourceAttr("echo.test", "data.api_key", "test-api-key-secret"),
					resource.TestCheckResourceAttr("echo.test", "data.tenant_id", "test-tenant-id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr("echo.test", "data.tags.*", "tf-test"),
				),
			},
		},
	})
}

func testAccEntryCredentialApiKeyEphemeralConfig(vaultName, entryName, ephemeralBlock string) string {
	echoConfig := ""
	if ephemeralBlock != "" {
		echoConfig = testAccEphemeralEchoConfig("ephemeral.dvls_entry_credential_api_key.test")
	}

	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder"
}

resource "dvls_entry_credential_api_key" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for ephemeral resource"
  folder      = "tf_test_folder"
  tags        = ["acceptance", "tf-test"]
  api_id      = "test-api-id"
  api_key     = "test-api-key-secret"
  tenant_id   = "test-tenant-id"

  depends_on = [dvls_entry_folder.default]
}

%s

%s
`, testAccProviderConfig(), vaultName, entryName, ephemeralBlock, echoConfig)
}
