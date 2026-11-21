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
var _ resource.ResourceWithUpgradeState = &EntryPermissionsResource{}

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
		Version: 1,
		MarkdownDescription: "The permission block of a DVLS entry or folder. Declare at most one `dvls_entry_permissions` per entry: creating it takes ownership of the permissions already set on the entry, and destroying it resets them to `default`.\n\n" +
			"The `view` right is set in `permissions` like any other right. " +
			"`permissions` can only be set when `role_override` is `custom` or `custom_inherited`; for any other value DVLS applies `role_override` to every right.\n\n" +
			"Saving replaces every permission of the entry: a permission set outside Terraform is overwritten.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The ID of the entry the permissions apply to.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description: "The ID of the vault holding the entry or folder.",
				Required:    true,
				Validators:  []validator.String{uuidValidator{"vault id"}},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
						// A state upgraded from version 0 has no vault id: setting it is an in-place update.
						resp.RequiresReplace = !req.StateValue.IsNull()
					}, "Changing an existing vault_id forces replacement.", "Changing an existing `vault_id` forces replacement."),
				},
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

	if state.VaultId.IsNull() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	security, err := r.client.Entries.Permissions.GetWithContext(ctx, state.VaultId.ValueString(), state.EntryId.ValueString())
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

	if state.VaultId.IsNull() {
		resp.Diagnostics.AddError("unable to reset entry permissions",
			"the state has no vault_id (it was created by an earlier provider version): set vault_id and apply before destroying, or remove the resource from the state with terraform state rm")
		return
	}

	err := r.client.Entries.Permissions.SetWithContext(ctx, state.VaultId.ValueString(), state.EntryId.ValueString(), dvls.EntrySecurity{})
	if err != nil {
		if isEntryNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("unable to reset entry permissions", err.Error())
		return
	}
}

func (r *EntryPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err == nil {
		_, err = uuid.Parse(vaultId)
	}
	if err == nil {
		_, err = uuid.Parse(entryId)
	}
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", fmt.Sprintf("expected <vault_id>/<entry_id> with both ids as UUIDs, got %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), entryId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("entry_id"), entryId)...)
}

func (r *EntryPermissionsResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	permissionObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"right":    schema.StringAttribute{Required: true},
			"override": schema.StringAttribute{Optional: true, Computed: true},
			"roles":    schema.SetAttribute{Optional: true, ElementType: types.StringType},
		},
	}

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":            schema.StringAttribute{Computed: true},
					"entry_id":      schema.StringAttribute{Required: true},
					"role_override": schema.StringAttribute{Optional: true, Computed: true},
					"view_override": schema.StringAttribute{Optional: true, Computed: true},
					"view_roles":    schema.SetAttribute{Optional: true, ElementType: types.StringType},
					"permissions":   schema.SetNestedAttribute{Optional: true, NestedObject: permissionObject},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var prior entryPermissionsResourceModelV0
				resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
				if resp.Diagnostics.HasError() {
					return
				}

				upgraded, diags := upgradeEntryPermissionsStateV0(ctx, prior)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
			},
		},
	}
}

func (r *EntryPermissionsResource) apply(ctx context.Context, plan *EntryPermissionsResourceModel) diag.Diagnostics {
	security, diags := newEntrySecurityFromResourceModel(ctx, plan)
	if diags.HasError() {
		return diags
	}

	vaultId := plan.VaultId.ValueString()
	entryId := plan.EntryId.ValueString()

	err := r.client.Entries.Permissions.SetWithContext(ctx, vaultId, entryId, security)
	if err != nil {
		diags.AddError("unable to save entry permissions", err.Error())
		return diags
	}

	saved, err := r.client.Entries.Permissions.GetWithContext(ctx, vaultId, entryId)
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
