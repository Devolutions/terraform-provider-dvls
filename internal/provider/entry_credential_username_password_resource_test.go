package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialUsernamePasswordResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCredentialUsernamePasswordResourceConfig(
					"tf_test_username_password", "tf_test_username_password", "test description", "tf_test_folder",
					"testuser", "testdomain", "testpassword",
				),
				Check: testAccEntryCredentialUsernamePasswordResourceCheck(
					"tf_test_username_password", "test description", "tf_test_folder",
					"testuser", "testdomain", "testpassword",
				),
			},
			{
				Config: testAccEntryCredentialUsernamePasswordResourceConfig(
					"tf_test_username_password", "tf_test_username_password_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updateddomain", "updatedpassword",
				),
				Check: testAccEntryCredentialUsernamePasswordResourceCheck(
					"tf_test_username_password_updated", "updated description", "tf_test_folder_updated",
					"updateduser", "updateddomain", "updatedpassword",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_username_password.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_username_password.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialUsernamePasswordResourceCheck(name, description, folder, username, domain, password string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_username_password.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_username_password.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_username_password.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_username_password.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "username", username),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "domain", domain),
		resource.TestCheckResourceAttr("dvls_entry_credential_username_password.test", "password", password),
	)
}

func testAccEntryCredentialUsernamePasswordResourceConfig(vaultName, name, description, folder, username, domain, password string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_username_password",
		vaultName, name, description, folder,
		fmt.Sprintf(`  username = %q
  domain = %q
  password = %q`, username, domain, password),
	)
}
