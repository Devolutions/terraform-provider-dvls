package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EntryPermissionsResource{}
var _ resource.ResourceWithImportState = &EntryPermissionsResource{}
var _ resource.ResourceWithValidateConfig = &EntryPermissionsResource{}

func NewEntryPermissionsResource() resource.Resource {
	return &EntryPermissionsResource{}
}

type EntryPermissionsResource struct {
	client *dvls.Client
}

func (r *EntryPermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_permissions"
}

func (r *EntryPermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	overrideValues := mapValues(securityRoleOverrides)
	overrideDescription := listMapValues(securityRoleOverrides)

	resp.Schema = schema.Schema{
		MarkdownDescription: "The permission block of a DVLS entry or folder. Declare at most one `dvls_entry_permissions` per entry: creating it takes ownership of the permissions already set on the entry, and destroying it resets them to `default`.\n\n" +
			"The `view` right is not part of `permissions`; it is controlled by `view_override` and `view_roles`. " +
			"`permissions` and `view_roles` can only be set when `role_override` is `custom` or `custom_inherited`; for any other value DVLS forces `view_override` to `role_override`.\n\n" +
			"DVLS has no partial update for permissions: the provider fetches the entry and saves it back, so a concurrent change to the same entry between the two requests is overwritten.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The ID of the entry the permissions apply to.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"entry_id": schema.StringAttribute{
				Description:   "The ID of the entry or folder.",
				Required:      true,
				Validators:    []validator.String{uuidValidator{"entry id"}},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"role_override": schema.StringAttribute{
				Description: fmt.Sprintf("How the entry permissions are resolved. Possible values are %s. Defaults to `default`.", overrideDescription),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(securityRoleOverrides[dvls.SecurityRoleOverrideDefault]),
				Validators:  []validator.String{stringvalidator.OneOf(overrideValues...)},
			},
			"view_override": schema.StringAttribute{
				Description: fmt.Sprintf("How the `view` right is resolved. Possible values are %s. Must match `role_override` unless `role_override` is `custom` or `custom_inherited`.", overrideDescription),
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.OneOf(overrideValues...)},
			},
			"view_roles": schema.SetAttribute{
				Description: "IDs of the principals (users, user groups or applications) granted the `view` right. Requires `view_override` to be `custom` or `custom_inherited`.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(uuidValidator{"principal id"}),
				},
			},
			"permissions": schema.SetNestedAttribute{
				Description: "Rights granted on the entry, one block per right.",
				Optional:    true,
				Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"right": schema.StringAttribute{
							Description: fmt.Sprintf("The right being granted. Possible values are %s.", formatValues(entryPermissionRightValues)),
							Required:    true,
							Validators:  []validator.String{stringvalidator.OneOf(entryPermissionRightValues...)},
						},
						"override": schema.StringAttribute{
							Description: fmt.Sprintf("How the right is resolved. Possible values are %s. Defaults to `custom`.", overrideDescription),
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(securityRoleOverrides[dvls.SecurityRoleOverrideCustom]),
							Validators:  []validator.String{stringvalidator.OneOf(overrideValues...)},
						},
						"roles": schema.SetAttribute{
							Description: "IDs of the principals (users, user groups or applications) granted the right. Requires `override` to be `custom` or `custom_inherited`.",
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(uuidValidator{"principal id"}),
							},
						},
					},
				},
			},
		},
	}
}

func (r *EntryPermissionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config EntryPermissionsResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleOverride := securityRoleOverrides[dvls.SecurityRoleOverrideDefault]
	if !config.RoleOverride.IsNull() {
		roleOverride = config.RoleOverride.ValueString()
	}

	if !config.RoleOverride.IsUnknown() && !isCustomSecurityRoleOverride(roleOverride) {
		if !config.Permissions.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "permissions require a custom role_override",
				fmt.Sprintf("permissions can only be set when role_override is custom or custom_inherited, got %q", roleOverride))
		}
		if !config.ViewRoles.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("view_roles"), "view_roles require a custom role_override",
				fmt.Sprintf("view_roles can only be set when role_override is custom or custom_inherited, got %q", roleOverride))
		}
		if !config.ViewOverride.IsNull() && !config.ViewOverride.IsUnknown() && config.ViewOverride.ValueString() != roleOverride {
			resp.Diagnostics.AddAttributeError(path.Root("view_override"), "view_override must match role_override",
				fmt.Sprintf("DVLS forces view_override to role_override (%q) unless role_override is custom or custom_inherited", roleOverride))
		}
	}

	if !config.ViewRoles.IsNull() && !config.ViewOverride.IsUnknown() &&
		(config.ViewOverride.IsNull() || !isCustomSecurityRoleOverride(config.ViewOverride.ValueString())) {
		resp.Diagnostics.AddAttributeError(path.Root("view_override"), "view_roles require a custom view_override",
			"view_override must be custom or custom_inherited when view_roles is set")
	}

	if config.Permissions.IsNull() || config.Permissions.IsUnknown() {
		return
	}

	var permissions []entryPermissionModel
	resp.Diagnostics.Append(config.Permissions.ElementsAs(ctx, &permissions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	seen := map[string]struct{}{}
	for _, permission := range permissions {
		if !permission.Roles.IsNull() && !permission.Override.IsUnknown() &&
			!permission.Override.IsNull() && !isCustomSecurityRoleOverride(permission.Override.ValueString()) {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "permission roles require a custom override",
				fmt.Sprintf("roles of the %q permission can only be set when its override is custom or custom_inherited", permission.Right.ValueString()))
		}

		if permission.Right.IsUnknown() {
			continue
		}
		right := permission.Right.ValueString()
		if _, ok := seen[right]; ok {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "duplicate permission right",
				fmt.Sprintf("the %q right is declared more than once", right))
		}
		seen[right] = struct{}{}
	}
}

func (r *EntryPermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryPermissionsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryPermissionsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	security, err := r.client.Entries.Permissions.GetWithContext(ctx, state.EntryId.ValueString())
	if err != nil {
		if isEntryNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read entry permissions", err.Error())
		return
	}

	resp.Diagnostics.Append(setEntryPermissionsResourceModel(ctx, security, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryPermissionsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryPermissionsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Entries.Permissions.SetWithContext(ctx, state.EntryId.ValueString(), dvls.EntrySecurity{})
	if err != nil {
		if isEntryNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("unable to reset entry permissions", err.Error())
		return
	}
}

func (r *EntryPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	_, err := uuid.Parse(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", fmt.Sprintf("expected the entry id as a UUID, got %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("entry_id"), req.ID)...)
}

func (r *EntryPermissionsResource) apply(ctx context.Context, plan *EntryPermissionsResourceModel) diag.Diagnostics {
	security, diags := newEntrySecurityFromResourceModel(ctx, plan)
	if diags.HasError() {
		return diags
	}

	entryId := plan.EntryId.ValueString()

	err := r.client.Entries.Permissions.SetWithContext(ctx, entryId, security)
	if err != nil {
		diags.AddError("unable to save entry permissions", err.Error())
		return diags
	}

	saved, err := r.client.Entries.Permissions.GetWithContext(ctx, entryId)
	if err != nil {
		diags.AddError("unable to fetch saved entry permissions", err.Error())
		return diags
	}

	if dropped := droppedPrincipals(security, saved); len(dropped) > 0 {
		diags.AddError(
			"unknown principals in entry permissions",
			fmt.Sprintf("DVLS ignored the following principal ids, verify they refer to existing users, user groups or applications: %s", strings.Join(dropped, ", ")),
		)
		return diags
	}

	diags.Append(setEntryPermissionsResourceModel(ctx, saved, plan)...)

	return diags
}

func isEntryNotFound(err error) bool {
	return errors.Is(err, dvls.ErrEntryNotFound) || dvls.IsNotFound(err)
}
