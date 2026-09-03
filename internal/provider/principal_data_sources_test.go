package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	client := testAccLiveClient(t)
	user := testAccFirstOf(t, "user", client.Users.List)

	testAccPrincipalDataSourceTest(t, "dvls_user", user.Id, user.Name, "authentication_type")
}

func TestAccApplicationDataSource(t *testing.T) {
	client := testAccLiveClient(t)
	application := testAccFirstOf(t, "application", client.Users.ListApplications)

	testAccPrincipalDataSourceTest(t, "dvls_application", application.Id, application.Name, "is_enabled")
}

func TestAccUserGroupDataSource(t *testing.T) {
	client := testAccLiveClient(t)
	group := testAccFirstOf(t, "user group", client.UserGroups.List)

	testAccPrincipalDataSourceTest(t, "dvls_user_group", group.Id, group.Name, "type")
}

func testAccPrincipalDataSourceTest(t *testing.T, dataSourceType, id, name, extraAttribute string) {
	t.Helper()

	byId := "data." + dataSourceType + ".by_id"
	byName := "data." + dataSourceType + ".by_name"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() +
					testAccDataSourceBlock(dataSourceType, "by_id", fmt.Sprintf("id = %q", id)) +
					testAccDataSourceBlock(dataSourceType, "by_name", fmt.Sprintf("name = %q", name)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(byId, "id", id),
					resource.TestCheckResourceAttr(byId, "name", name),
					resource.TestCheckResourceAttrSet(byId, extraAttribute),
					resource.TestCheckResourceAttr(byName, "id", id),
					resource.TestCheckResourceAttrPair(byName, extraAttribute, byId, extraAttribute),
				),
			},
		},
	})
}
