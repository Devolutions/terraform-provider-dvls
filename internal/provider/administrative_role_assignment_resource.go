package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &AdministrativeRoleAssignmentResource{}
var _ resource.ResourceWithImportState = &AdministrativeRoleAssignmentResource{}
var _ resource.ResourceWithConfigValidators = &AdministrativeRoleAssignmentResource{}
var _ resource.ResourceWithValidateConfig = &AdministrativeRoleAssignmentResource{}

func NewAdministrativeRoleAssignmentResource() resource.Resource {
	return &AdministrativeRoleAssignmentResource{}
}

type AdministrativeRoleAssignmentResource struct {
	client *dvls.Client
}

type AdministrativeRoleAssignmentResourceModel struct {
	Id           types.String `tfsdk:"id"`
	RoleId       types.String `tfsdk:"role_id"`
	RoleName     types.String `tfsdk:"role_name"`
	AssigneeId   types.String `tfsdk:"assignee_id"`
	AssigneeName types.String `tfsdk:"assignee_name"`
	AssigneeType types.String `tfsdk:"assignee_type"`
	ScopeType    types.String `tfsdk:"scope_type"`
	ScopeId      types.String `tfsdk:"scope_id"`
}

func (r *AdministrativeRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administrative_role_assignment"
}

func (r *AdministrativeRoleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A DVLS administrative role assignment: grants one administrative role to one principal (user, user group or application) on one scope.\n\n" +
			"Requires DVLS 2026.3 or later. Every attribute forces replacement.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The ID of the assignment. This is set by the provider after creation.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"role_id": schema.StringAttribute{
				Description: "The ID of the administrative role. Exactly one of `role_id` or `role_name` must be set.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{uuidValidator{"role id"}},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_name": schema.StringAttribute{
				Description:   "The name of the administrative role. Exactly one of `role_id` or `role_name` must be set.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignee_id": schema.StringAttribute{
				Description:   "The ID of the principal receiving the role (user, user group or application).",
				Required:      true,
				Validators:    []validator.String{uuidValidator{"assignee id"}},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignee_name": schema.StringAttribute{
				Description: "The display name of the principal, as reported by DVLS.",
				Computed:    true,
			},
			"assignee_type": schema.StringAttribute{
				Description: fmt.Sprintf("The type of the principal. Possible values are %s.", listMapValues(administrativeRoleAssigneeTypes)),
				Computed:    true,
			},
			"scope_type": schema.StringAttribute{
				Description:   fmt.Sprintf("The scope type of the assignment. Possible values are %s.", listMapValues(administrativeRoleScopeTypes)),
				Required:      true,
				Validators:    []validator.String{stringvalidator.OneOf(mapValues(administrativeRoleScopeTypes)...)},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"scope_id": schema.StringAttribute{
				Description:   "The ID of the scoped resource (vault, gateway, PAM provider or organizational unit). Required unless `scope_type` is `global`, in which case it must be omitted.",
				Optional:      true,
				Validators:    []validator.String{uuidValidator{"scope id"}},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *AdministrativeRoleAssignmentResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("role_id"),
			path.MatchRoot("role_name"),
		),
	}
}

func (r *AdministrativeRoleAssignmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config AdministrativeRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ScopeType.IsUnknown() || config.ScopeId.IsUnknown() {
		return
	}

	isGlobal := config.ScopeType.ValueString() == administrativeRoleScopeTypes[dvls.AdministrativeRoleScopeGlobal]
	if isGlobal && !config.ScopeId.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("scope_id"), "scope_id must not be set", "scope_id must be omitted when scope_type is global")
		return
	}
	if !isGlobal && config.ScopeId.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("scope_id"), "scope_id is required", fmt.Sprintf("scope_id is required when scope_type is %q", config.ScopeType.ValueString()))
	}
}

func (r *AdministrativeRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dvls.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *dvls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *AdministrativeRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *AdministrativeRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleId, diags := resolveAdministrativeRoleId(ctx, r.client, plan.RoleId, plan.RoleName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scopeType, err := lookupMapValue(administrativeRoleScopeTypes, plan.ScopeType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid scope type", err.Error())
		return
	}

	err = r.client.AdministrativeRoleAssignments.AddMemberWithContext(ctx, dvls.AdministrativeRoleMemberRequest{
		AdministrativeRoleId: roleId,
		AssigneeId:           plan.AssigneeId.ValueString(),
		ScopeType:            scopeType,
		ScopeResourceId:      plan.ScopeId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("unable to create administrative role assignment", administrativeRoleErrorDetail(err))
		return
	}

	member, found, diags := r.findMember(ctx, roleId, plan.ScopeType.ValueString(), plan.ScopeId.ValueString(), plan.AssigneeId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError(
			"administrative role assignment not found after creation",
			fmt.Sprintf("DVLS accepted the assignment but does not list %q as a member of role %q, verify the assignee id refers to an existing principal", plan.AssigneeId.ValueString(), roleId),
		)
		return
	}

	setAdministrativeRoleAssignmentResourceModel(member, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AdministrativeRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *AdministrativeRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	member, found, diags := r.findMember(ctx, state.RoleId.ValueString(), state.ScopeType.ValueString(), state.ScopeId.ValueString(), state.AssigneeId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	setAdministrativeRoleAssignmentResourceModel(member, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AdministrativeRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *AdministrativeRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AdministrativeRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *AdministrativeRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AdministrativeRoleAssignments.RemoveMemberWithContext(ctx, state.Id.ValueString())
	if err != nil {
		if dvls.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("unable to delete administrative role assignment", err.Error())
		return
	}
}

func (r *AdministrativeRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	roleId, scopeType, scopeId, assigneeId, err := parseAdministrativeRoleAssignmentImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	member, found, diags := r.findMember(ctx, roleId, scopeType, scopeId, assigneeId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("administrative role assignment not found", fmt.Sprintf("%q is not a member of role %q on the given scope", assigneeId, roleId))
		return
	}

	state := &AdministrativeRoleAssignmentResourceModel{ScopeType: types.StringValue(scopeType)}
	if scopeId != "" {
		state.ScopeId = types.StringValue(scopeId)
	}
	setAdministrativeRoleAssignmentResourceModel(member, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AdministrativeRoleAssignmentResource) findMember(ctx context.Context, roleId, scopeTypeValue, scopeId, assigneeId string) (dvls.AdministrativeRoleMember, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopeType, err := lookupMapValue(administrativeRoleScopeTypes, scopeTypeValue)
	if err != nil {
		diags.AddError("invalid scope type", fmt.Sprintf("unknown scope type %q, expected one of %s", scopeTypeValue, listMapValues(administrativeRoleScopeTypes)))
		return dvls.AdministrativeRoleMember{}, false, diags
	}

	members, err := r.client.AdministrativeRoleAssignments.GetMembersWithContext(ctx, roleId, scopeType, scopeId)
	if err != nil {
		diags.AddError("unable to read administrative role assignment", administrativeRoleErrorDetail(err))
		return dvls.AdministrativeRoleMember{}, false, diags
	}

	for _, member := range members {
		if member.AssigneeId == assigneeId {
			return member, true, diags
		}
	}

	return dvls.AdministrativeRoleMember{}, false, diags
}

func setAdministrativeRoleAssignmentResourceModel(member dvls.AdministrativeRoleMember, data *AdministrativeRoleAssignmentResourceModel) {
	data.Id = types.StringValue(member.Id)
	data.RoleId = types.StringValue(member.AdministrativeRoleId)
	data.AssigneeId = types.StringValue(member.AssigneeId)
	data.AssigneeName = types.StringValue(member.AssigneeName)
	data.AssigneeType = types.StringValue(administrativeRoleAssigneeTypes[member.AssigneeType])
}

func parseAdministrativeRoleAssignmentImportId(id string) (string, string, string, string, error) {
	parts := strings.SplitN(id, "/", 4)
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[3] == "" {
		return "", "", "", "", fmt.Errorf("unexpected format of ID (%s), expected <role_id>/<scope_type>/<scope_id>/<assignee_id> (leave <scope_id> empty for the global scope)", id)
	}

	return parts[0], parts[1], parts[2], parts[3], nil
}
