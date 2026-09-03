package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AdministrativeRoleDataSource{}
var _ datasource.DataSourceWithConfigValidators = &AdministrativeRoleDataSource{}

func NewAdministrativeRoleDataSource() datasource.DataSource {
	return &AdministrativeRoleDataSource{}
}

type AdministrativeRoleDataSource struct {
	client *dvls.Client
}

type AdministrativeRoleDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Permissions     types.List   `tfsdk:"permissions"`
	SupportedScopes types.List   `tfsdk:"supported_scopes"`
	IsAssignable    types.Bool   `tfsdk:"is_assignable"`
	IsBuiltIn       types.Bool   `tfsdk:"is_built_in"`
	IsPam           types.Bool   `tfsdk:"is_pam"`
	IsPrivileged    types.Bool   `tfsdk:"is_privileged"`
}

func (d *AdministrativeRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administrative_role"
}

func (d *AdministrativeRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS administrative role definition. Requires DVLS 2026.3 or later.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the administrative role. Exactly one of `id` or `name` must be set.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{uuidValidator{"role id"}},
			},
			"name": schema.StringAttribute{
				Description: "The name of the administrative role. Exactly one of `id` or `name` must be set.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the administrative role.",
				Computed:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "The administrative permissions granted by the role.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"supported_scopes": schema.ListAttribute{
				Description: fmt.Sprintf("The scope types the role can be assigned on. Possible values are %s.", listMapValues(administrativeRoleScopeTypes)),
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_assignable": schema.BoolAttribute{
				Description: "Whether the role can be assigned to principals.",
				Computed:    true,
			},
			"is_built_in": schema.BoolAttribute{
				Description: "Whether the role is a DVLS built-in role.",
				Computed:    true,
			},
			"is_pam": schema.BoolAttribute{
				Description: "Whether the role is a PAM role.",
				Computed:    true,
			},
			"is_privileged": schema.BoolAttribute{
				Description: "Whether the role is privileged.",
				Computed:    true,
			},
		},
	}
}

func (d *AdministrativeRoleDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return exactlyOneOfConfigValidators("id", "name")
}

func (d *AdministrativeRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dvls.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dvls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *AdministrativeRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AdministrativeRoleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var role dvls.AdministrativeRole
	var err error

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		role, err = d.client.AdministrativeRoles.GetWithContext(ctx, data.Id.ValueString())
	} else {
		role, err = d.client.AdministrativeRoles.GetByNameWithContext(ctx, data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("unable to read administrative role", administrativeRoleErrorDetail(err))
		return
	}

	resp.Diagnostics.Append(setAdministrativeRoleDataModel(ctx, role, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setAdministrativeRoleDataModel(ctx context.Context, role dvls.AdministrativeRole, data *AdministrativeRoleDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	permissions := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.String())
	}

	scopes := make([]string, 0, len(role.SupportedScopes))
	for _, scope := range role.SupportedScopes {
		scopes = append(scopes, administrativeRoleScopeTypes[scope])
	}

	permissionsList, permissionDiags := types.ListValueFrom(ctx, types.StringType, permissions)
	diags.Append(permissionDiags...)
	scopesList, scopeDiags := types.ListValueFrom(ctx, types.StringType, scopes)
	diags.Append(scopeDiags...)
	if diags.HasError() {
		return diags
	}

	*data = AdministrativeRoleDataSourceModel{
		Id:              types.StringValue(role.Id),
		Name:            types.StringValue(role.Name),
		Description:     types.StringValue(role.Description),
		Permissions:     permissionsList,
		SupportedScopes: scopesList,
		IsAssignable:    types.BoolValue(role.IsAssignable),
		IsBuiltIn:       types.BoolValue(role.IsBuiltIn),
		IsPam:           types.BoolValue(role.IsPam),
		IsPrivileged:    types.BoolValue(role.IsPrivileged),
	}

	return diags
}
