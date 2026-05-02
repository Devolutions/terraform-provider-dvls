package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// upgradeFromV0Test asserts that state created with the published v0.6.0
// provider (where `tags` was a list) plans clean against the local provider
// (where `tags` is a set). A clean plan after the switch means the upgrader
// migrated the value silently.
func upgradeFromV0Test(t *testing.T, address, cfg string) {
	t.Helper()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"dvls": {Source: "devolutions/dvls", VersionConstraint: "0.6.0"},
				},
				Config: cfg,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(address, "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr(address, "tags.*", "acceptance"),
					resource.TestCheckTypeSetElemAttr(address, "tags.*", "tf-test"),
				),
			},
		},
	})
}

func TestAccEntryCredentialSecretResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_secret.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_secret_upgrade"
}

resource "dvls_entry_credential_secret" "test" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_secret_upgrade"
  tags     = ["acceptance", "tf-test"]
  secret   = "test-secret"
}
`, testAccProviderConfig()))
}

func TestAccEntryCredentialApiKeyResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_api_key.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_api_key_upgrade"
}

resource "dvls_entry_credential_api_key" "test" {
  vault_id  = dvls_vault.test.id
  name      = "tf_test_api_key_upgrade"
  tags      = ["acceptance", "tf-test"]
  api_id    = "id"
  api_key   = "key"
  tenant_id = "tenant"
}
`, testAccProviderConfig()))
}

func TestAccEntryCredentialSSHKeyResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_ssh_key.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_ssh_key_upgrade"
}

resource "dvls_entry_credential_ssh_key" "test" {
  vault_id         = dvls_vault.test.id
  name             = "tf_test_ssh_key_upgrade"
  tags             = ["acceptance", "tf-test"]
  username         = "user"
  private_key_data = "data"
}
`, testAccProviderConfig()))
}

func TestAccEntryCredentialConnectionStringResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_connection_string.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_connection_string_upgrade"
}

resource "dvls_entry_credential_connection_string" "test" {
  vault_id          = dvls_vault.test.id
  name              = "tf_test_connection_string_upgrade"
  tags              = ["acceptance", "tf-test"]
  connection_string = "Server=foo"
}
`, testAccProviderConfig()))
}

func TestAccEntryCredentialAzureServicePrincipalResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_azure_service_principal.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_azure_sp_upgrade"
}

resource "dvls_entry_credential_azure_service_principal" "test" {
  vault_id      = dvls_vault.test.id
  name          = "tf_test_azure_sp_upgrade"
  tags          = ["acceptance", "tf-test"]
  client_id     = "client"
  client_secret = "secret"
  tenant_id     = "tenant"
}
`, testAccProviderConfig()))
}

func TestAccEntryCredentialUsernamePasswordResource_UpgradeFromV0(t *testing.T) {
	upgradeFromV0Test(t, "dvls_entry_credential_username_password.test", fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_username_password_upgrade"
}

resource "dvls_entry_credential_username_password" "test" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_username_password_upgrade"
  tags     = ["acceptance", "tf-test"]
  username = "user"
  password = "pass"
}
`, testAccProviderConfig()))
}
