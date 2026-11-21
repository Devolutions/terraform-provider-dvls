package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUpgradeEntryPermissionsStateV0(t *testing.T) {
	ctx := context.Background()
	principal := "00000000-0000-0000-0000-000000000001"

	permissionSet := func(t *testing.T, permissions ...entryPermissionModel) types.Set {
		t.Helper()
		if len(permissions) == 0 {
			return types.SetNull(entryPermissionObjectType)
		}
		set, diags := types.SetValueFrom(ctx, entryPermissionObjectType, permissions)
		if diags.HasError() {
			t.Fatalf("unable to build permissions: %v", diags)
		}
		return set
	}
	permission := func(right, override string, roles ...string) entryPermissionModel {
		return entryPermissionModel{Right: types.StringValue(right), Override: types.StringValue(override), Roles: tagsSliceToSet(roles)}
	}

	tests := []struct {
		name  string
		prior func(t *testing.T) entryPermissionsResourceModelV0
		want  func(t *testing.T) types.Set
	}{
		{
			name: "custom view moves into permissions and sealed is dropped",
			prior: func(t *testing.T) entryPermissionsResourceModelV0 {
				return entryPermissionsResourceModelV0{
					RoleOverride: types.StringValue("custom"),
					ViewOverride: types.StringValue("custom"),
					ViewRoles:    tagsSliceToSet([]string{principal}),
					Permissions:  permissionSet(t, permission("edit", "custom", principal), permission("sealed", "everyone")),
				}
			},
			want: func(t *testing.T) types.Set {
				return permissionSet(t, permission("edit", "custom", principal), permission("view", "custom", principal))
			},
		},
		{
			name: "custom_inherited role with a non custom view",
			prior: func(t *testing.T) entryPermissionsResourceModelV0 {
				return entryPermissionsResourceModelV0{
					RoleOverride: types.StringValue("custom_inherited"),
					ViewOverride: types.StringValue("everyone"),
					ViewRoles:    types.SetNull(types.StringType),
					Permissions:  types.SetNull(entryPermissionObjectType),
				}
			},
			want: func(t *testing.T) types.Set {
				return permissionSet(t, permission("view", "everyone"))
			},
		},
		{
			name: "a default view is not folded into permissions",
			prior: func(t *testing.T) entryPermissionsResourceModelV0 {
				return entryPermissionsResourceModelV0{
					RoleOverride: types.StringValue("custom"),
					ViewOverride: types.StringValue("default"),
					ViewRoles:    types.SetNull(types.StringType),
					Permissions:  permissionSet(t, permission("edit", "custom", principal)),
				}
			},
			want: func(t *testing.T) types.Set {
				return permissionSet(t, permission("edit", "custom", principal))
			},
		},
		{
			name: "non custom role keeps no permissions",
			prior: func(t *testing.T) entryPermissionsResourceModelV0 {
				return entryPermissionsResourceModelV0{
					RoleOverride: types.StringValue("everyone"),
					ViewOverride: types.StringValue("everyone"),
					ViewRoles:    types.SetNull(types.StringType),
					Permissions:  types.SetNull(entryPermissionObjectType),
				}
			},
			want: func(t *testing.T) types.Set { return types.SetNull(entryPermissionObjectType) },
		},
		{
			name: "only a sealed permission leaves none",
			prior: func(t *testing.T) entryPermissionsResourceModelV0 {
				return entryPermissionsResourceModelV0{
					RoleOverride: types.StringValue("custom"),
					ViewOverride: types.StringNull(),
					ViewRoles:    types.SetNull(types.StringType),
					Permissions:  permissionSet(t, permission("sealed", "never")),
				}
			},
			want: func(t *testing.T) types.Set { return types.SetNull(entryPermissionObjectType) },
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prior := tc.prior(t)
			prior.Id = types.StringValue(principal)
			prior.EntryId = types.StringValue(principal)

			got, diags := upgradeEntryPermissionsStateV0(ctx, prior)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !got.VaultId.IsNull() {
				t.Fatalf("vault_id = %s, want null", got.VaultId)
			}
			if !got.Id.Equal(prior.Id) || !got.EntryId.Equal(prior.EntryId) || !got.RoleOverride.Equal(prior.RoleOverride) {
				t.Fatalf("id, entry_id or role_override changed: %+v", got)
			}
			if want := tc.want(t); !got.Permissions.Equal(want) {
				t.Fatalf("permissions = %s, want %s", got.Permissions, want)
			}
		})
	}
}
