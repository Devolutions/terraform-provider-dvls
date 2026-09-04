package provider

import (
	"context"
	"errors"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var administrativeRoleScopeTypes = map[dvls.AdministrativeRoleScopeType]string{
	dvls.AdministrativeRoleScopeGlobal:             "global",
	dvls.AdministrativeRoleScopeOrganizationalUnit: "organizational_unit",
	dvls.AdministrativeRoleScopeVault:              "vault",
	dvls.AdministrativeRoleScopeGateway:            "gateway",
	dvls.AdministrativeRoleScopePamProvider:        "pam_provider",
}

var administrativeRoleAssigneeTypes = map[dvls.AdministrativeRoleAssigneeType]string{
	dvls.AdministrativeRoleAssigneeUser:        "user",
	dvls.AdministrativeRoleAssigneeApplication: "application",
	dvls.AdministrativeRoleAssigneeUserGroup:   "user_group",
}

func administrativeRoleErrorDetail(err error) string {
	switch {
	case dvls.IsNotFound(err):
		return "administrative roles require DVLS 2026.3 or later"
	case errors.Is(err, dvls.ErrMultipleAdministrativeRolesFound):
		return "more than one administrative role has that name, look it up by id instead"
	default:
		return err.Error()
	}
}

func resolveAdministrativeRoleId(ctx context.Context, client *dvls.Client, roleId, roleName types.String) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !roleId.IsNull() && !roleId.IsUnknown() {
		return roleId.ValueString(), diags
	}

	role, err := client.AdministrativeRoles.GetByNameWithContext(ctx, roleName.ValueString())
	if err != nil {
		diags.AddError("unable to look up administrative role by name", administrativeRoleErrorDetail(err))
		return "", diags
	}

	return role.Id, diags
}
