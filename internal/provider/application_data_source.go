package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ApplicationDataSource{}
var _ datasource.DataSourceWithConfigValidators = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

type ApplicationDataSource struct {
	client *dvls.Client
}

type ApplicationDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	FullName        types.String `tfsdk:"full_name"`
	IsAdministrator types.Bool   `tfsdk:"is_administrator"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
	UserGroups      types.Set    `tfsdk:"user_groups"`
}

func (d *ApplicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *ApplicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS application account (application identity).",

		Attributes: map[string]schema.Attribute{
			"id":        principalIdAttribute("application", userLookupKeys),
			"name":      lookupAttribute("The application key.", userLookupKeys),
			"full_name": lookupAttribute("The display name of the application.", userLookupKeys),
			"is_administrator": schema.BoolAttribute{
				Description: "Whether the application is a DVLS administrator.",
				Computed:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the application is enabled.",
				Computed:    true,
			},
			"user_groups": schema.SetAttribute{
				Description: "IDs of the user groups the application belongs to.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *ApplicationDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return exactlyOneOfConfigValidators(userLookupKeys...)
}

func (d *ApplicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ApplicationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := findPrincipal(ctx, d.client.Users.ListApplicationsWithContext, userLookup(data.Id, data.Name, data.FullName))
	if err != nil {
		resp.Diagnostics.AddError("unable to read application", principalErrorDetail(err, "application"))
		return
	}

	*data = ApplicationDataSourceModel{
		Id:              types.StringValue(application.Id),
		Name:            types.StringValue(application.Name),
		FullName:        types.StringValue(application.FullName),
		IsAdministrator: types.BoolValue(application.IsAdministrator),
		IsEnabled:       types.BoolValue(application.IsEnabled),
		UserGroups:      tagsSliceToSet(application.UserGroups),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
