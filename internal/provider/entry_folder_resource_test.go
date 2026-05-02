package provider

import (
	"fmt"
	"testing"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEntryFolderResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceConfig(
					"tf_test_folder_basic", "tf_test_root", "root description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_folder.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_folder.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "name", "tf_test_root"),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "description", "root description"),
					resource.TestCheckNoResourceAttr("dvls_entry_folder.test", "parent_folder"),
				),
			},
			{
				Config: testAccEntryFolderResourceConfig(
					"tf_test_folder_basic", "tf_test_root_renamed", "renamed description",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "name", "tf_test_root_renamed"),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "description", "renamed description"),
				),
			},
			{
				ResourceName:      "dvls_entry_folder.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_folder.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEntryFolderResource_nested(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceNestedConfig("tf_test_folder_nested"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_folder.parent", "name", "tf_test_parent"),
					resource.TestCheckNoResourceAttr("dvls_entry_folder.parent", "parent_folder"),
					resource.TestCheckResourceAttr("dvls_entry_folder.child", "name", "tf_test_child"),
					resource.TestCheckResourceAttr("dvls_entry_folder.child", "parent_folder", "tf_test_parent"),
				),
			},
		},
	})
}

func testAccCheckEntryFolderDestroy(s *terraform.State) error {
	client, err := getTestAccClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dvls_entry_folder" {
			continue
		}

		_, err := client.Entries.Folder.GetById(rs.Primary.Attributes["vault_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("folder %s/%s still exists", rs.Primary.Attributes["vault_id"], rs.Primary.ID)
		}

		if !dvls.IsNotFound(err) {
			return fmt.Errorf("unexpected error checking folder %s/%s: %s", rs.Primary.Attributes["vault_id"], rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccEntryFolderResourceConfig(vaultName, name, description string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "test" {
  vault_id    = dvls_vault.test.id
  name        = %[3]q
  description = %[4]q
}
`, testAccProviderConfig(), vaultName, name, description)
}

func testAccEntryFolderResourceNestedConfig(vaultName string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "parent" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_parent"
}

resource "dvls_entry_folder" "child" {
  vault_id      = dvls_vault.test.id
  name          = "tf_test_child"
  parent_folder = dvls_entry_folder.parent.name
}
`, testAccProviderConfig(), vaultName)
}

// DVLS does not support folder moves: PUT /entry silently ignores `path`
// changes. The resource declares parent_folder with RequiresReplace; this
// test confirms a parent_folder change plans as destroy-then-create.
func TestAccEntryFolderResource_parentFolderRequiresReplace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceMoveConfig("tf_test_folder_move", "alpha"),
			},
			{
				Config: testAccEntryFolderResourceMoveConfig("tf_test_folder_move", "beta"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("dvls_entry_folder.movable", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func testAccEntryFolderResourceMoveConfig(vaultName, parent string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "alpha" {
  vault_id = dvls_vault.test.id
  name     = "alpha"
}

resource "dvls_entry_folder" "beta" {
  vault_id = dvls_vault.test.id
  name     = "beta"
}

resource "dvls_entry_folder" "movable" {
  vault_id      = dvls_vault.test.id
  name          = "movable"
  parent_folder = %[3]q

  depends_on = [dvls_entry_folder.alpha, dvls_entry_folder.beta]
}
`, testAccProviderConfig(), vaultName, parent)
}
