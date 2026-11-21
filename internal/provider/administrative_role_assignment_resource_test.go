package provider

import (
	"fmt"
	"testing"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAdministrativeRoleAssignmentResource_byRoleId(t *testing.T) {
	testAccPreCheckRoleAssignment(t)
	assigneeId := testAccFindAssigneeId(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAdministrativeRoleAssignmentDestroy,
			testAccCheckVaultDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeRoleAssignmentResourceConfig(
					"tf_test_role_assignment_id",
					fmt.Sprintf("role_id = %q", dvls.BuiltinRoleVaultUserId),
					assigneeId,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_administrative_role_assignment.test", "id"),
					resource.TestCheckResourceAttr("dvls_administrative_role_assignment.test", "role_id", dvls.BuiltinRoleVaultUserId),
					resource.TestCheckResourceAttr("dvls_administrative_role_assignment.test", "assignee_id", assigneeId),
					resource.TestCheckResourceAttrSet("dvls_administrative_role_assignment.test", "assignee_name"),
					resource.TestCheckResourceAttrSet("dvls_administrative_role_assignment.test", "assignee_type"),
					resource.TestCheckResourceAttr("dvls_administrative_role_assignment.test", "scope_type", "vault"),
					resource.TestCheckResourceAttrPair("dvls_administrative_role_assignment.test", "scope_id", "dvls_vault.test", "id"),
				),
			},
			{
				ResourceName:      "dvls_administrative_role_assignment.test",
				ImportState:       true,
				ImportStateIdFunc: testAccAdministrativeRoleAssignmentImportStateIdFunc("dvls_administrative_role_assignment.test"),
				ImportStateVerify: true,
			},
			{
				ResourceName:      "dvls_administrative_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAdministrativeRoleAssignmentResource_deletedOutOfBand(t *testing.T) {
	testAccPreCheckRoleAssignment(t)
	assigneeId := testAccFindAssigneeId(t)
	config := testAccAdministrativeRoleAssignmentResourceConfig(
		"tf_test_role_assignment_oob",
		fmt.Sprintf("role_id = %q", dvls.BuiltinRoleVaultUserId),
		assigneeId,
		"",
	)

	var firstId string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAdministrativeRoleAssignmentDestroy,
			testAccCheckVaultDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					firstId = s.RootModule().Resources["dvls_administrative_role_assignment.test"].Primary.ID
					return nil
				},
			},
			{
				PreConfig: func() {
					if err := testAccLiveClient(t).AdministrativeRoleAssignments.Delete(firstId); err != nil {
						t.Fatalf("unable to delete role assignment: %s", err)
					}
				},
				Config: config,
				Check: func(s *terraform.State) error {
					if id := s.RootModule().Resources["dvls_administrative_role_assignment.test"].Primary.ID; id == firstId {
						return fmt.Errorf("the assignment deleted out of band was not recreated, id is still %s", id)
					}
					return nil
				},
			},
		},
	})
}

func TestAccAdministrativeRoleAssignmentResource_byRoleName(t *testing.T) {
	testAccPreCheckRoleAssignment(t)
	assigneeId := testAccFindAssigneeId(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAdministrativeRoleAssignmentDestroy,
			testAccCheckVaultDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeRoleAssignmentResourceConfig(
					"tf_test_role_assignment_name",
					"role_name = data.dvls_administrative_role.vault_user.name",
					assigneeId,
					testAccDataSourceBlock("dvls_administrative_role", "vault_user", fmt.Sprintf("id = %q", dvls.BuiltinRoleVaultUserId)),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_administrative_role_assignment.test", "role_id", dvls.BuiltinRoleVaultUserId),
					resource.TestCheckResourceAttrPair("dvls_administrative_role_assignment.test", "role_name", "data.dvls_administrative_role.vault_user", "name"),
				),
			},
			{
				ResourceName:            "dvls_administrative_role_assignment.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccAdministrativeRoleAssignmentImportStateIdFunc("dvls_administrative_role_assignment.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role_name"},
			},
		},
	})
}

// TestAccAdministrativeRoleAssignmentResource_adoptsExisting creates the
// assignment outside Terraform first: DVLS refuses the duplicate and the
// resource adopts the existing assignment.
func TestAccAdministrativeRoleAssignmentResource_adoptsExisting(t *testing.T) {
	testAccPreCheckRoleAssignment(t)
	assigneeId := testAccFindAssigneeId(t)
	client := testAccLiveClient(t)

	vault, err := client.Vaults.New(dvls.Vault{Name: "tf_test_role_assignment_adopt"})
	if err != nil {
		t.Fatalf("unable to create vault: %s", err)
	}
	t.Cleanup(func() { _ = client.Vaults.Delete(vault.Id) })

	existing, err := client.AdministrativeRoleAssignments.New(dvls.AdministrativeRoleAssignmentRequest{
		RoleId:     dvls.BuiltinRoleVaultUserId,
		AssigneeId: assigneeId,
		ScopeType:  dvls.AdministrativeRoleScopeVault,
		ScopeId:    vault.Id,
	})
	if err != nil {
		t.Fatalf("unable to create role assignment: %s", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAdministrativeRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "dvls_administrative_role_assignment" "test" {
  role_id     = %q
  assignee_id = %q
  scope_type  = "vault"
  scope_id    = %q
}
`, dvls.BuiltinRoleVaultUserId, assigneeId, vault.Id),
				Check: resource.TestCheckResourceAttr("dvls_administrative_role_assignment.test", "id", existing.Id),
			},
		},
	})
}

// testAccAdministrativeRoleAssignmentResourceConfig builds provider + vault +
// a vault-scoped role assignment. `roleRef` is the `role_id = ...` or
// `role_name = ...` line; `extraBlocks` is appended verbatim (e.g. a data source).
func testAccAdministrativeRoleAssignmentResourceConfig(vaultName, roleRef, assigneeId, extraBlocks string) string {
	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_administrative_role_assignment" "test" {
  %[3]s
  assignee_id = %[4]q
  scope_type  = "vault"
  scope_id    = dvls_vault.test.id
}

%[5]s
`, testAccProviderConfig(), vaultName, roleRef, assigneeId, extraBlocks)
}

func testAccAdministrativeRoleAssignmentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		attrs := rs.Primary.Attributes
		return fmt.Sprintf("%s/%s/%s/%s", attrs["role_id"], attrs["scope_type"], attrs["scope_id"], attrs["assignee_id"]), nil
	}
}
