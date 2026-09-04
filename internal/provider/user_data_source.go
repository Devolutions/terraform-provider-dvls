package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserDataSource{}
var _ datasource.DataSourceWithConfigValidators = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *dvls.Client
}

type UserDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	FullName           types.String `tfsdk:"full_name"`
	Email              types.String `tfsdk:"email"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	IsAdministrator    types.Bool   `tfsdk:"is_administrator"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	UserGroups         types.Set    `tfsdk:"user_groups"`
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS user. Application accounts are exposed by `dvls_application` instead.",

		Attributes: map[string]schema.Attribute{
			"id":        principalIdAttribute("user", userLookupKeys),
			"name":      lookupAttribute("The login name of the user.", userLookupKeys),
			"full_name": lookupAttribute("The display name of the user.", userLookupKeys),
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				Computed:    true,
			},
			"authentication_type": schema.StringAttribute{
				Description: fmt.Sprintf("How the user authenticates. Possible values are %s.", listMapValues(userAuthenticationTypes)),
				Computed:    true,
			},
			"is_administrator": schema.BoolAttribute{
				Description: "Whether the user is a DVLS administrator.",
				Computed:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the user is enabled.",
				Computed:    true,
			},
			"user_groups": schema.SetAttribute{
				Description: "IDs of the user groups the user belongs to.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *UserDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return exactlyOneOfConfigValidators(userLookupKeys...)
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := findPrincipal(ctx, d.client.Users.ListWithContext, userLookup(data.Id, data.Name, data.FullName))
	if err != nil {
		resp.Diagnostics.AddError("unable to read user", principalErrorDetail(err, "user"))
		return
	}

	*data = UserDataSourceModel{
		Id:                 types.StringValue(user.Id),
		Name:               types.StringValue(user.Name),
		FullName:           types.StringValue(user.FullName),
		Email:              types.StringValue(user.Email),
		AuthenticationType: types.StringValue(userAuthenticationTypes[user.AuthenticationType]),
		IsAdministrator:    types.BoolValue(user.IsAdministrator),
		IsEnabled:          types.BoolValue(user.IsEnabled),
		UserGroups:         tagsSliceToSet(user.UserGroups),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
