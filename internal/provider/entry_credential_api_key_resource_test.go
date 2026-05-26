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
			{
				Config: testAccEntryCredentialApiKeyResourceConfig(
					"tf_test_api_key", "tf_test_api_key", "test description", "tf_test_folder",
					"test-api-id", "test-api-key-secret", "test-tenant-id",
				),
				Check: testAccEntryCredentialApiKeyResourceCheck(
					"tf_test_api_key", "test description", "tf_test_folder",
					"test-api-id", "test-api-key-secret", "test-tenant-id",
				),
			},
			{
				Config: testAccEntryCredentialApiKeyResourceConfig(
					"tf_test_api_key", "tf_test_api_key_updated", "updated description", "tf_test_folder_updated",
					"updated-api-id", "updated-api-key-secret", "updated-tenant-id",
				),
				Check: testAccEntryCredentialApiKeyResourceCheck(
					"tf_test_api_key_updated", "updated description", "tf_test_folder_updated",
					"updated-api-id", "updated-api-key-secret", "updated-tenant-id",
				),
			},
			{
				ResourceName:      "dvls_entry_credential_api_key.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_credential_api_key.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEntryCredentialApiKeyResourceCheck(name, description, folder, apiId, apiKey, tenantId string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("dvls_entry_credential_api_key.test", "id"),
		resource.TestCheckResourceAttrPair("dvls_entry_credential_api_key.test", "vault_id", "dvls_vault.test", "id"),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "name", name),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "description", description),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "folder", folder),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tags.#", "2"),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_api_key.test", "tags.*", testAccTestTags[0]),
		resource.TestCheckTypeSetElemAttr("dvls_entry_credential_api_key.test", "tags.*", testAccTestTags[1]),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_id", apiId),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "api_key", apiKey),
		resource.TestCheckResourceAttr("dvls_entry_credential_api_key.test", "tenant_id", tenantId),
	)
}

func testAccEntryCredentialApiKeyResourceConfig(vaultName, name, description, folder, apiId, apiKey, tenantId string) string {
	return testAccEntryCredentialResourceConfig(
		"dvls_entry_credential_api_key",
		vaultName, name, description, folder,
		fmt.Sprintf(`  api_id = %q
  api_key = %q
  tenant_id = %q`, apiId, apiKey, tenantId),
	)
}
