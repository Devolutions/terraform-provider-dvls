package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCredentialApiKeyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryCredentialDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccEntryCredentialApiKeyResourceConfig(
					"tf_test_api_key", "tf_test_api_key", "test description", "tf_test_folder",
					"test-api-id", "test-api-key-secret", "test-tenant-id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_credential_api_key.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_credential_api_key.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "name", "tf_test_api_key"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "description", "test description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "folder", "tf_test_folder"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tags.0", "tf-test"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tags.1", "acceptance"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_id", "test-api-id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_key", "test-api-key-secret"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tenant_id", "test-tenant-id"),
				),
			},
			// Update
			{
				Config: testAccEntryCredentialApiKeyResourceConfig(
					"tf_test_api_key", "tf_test_api_key_updated", "updated description", "tf_test_folder_updated",
					"updated-api-id", "updated-api-key-secret", "updated-tenant-id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "name", "tf_test_api_key_updated"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "description", "updated description"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_id", "updated-api-id"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_key", "updated-api-key-secret"),
					resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tenant_id", "updated-tenant-id"),
				),
			},
			// ImportState
			{
				ResourceName:      "dvls_entry_credential_api_key.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryCredentialImportStateIdFunc("dvls_entry_credential_api_key.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialApiKeyResourceConfig(vaultName, name, description, folder, apiId, apiKey, tenantId string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_credential_api_key" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = %[4]q
  folder      = %[5]q
  tags        = ["tf-test", "acceptance"]
  api_id      = %[6]q
  api_key     = %[7]q
  tenant_id   = %[8]q
}
`, testAccProviderConfig(), vaultName, name, description, folder, apiId, apiKey, tenantId)
}
