package provider

import (
	"fmt"
	"testing"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdministrativeRoleDataSource(t *testing.T) {
	testAccPreCheckRoleAssignment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() + testAccAdministrativeRoleDataSourceBlock("by_id", fmt.Sprintf("id = %q", dvls.BuiltinRoleVaultUserId)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dvls_administrative_role.by_id", "id", dvls.BuiltinRoleVaultUserId),
					resource.TestCheckResourceAttrSet("data.dvls_administrative_role.by_id", "name"),
					resource.TestCheckResourceAttr("data.dvls_administrative_role.by_id", "is_built_in", "true"),
					resource.TestCheckTypeSetElemAttr("data.dvls_administrative_role.by_id", "supported_scopes.*", "vault"),
				),
			},
			{
				Config: testAccProviderConfig() +
					testAccAdministrativeRoleDataSourceBlock("by_id", fmt.Sprintf("id = %q", dvls.BuiltinRoleVaultUserId)) +
					testAccAdministrativeRoleDataSourceBlock("by_name", "name = data.dvls_administrative_role.by_id.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dvls_administrative_role.by_name", "id", dvls.BuiltinRoleVaultUserId),
					resource.TestCheckResourceAttrPair("data.dvls_administrative_role.by_name", "name", "data.dvls_administrative_role.by_id", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_administrative_role.by_name", "permissions.#", "data.dvls_administrative_role.by_id", "permissions.#"),
				),
			},
		},
	})
}
