package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialSecretResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			testAccVaultWithFoldersStep("tf_test_secret", "tf_test_folder", "tf_test_folder_updated"),
			// Create
			{
				Config: testAccEntryCredentialSecretResourceConfig(
					"tf_test_secret", "tf_test_secret", "test description", "tf_test_folder",
					"my-secret-value-123",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_secret.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_secret.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "name", "tf_test_secret"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "tags.0", "acceptance"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "tags.1", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "secret", "my-secret-value-123"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialSecretResourceConfig(
					"tf_test_secret", "tf_test_secret_updated", "updated description", "tf_test_folder_updated",
					"updated-secret-value-456",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "name", "tf_test_secret_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "secret", "updated-secret-value-456"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_secret.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryCredentialImportStateIdFunc("dvls_entry_credential_secret.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialSecretResourceConfig(vaultName, name, description, folder, secret string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_secret" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = %[4]q
  folder      = %[5]q
  tags        = ["acceptance", "tf-test"]
  secret      = %[6]q
}
`, testAccProviderConfig(), vaultName, name, description, folder, secret)
}
