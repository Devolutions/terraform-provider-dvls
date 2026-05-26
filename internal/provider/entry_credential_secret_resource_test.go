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
			{
				Config: testAccEntryCredentialSecretResourceConfig(
					"tf_test_secret", "tf_test_secret", "test description", "tf_test_folder",
					"my-secret-value-123",
				),
				Check: testAccEntryCredentialSecretResourceCheck(
					"tf_test_secret", "test description", "tf_test_folder", "my-secret-value-123",
				),
			},
			{
				Config: testAccEntryCredentialSecretResourceConfig(
					"tf_test_secret", "tf_test_secret_updated", "updated description", "tf_test_folder_updated",
					"updated-secret-value-456",
				),
				Check: testAccEntryCredentialSecretResourceCheck(
					"tf_test_secret_updated", "updated description", "tf_test_folder_updated", "updated-secret-value-456",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_secret.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_secret.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialSecretResourceCheck(name, description, folder, secret string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_secret.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_secret.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_secret.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_secret.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_secret.test", "secret", secret),
	)
}

func testAccEntryCredentialSecretResourceConfig(vaultName, name, description, folder, secret string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_secret",
		vaultName, name, description, folder,
		fmt.Sprintf(`  secret = %q`, secret),
	)
}
