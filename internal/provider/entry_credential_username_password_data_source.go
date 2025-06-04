package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntryCredentialUsernamePasswordDataSource{}

func NewEntryCredentialUsernamePasswordDataSource() datasource.DataSource {
	return &EntryCredentialUsernamePasswordDataSource{}
}

// EntryCredentialUsernamePasswordDataSource defines the data source implementation.
type EntryCredentialUsernamePasswordDataSource struct {
	client *dvls.Client
}

// EntryCredentialUsernamePasswordDataSourceModel describes the data source data model.
type EntryCredentialUsernamePasswordDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	VaultId types.String `tfsdk:"vault_id"`
	Name    types.String `tfsdk:"name"`
	Folder  types.String `tfsdk:"folder"`

	// General
	Username types.String `tfsdk:"username"`
	Domain   types.String `tfsdk:"domain"`
	Password types.String `tfsdk:"password"`

	// More
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`
}

func (d *EntryCredentialUsernamePasswordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_username_password"
}

func (d *EntryCredentialUsernamePasswordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Username and Password Credential Entry",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the entry.",
				Required:    true,
				Validators:  []validator.String{entryIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "The ID of the vault.",
				Required:    true,
				Validators:  []validator.String{vaultIdValidator{}},
			},
			"name": schema.StringAttribute{
				Description: "The name of the entry.",
				Computed:    true,
			},
			"folder": schema.StringAttribute{
				Description: "The folder path of the entry.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The entry credential username.",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The entry credential domain.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The entry credential password.",
				Computed:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the entry.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of tags added to the entry.",
				Computed:    true,
			},
		},
	}
}

func (d *EntryCredentialUsernamePasswordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *EntryCredentialUsernamePasswordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryCredentialUsernamePasswordDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialUsernamePassword, err := d.client.Entries.Credential.GetById(data.VaultId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read username password credential entry", err.Error())
		return
	}

	setEntryCredentialUsernamePasswordDataModel(entryCredentialUsernamePassword, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
