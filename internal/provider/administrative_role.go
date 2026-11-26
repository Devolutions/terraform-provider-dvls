package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var administrativeRoleScopeTypes = map[dvls.AdministrativeRoleScopeType]string{
	dvls.AdministrativeRoleScopeWorkspace:    "global",
	dvls.AdministrativeRoleScopeBusinessUnit: "organizational_unit",
	dvls.AdministrativeRoleScopeVault:        "vault",
}

var administrativeRoleAssigneeTypes = map[dvls.AdministrativeRoleAssigneeType]string{
	dvls.AdministrativeRoleAssigneeUser:                "user",
	dvls.AdministrativeRoleAssigneeApplicationIdentity: "application",
	dvls.AdministrativeRoleAssigneeUserGroup:           "user_group",
}

func administrativeRoleErrorDetail(err error) string {
	switch {
	case dvls.IsNotFound(err):
		return "administrative roles require a DVLS release with the public API v2"
	case errors.Is(err, dvls.ErrMultipleAdministrativeRolesFound):
		return "more than one administrative role has that name, look it up by id instead"
	default:
		return err.Error()
	}
}

// DVLS answers a role assignment list filter naming a role, scope or assignee that no longer exists with a 400 on
// that field, where it used to return an empty list.
func isAdministrativeRoleAssignmentFilterGone(err error) bool {
	var reqErr *dvls.RequestError
	if !errors.As(err, &reqErr) || reqErr.StatusCode != http.StatusBadRequest {
		return false
	}

	var body struct {
		Errors []struct {
			Field string `json:"field"`
		} `json:"errors"`
	}
	if json.Unmarshal(reqErr.Body, &body) != nil || len(body.Errors) == 0 {
		return false
	}

	for _, fieldErr := range body.Errors {
		if !slices.Contains([]string{"roleId", "scopeId", "assigneeId"}, fieldErr.Field) {
			return false
		}
	}

	return true
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
