package provider

import (
	"fmt"
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
			{
				Config: testAccEntryCredentialSecretEphemeralConfig("tf_test_secret_eph_byname", "tf_test_secret_eph_byname", `
ephemeral "dvls_entry_credential_secret" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_secret.test.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.secret", "my-secret-value-123"),
					resource.TestCheckResourceAttrPair("echo.test", "data.id", "dvls_entry_credential_secret.test", "id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.0", "tf-test"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.1", "acceptance"),
				),
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
			{
				Config: testAccEntryCredentialSecretEphemeralConfig("tf_test_secret_eph_byid", "tf_test_secret_eph_byid", `
ephemeral "dvls_entry_credential_secret" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_secret.test.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.secret", "my-secret-value-123"),
					resource.TestCheckResourceAttrPair("echo.test", "data.id", "dvls_entry_credential_secret.test", "id"),
					resource.TestCheckResourceAttr("echo.test", "data.description", "test entry for ephemeral resource"),
					resource.TestCheckResourceAttr("echo.test", "data.folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.#", "2"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.0", "tf-test"),
					resource.TestCheckResourceAttr("echo.test", "data.tags.1", "acceptance"),
				),
			},
		},
	})
}

func testAccEntryCredentialSecretEphemeralConfig(vaultName, entryName, ephemeralBlock string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_secret" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = "test entry for ephemeral resource"
  folder      = "tf_test_folder"
  tags        = ["tf-test", "acceptance"]
  secret      = "my-secret-value-123"
}

%s

%s
`, testAccProviderConfig(), vaultName, entryName, ephemeralBlock, testAccEphemeralEchoConfig("ephemeral.dvls_entry_credential_secret.test"))
}
