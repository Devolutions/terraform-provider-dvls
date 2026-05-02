package provider

import (
	"fmt"
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
				Config: testAccEntryCredentialSSHKeyEphemeralConfig("tf_test_sshkey_eph_byname", "tf_test_sshkey_eph_byname", `
ephemeral "dvls_entry_credential_ssh_key" "test" {
  vault_id = dvls_vault.test.id
  name     = dvls_entry_credential_ssh_key.test.name
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.username", "testuser"),
					resource.TestCheckResourceAttr("echo.test", "data.password", "testpassword"),
					resource.TestCheckResourceAttr("echo.test", "data.passphrase", "testpassphrase"),
					resource.TestCheckResourceAttrSet("echo.test", "data.private_key_data"),
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

func TestAccEntryCredentialSSHKeyEphemeralResource_byId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks:   testAccEphemeralTerraformVersionCheck,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialSSHKeyEphemeralConfig("tf_test_sshkey_eph_byid", "tf_test_sshkey_eph_byid", `
ephemeral "dvls_entry_credential_ssh_key" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_credential_ssh_key.test.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.username", "testuser"),
					resource.TestCheckResourceAttr("echo.test", "data.password", "testpassword"),
					resource.TestCheckResourceAttr("echo.test", "data.passphrase", "testpassphrase"),
					resource.TestCheckResourceAttrSet("echo.test", "data.private_key_data"),
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

func testAccEntryCredentialSSHKeyEphemeralConfig(vaultName, entryName, ephemeralBlock string) string {
	echoConfig := ""
	if ephemeralBlock != "" {
		echoConfig = testAccEphemeralEchoConfig("ephemeral.dvls_entry_credential_ssh_key.test")
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

resource "dvls_entry_credential_ssh_key" "test" {
  vault_id         = dvls_vault.test.id
  name             = %[3]q
  description      = "test entry for ephemeral resource"
  folder           = "tf_test_folder"
  tags             = ["acceptance", "tf-test"]
  username         = "testuser"
  password         = "testpassword"
  passphrase       = "testpassphrase"
  private_key_data = "-----BEGIN OPENSSH PRIVATE KEY-----\nfake-key-data\n-----END OPENSSH PRIVATE KEY-----"

  depends_on = [dvls_entry_folder.default]
}

%s

%s
`, testAccProviderConfig(), vaultName, entryName, ephemeralBlock, echoConfig)
}
