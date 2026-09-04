package provider

import (
	"fmt"
	"testing"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEntryPermissionsResource_override(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryPermissionsResourceConfig("tf_test_permissions_override", `
  role_override = "everyone"
  view_override = "everyone"
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("dvls_entry_permissions.test", "id", "dvls_entry_folder.test", "id"),
					resource.TestCheckResourceAttrPair("dvls_entry_permissions.test", "entry_id", "dvls_entry_folder.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "role_override", "everyone"),
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "view_override", "everyone"),
					resource.TestCheckNoResourceAttr("dvls_entry_permissions.test", "view_roles"),
					resource.TestCheckNoResourceAttr("dvls_entry_permissions.test", "permissions"),
				),
			},
			{
				Config: testAccEntryPermissionsResourceConfig("tf_test_permissions_override", `
  role_override = "never"
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "role_override", "never"),
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "view_override", "never"),
				),
			},
			{
				ResourceName:      "dvls_entry_permissions.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEntryPermissionsResourceConfig("tf_test_permissions_override", ""),
				Check:  testAccCheckEntryPermissionsReset("dvls_entry_folder.test"),
			},
		},
	})
}

func TestAccEntryPermissionsResource_custom(t *testing.T) {
	testAccPreCheckRoleAssignment(t)
	assigneeId := testAccFindAssigneeId(t)

	config := testAccEntryPermissionsResourceConfig("tf_test_permissions_custom", fmt.Sprintf(`
  role_override = "custom"
  view_override = "custom"
  view_roles    = [%[1]q]

  permissions = [
    {
      right = "edit"
      roles = [%[1]q]
    },
    {
      right    = "checkout"
      override = "everyone"
    },
  ]
`, assigneeId))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEntryFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "role_override", "custom"),
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "view_override", "custom"),
					resource.TestCheckTypeSetElemAttr("dvls_entry_permissions.test", "view_roles.*", assigneeId),
					resource.TestCheckResourceAttr("dvls_entry_permissions.test", "permissions.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("dvls_entry_permissions.test", "permissions.*", map[string]string{
						"right":    "edit",
						"override": "custom",
						"roles.#":  "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("dvls_entry_permissions.test", "permissions.*", map[string]string{
						"right":    "checkout",
						"override": "everyone",
					}),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
			{
				ResourceName:      "dvls_entry_permissions.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccEntryPermissionsResourceConfig builds provider + vault + folder and,
// when `attributes` is non-empty, a dvls_entry_permissions resource on the
// folder carrying those attribute lines.
func testAccEntryPermissionsResourceConfig(vaultName, attributes string) string {
	permissions := ""
	if attributes != "" {
		permissions = fmt.Sprintf(`
resource "dvls_entry_permissions" "test" {
  entry_id = dvls_entry_folder.test.id
%s}
`, attributes)
	}

	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "test" {
  vault_id = dvls_vault.test.id
  name     = %[3]q
}

%[4]s
`, testAccProviderConfig(), vaultName, testAccTestFolder, permissions)
}

func testAccCheckEntryPermissionsReset(folderResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[folderResourceName]
		if !ok {
			return fmt.Errorf("not found: %s", folderResourceName)
		}

		client, err := getTestAccClient()
		if err != nil {
			return err
		}

		security, err := client.Entries.Permissions.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to read permissions of %s: %s", rs.Primary.ID, err)
		}

		if security.RoleOverride != dvls.SecurityRoleOverrideDefault {
			return fmt.Errorf("expected permissions of %s to be reset, got role override %s", rs.Primary.ID, security.RoleOverride)
		}

		return nil
	}
}
