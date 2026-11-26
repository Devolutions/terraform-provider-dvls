package provider

import (
	"context"
	"slices"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var securityRoleOverrides = map[dvls.SecurityRoleOverride]string{
	dvls.SecurityRoleOverrideDefault:         "default",
	dvls.SecurityRoleOverrideCustom:          "custom",
	dvls.SecurityRoleOverrideInherited:       "inherited",
	dvls.SecurityRoleOverrideEveryone:        "everyone",
	dvls.SecurityRoleOverrideNever:           "never",
	dvls.SecurityRoleOverrideCustomInherited: "custom_inherited",
}

var securityRoleRights = map[dvls.SecurityRoleRight]string{
	dvls.SecurityRoleRightView:                              "view",
	dvls.SecurityRoleRightViewPassword:                      "view_password",
	dvls.SecurityRoleRightAdd:                               "add",
	dvls.SecurityRoleRightDelete:                            "delete",
	dvls.SecurityRoleRightEdit:                              "edit",
	dvls.SecurityRoleRightEditStatus:                        "edit_status",
	dvls.SecurityRoleRightEditDescription:                   "edit_description",
	dvls.SecurityRoleRightEditSecurity:                      "edit_security",
	dvls.SecurityRoleRightPasswordHistory:                   "password_history",
	dvls.SecurityRoleRightEntryHistory:                      "connection_history",
	dvls.SecurityRoleRightRemoteTools:                       "remote_tools",
	dvls.SecurityRoleRightAttachment:                        "attachment",
	dvls.SecurityRoleRightEditAttachment:                    "edit_attachment",
	dvls.SecurityRoleRightInventory:                         "inventory",
	dvls.SecurityRoleRightViewLogs:                          "view_logs",
	dvls.SecurityRoleRightDocumentation:                     "handbook",
	dvls.SecurityRoleRightEditDocumentation:                 "edit_handbook",
	dvls.SecurityRoleRightWebManagementTools:                "web_management_tools",
	dvls.SecurityRoleRightConsoleManagementTools:            "console_management_tools",
	dvls.SecurityRoleRightMacroScriptTools:                  "macro_script_tools",
	dvls.SecurityRoleRightMacroScriptToolsEntry:             "macro_script_tools_entry",
	dvls.SecurityRoleRightEditPassword:                      "edit_password",
	dvls.SecurityRoleRightExecute:                           "execute",
	dvls.SecurityRoleRightViewSessionRecording:              "view_session_recording",
	dvls.SecurityRoleRightViewInformation:                   "view_information",
	dvls.SecurityRoleRightExport:                            "export",
	dvls.SecurityRoleRightEditInformation:                   "edit_information",
	dvls.SecurityRoleRightMove:                              "move",
	dvls.SecurityRoleRightDeleteDocumentation:               "delete_handbook",
	dvls.SecurityRoleRightViewSensitiveInformation:          "view_sensitive_information",
	dvls.SecurityRoleRightResetPassword:                     "reset_password",
	dvls.SecurityRoleRightApproveCheckoutRequest:            "approve_checkout_request",
	dvls.SecurityRoleRightForceCheckin:                      "force_checkin",
	dvls.SecurityRoleRightCheckout:                          "checkout",
	dvls.SecurityRoleRightReadLogs:                          "read_logs",
	dvls.SecurityRoleRightEditVPNSSHGatewayConfiguration:    "edit_vpn_ssh_gateway_configuration",
	dvls.SecurityRoleRightEditSessionRecordingConfiguration: "edit_session_recording_configuration",
}

var entryPermissionRightValues = mapValues(securityRoleRights)

var entryPermissionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"right":    types.StringType,
		"override": types.StringType,
		"roles":    types.SetType{ElemType: types.StringType},
	},
}

func isCustomSecurityRoleOverride(value string) bool {
	return value == securityRoleOverrides[dvls.SecurityRoleOverrideCustom] ||
		value == securityRoleOverrides[dvls.SecurityRoleOverrideCustomInherited]
}

type EntryPermissionsResourceModel struct {
	Id           types.String `tfsdk:"id"`
	VaultId      types.String `tfsdk:"vault_id"`
	EntryId      types.String `tfsdk:"entry_id"`
	RoleOverride types.String `tfsdk:"role_override"`
	Permissions  types.Set    `tfsdk:"permissions"`
}

type entryPermissionModel struct {
	Right    types.String `tfsdk:"right"`
	Override types.String `tfsdk:"override"`
	Roles    types.Set    `tfsdk:"roles"`
}

func newEntrySecurityFromResourceModel(ctx context.Context, data *EntryPermissionsResourceModel) (dvls.EntrySecurity, diag.Diagnostics) {
	var diags diag.Diagnostics

	roleOverride, err := lookupMapValue(securityRoleOverrides, data.RoleOverride.ValueString())
	if err != nil {
		diags.AddError("invalid role_override", err.Error())
		return dvls.EntrySecurity{}, diags
	}

	security := dvls.EntrySecurity{RoleOverride: roleOverride}

	if data.Permissions.IsNull() || data.Permissions.IsUnknown() {
		return security, diags
	}

	var permissions []entryPermissionModel
	diags.Append(data.Permissions.ElementsAs(ctx, &permissions, false)...)
	if diags.HasError() {
		return dvls.EntrySecurity{}, diags
	}

	for _, permission := range permissions {
		right, err := lookupMapValue(securityRoleRights, permission.Right.ValueString())
		if err != nil {
			diags.AddError("invalid permission right", err.Error())
			return dvls.EntrySecurity{}, diags
		}

		override, err := lookupMapValue(securityRoleOverrides, permission.Override.ValueString())
		if err != nil {
			diags.AddError("invalid permission override", err.Error())
			return dvls.EntrySecurity{}, diags
		}

		security.Permissions = append(security.Permissions, dvls.EntryPermission{
			Right:    right,
			Override: override,
			Roles:    tagsSetToSlice(permission.Roles),
		})
	}

	return security, diags
}

func setEntryPermissionsResourceModel(ctx context.Context, security dvls.EntrySecurity, data *EntryPermissionsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Outside the custom modes DVLS ignores the permissions sent and answers a single view right mirroring the mode.
	permissionsSet := types.SetNull(entryPermissionObjectType)
	if isCustomSecurityRoleOverride(securityRoleOverrides[security.RoleOverride]) && len(security.Permissions) > 0 {
		permissions := make([]entryPermissionModel, 0, len(security.Permissions))
		for _, permission := range security.Permissions {
			permissions = append(permissions, entryPermissionModel{
				Right:    types.StringValue(securityRoleRights[permission.Right]),
				Override: types.StringValue(securityRoleOverrides[permission.Override]),
				Roles:    tagsSliceToSet(permission.Roles),
			})
		}

		permissionsSet, diags = types.SetValueFrom(ctx, entryPermissionObjectType, permissions)
		if diags.HasError() {
			return diags
		}
	}

	data.Id = data.EntryId
	data.RoleOverride = types.StringValue(securityRoleOverrides[security.RoleOverride])
	data.Permissions = permissionsSet

	return diags
}

func droppedPrincipals(requested, actual dvls.EntrySecurity) []string {
	actualRights := map[dvls.SecurityRoleRight][]string{}
	for _, permission := range actual.Permissions {
		actualRights[permission.Right] = permission.Roles
	}

	var dropped []string
	for _, permission := range requested.Permissions {
		dropped = append(dropped, missingValues(permission.Roles, actualRights[permission.Right])...)
	}

	slices.Sort(dropped)
	return slices.Compact(dropped)
}

type entryPermissionsResourceModelV0 struct {
	Id           types.String `tfsdk:"id"`
	EntryId      types.String `tfsdk:"entry_id"`
	RoleOverride types.String `tfsdk:"role_override"`
	ViewOverride types.String `tfsdk:"view_override"`
	ViewRoles    types.Set    `tfsdk:"view_roles"`
	Permissions  types.Set    `tfsdk:"permissions"`
}

// upgradeEntryPermissionsStateV0 folds view_override and view_roles into a view permission, unless the view is left
// to its default, and drops the sealed right. A v0 state has no vault id: it stays null until the next apply sets it from the configuration.
func upgradeEntryPermissionsStateV0(ctx context.Context, prior entryPermissionsResourceModelV0) (EntryPermissionsResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var permissions []entryPermissionModel
	if !prior.Permissions.IsNull() && !prior.Permissions.IsUnknown() {
		diags.Append(prior.Permissions.ElementsAs(ctx, &permissions, false)...)
		if diags.HasError() {
			return EntryPermissionsResourceModel{}, diags
		}
	}

	permissions = slices.DeleteFunc(permissions, func(permission entryPermissionModel) bool {
		return permission.Right.ValueString() == "sealed"
	})

	// DVLS does not store a view right left to its default without roles: such a view permission would not read back.
	defaultView := prior.ViewOverride.ValueString() == securityRoleOverrides[dvls.SecurityRoleOverrideDefault] &&
		len(prior.ViewRoles.Elements()) == 0
	if isCustomSecurityRoleOverride(prior.RoleOverride.ValueString()) && !prior.ViewOverride.IsNull() && !prior.ViewOverride.IsUnknown() && !defaultView {
		permissions = append(permissions, entryPermissionModel{
			Right:    types.StringValue(securityRoleRights[dvls.SecurityRoleRightView]),
			Override: prior.ViewOverride,
			Roles:    prior.ViewRoles,
		})
	}

	permissionsSet := types.SetNull(entryPermissionObjectType)
	if len(permissions) > 0 {
		var setDiags diag.Diagnostics
		permissionsSet, setDiags = types.SetValueFrom(ctx, entryPermissionObjectType, permissions)
		diags.Append(setDiags...)
		if diags.HasError() {
			return EntryPermissionsResourceModel{}, diags
		}
	}

	return EntryPermissionsResourceModel{
		Id:           prior.Id,
		VaultId:      types.StringNull(),
		EntryId:      prior.EntryId,
		RoleOverride: prior.RoleOverride,
		Permissions:  permissionsSet,
	}, diags
}

func missingValues(requested, actual []string) []string {
	var missing []string
	for _, value := range requested {
		if !slices.Contains(actual, value) {
			missing = append(missing, value)
		}
	}
	return missing
}
