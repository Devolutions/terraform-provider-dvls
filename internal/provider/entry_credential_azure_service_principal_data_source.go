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
var _ datasource.DataSource = &EntryCredentialAzureServicePrincipalDataSource{}

func NewEntryCredentialAzureServicePrincipalDataSource() datasource.DataSource {
	return &EntryCredentialAzureServicePrincipalDataSource{}
}

// EntryCredentialAzureServicePrincipalDataSource defines the data source implementation.
type EntryCredentialAzureServicePrincipalDataSource struct {
	client *dvls.Client
}

// EntryCredentialAzureServicePrincipalDataSourceModel describes the data source data model.
type EntryCredentialAzureServicePrincipalDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Folder      types.String   `tfsdk:"folder"`
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`

	// General
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	TenantId     types.String `tfsdk:"tenant_id"`
}

func (d *EntryCredentialAzureServicePrincipalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_azure_service_principal"
}

func (d *EntryCredentialAzureServicePrincipalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Azure Service Principal Credential Entry",

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
			"description": schema.StringAttribute{
				Description: "The description of the entry.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of tags added to the entry.",
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The entry credential client ID.",
				Computed:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The entry credential client secret.",
				Computed:    true,
				Sensitive:   true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The entry credential tenant ID.",
				Computed:    true,
			},
		},
	}
}

func (d *EntryCredentialAzureServicePrincipalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryCredentialAzureServicePrincipalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryCredentialAzureServicePrincipalDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialAzureServicePrincipal, err := d.client.Entries.Credential.GetById(data.VaultId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read azure service principal credential entry", err.Error())
		return
	}

	setEntryCredentialAzureServicePrincipalDataModel(entryCredentialAzureServicePrincipal, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
