package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var userGroupLookupKeys = []string{"id", "name"}

var _ datasource.DataSource = &UserGroupDataSource{}
var _ datasource.DataSourceWithConfigValidators = &UserGroupDataSource{}

func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

type UserGroupDataSource struct {
	client *dvls.Client
}

type UserGroupDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Type            types.String `tfsdk:"type"`
	IsAdministrator types.Bool   `tfsdk:"is_administrator"`
}

func (d *UserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *UserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS user group.",

		Attributes: map[string]schema.Attribute{
			"id":   principalIdAttribute("user group", userGroupLookupKeys),
			"name": lookupAttribute("The name of the user group.", userGroupLookupKeys),
			"description": schema.StringAttribute{
				Description: "The description of the user group.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf("The source of the user group. Possible values are %s.", listMapValues(userGroupTypes)),
				Computed:    true,
			},
			"is_administrator": schema.BoolAttribute{
				Description: "Whether members of the group are DVLS administrators.",
				Computed:    true,
			},
		},
	}
}

func (d *UserGroupDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return exactlyOneOfConfigValidators(userGroupLookupKeys...)
}

func (d *UserGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UserGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keys := []lookupKey[dvls.UserGroup]{
		{data.Id, func(group dvls.UserGroup) string { return group.Id }},
		{data.Name, func(group dvls.UserGroup) string { return group.Name }},
	}

	group, err := findPrincipal(ctx, d.client.UserGroups.ListWithContext, keys)
	if err != nil {
		resp.Diagnostics.AddError("unable to read user group", principalErrorDetail(err, "user group"))
		return
	}

	*data = UserGroupDataSourceModel{
		Id:              types.StringValue(group.Id),
		Name:            types.StringValue(group.Name),
		Description:     types.StringValue(group.Description),
		Type:            types.StringValue(userGroupTypes[group.Type]),
		IsAdministrator: types.BoolValue(group.IsAdministrator),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
