package provider

import (
	"fmt"
	"regexp"
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
				Config: testAccEntryFolderResourceNestedConfig("tf_test_folder_nested", "tf_test_child"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_folder.parent", "name", "tf_test_parent"),
					resource.TestCheckNoResourceAttr("dvls_entry_folder.parent", "parent_folder"),
					resource.TestCheckResourceAttr("dvls_entry_folder.child", "name", "tf_test_child"),
					resource.TestCheckResourceAttr("dvls_entry_folder.child", "parent_folder", "tf_test_parent"),
				),
			},
			{
				Config: testAccEntryFolderResourceNestedConfig("tf_test_folder_nested", "tf_test_child_renamed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("dvls_entry_folder.child", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_folder.child", "name", "tf_test_child_renamed"),
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

func testAccEntryFolderResourceNestedConfig(vaultName, childName string) string {
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
  name          = %[3]q
  parent_folder = dvls_entry_folder.parent.name
}
`, testAccProviderConfig(), vaultName, childName)
}

// A parent_folder change moves the folder in place: same id, content follows.
// A move into a folder that does not exist fails and leaves the folder where it was.
func TestAccEntryFolderResource_move(t *testing.T) {
	var folderId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceMoveConfig("alpha"),
				Check:  testAccStoreResourceId("dvls_entry_folder.movable", &folderId),
			},
			{
				Config: testAccEntryFolderResourceMoveConfig("beta"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("dvls_entry_folder.movable", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_folder.movable", "parent_folder", "beta"),
					testAccCheckResourceId("dvls_entry_folder.movable", &folderId),
					testAccCheckEntryFolderPath("dvls_entry_folder.child", `beta\movable`),
				),
			},
			{
				Config:      testAccEntryFolderResourceMoveConfig("gamma"),
				ExpectError: regexp.MustCompile(`(?s)parent folder not found.*Declare "gamma" as its own\s+dvls_entry_folder`),
			},
			{
				Config:   testAccEntryFolderResourceMoveConfig("beta"),
				PlanOnly: true,
			},
			{
				Config: testAccEntryFolderResourceMoveConfig("beta"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckResourceId("dvls_entry_folder.movable", &folderId),
					testAccCheckEntryFolderPath("dvls_entry_folder.movable", "beta"),
				),
			},
		},
	})
}

func testAccStoreResourceId(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckResourceId(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID != *id {
			return fmt.Errorf("%s id changed from %s to %s", resourceName, *id, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckEntryFolderPath(resourceName, wantPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client, err := getTestAccClient()
		if err != nil {
			return err
		}

		folder, err := client.Entries.Folder.GetById(rs.Primary.Attributes["vault_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to read folder %s: %s", rs.Primary.ID, err)
		}
		if folder.Path != wantPath {
			return fmt.Errorf("folder %s has path %q, want %q", rs.Primary.ID, folder.Path, wantPath)
		}

		return nil
	}
}

// A folder update sends no credentials, and DVLS keeps the ones stored on the folder.
func TestAccEntryFolderResource_keepsCredentials(t *testing.T) {
	var folderId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceConfig("tf_test_folder_credentials", "tf_test_creds", "before"),
				Check:  testAccStoreResourceId("dvls_entry_folder.test", &folderId),
			},
			{
				PreConfig: func() {
					client, err := getTestAccClient()
					if err != nil {
						t.Fatal(err)
					}
					folder, err := client.Entries.Folder.GetById(testAccFolderVaultId(t, client, "tf_test_folder_credentials"), folderId)
					if err != nil {
						t.Fatalf("unable to read folder: %s", err)
					}
					folder.Data = &dvls.EntryFolderData{Domain: "keep.local", Username: "keepuser"}
					if _, err := client.Entries.Folder.Update(folder); err != nil {
						t.Fatalf("unable to set folder credentials: %s", err)
					}
				},
				Config: testAccEntryFolderResourceConfig("tf_test_folder_credentials", "tf_test_creds_renamed", "after"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckResourceId("dvls_entry_folder.test", &folderId),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "name", "tf_test_creds_renamed"),
					func(s *terraform.State) error {
						client, err := getTestAccClient()
						if err != nil {
							return err
						}
						folder, err := client.Entries.Folder.GetById(s.RootModule().Resources["dvls_entry_folder.test"].Primary.Attributes["vault_id"], folderId)
						if err != nil {
							return err
						}
						data, ok := folder.GetFolderData()
						if !ok || data.Domain != "keep.local" || data.Username != "keepuser" {
							return fmt.Errorf("the rename cleared the folder credentials: %+v", data)
						}
						return nil
					},
				),
			},
		},
	})
}

// A folder deleted outside Terraform, in a vault the provider can still read, is gone: it is created again.
func TestAccEntryFolderResource_deletedOutOfBand(t *testing.T) {
	var folderId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceConfig("tf_test_folder_oob", "tf_test_oob", "oob"),
				Check:  testAccStoreResourceId("dvls_entry_folder.test", &folderId),
			},
			{
				PreConfig: func() {
					client, err := getTestAccClient()
					if err != nil {
						t.Fatal(err)
					}
					if err := client.Entries.Folder.DeleteById(testAccFolderVaultId(t, client, "tf_test_folder_oob"), folderId); err != nil {
						t.Fatalf("unable to delete folder: %s", err)
					}
				},
				Config: testAccEntryFolderResourceConfig("tf_test_folder_oob", "tf_test_oob", "oob"),
				Check: func(s *terraform.State) error {
					if id := s.RootModule().Resources["dvls_entry_folder.test"].Primary.ID; id == folderId {
						return fmt.Errorf("the folder deleted out of band was not created again, id is still %s", id)
					}
					return nil
				},
			},
		},
	})
}

func testAccFolderVaultId(t *testing.T, client *dvls.Client, vaultName string) string {
	t.Helper()

	vault, err := client.Vaults.GetByName(vaultName)
	if err != nil {
		t.Fatalf("unable to find vault %s: %s", vaultName, err)
	}

	return vault.Id
}

// TestAccEntryFolderResource_byVaultName exercises the vault_name path
// (resolved via dvls.Vaults.GetByName) instead of the direct vault_id input.
func TestAccEntryFolderResource_byVaultName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryFolderResourceByVaultNameConfig("tf_test_folder_by_vault_name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("dvls_entry_folder.test", "vault_id", "dvls_vault.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "vault_name", "tf_test_folder_by_vault_name"),
					resource.TestCheckResourceAttr("dvls_entry_folder.test", "name", "tf_test_root"),
				),
			},
		},
	})
}

func testAccEntryFolderResourceByVaultNameConfig(vaultName string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "test" {
  vault_name = dvls_vault.test.name
  name       = "tf_test_root"
}
`, testAccProviderConfig(), vaultName)
}

func testAccEntryFolderResourceMoveConfig(parent string) string {
	return fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = "tf_test_folder_move"
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
  parent_folder = %[2]q

  depends_on = [dvls_entry_folder.alpha, dvls_entry_folder.beta]
}

# The child folder ignores its own parent_folder so the check sees where DVLS put it after the move.
resource "dvls_entry_folder" "child" {
  vault_id      = dvls_vault.test.id
  name          = "child"
  parent_folder = "${dvls_entry_folder.movable.parent_folder}\\${dvls_entry_folder.movable.name}"

  lifecycle {
    ignore_changes = [parent_folder]
  }
}
`, testAccProviderConfig(), parent)
}
