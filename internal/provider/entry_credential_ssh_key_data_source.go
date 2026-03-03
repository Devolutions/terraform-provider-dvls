package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EntryCredentialSSHKeyDataSource{}
var _ datasource.DataSourceWithConfigValidators = &EntryCredentialSSHKeyDataSource{}

func NewEntryCredentialSSHKeyDataSource() datasource.DataSource {
	return &EntryCredentialSSHKeyDataSource{}
}

// EntryCredentialSSHKeyDataSource defines the data source implementation.
type EntryCredentialSSHKeyDataSource struct {
	client *dvls.Client
}

// EntryCredentialSSHKeyDataSourceModel describes the data source data model.
type EntryCredentialSSHKeyDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Folder      types.String   `tfsdk:"folder"`
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`

	// General
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Passphrase     types.String `tfsdk:"passphrase"`
	PrivateKeyData types.String `tfsdk:"private_key_data"`
	PublicKey      types.String `tfsdk:"public_key"`
}

func (d *EntryCredentialSSHKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_ssh_key"
}

func (d *EntryCredentialSSHKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS SSH Key Credential Entry",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the entry.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{entryIdValidator{}},
			},
			"vault_id": schema.StringAttribute{
				Description: "The ID of the vault.",
				Required:    true,
				Validators:  []validator.String{vaultIdValidator{}},
			},
			"name": schema.StringAttribute{
				Description: "The name of the entry.",
				Optional:    true,
				Computed:    true,
			},
			"folder": schema.StringAttribute{
				Description: "The folder path to search in. Returns entries in the specified folder and all sub-folders.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.AlsoRequires(path.MatchRoot("name"))},
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
			"username": schema.StringAttribute{
				Description: "The entry credential username.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The entry credential password.",
				Computed:    true,
				Sensitive:   true,
			},
			"passphrase": schema.StringAttribute{
				Description: "The entry credential passphrase.",
				Computed:    true,
				Sensitive:   true,
			},
			"private_key_data": schema.StringAttribute{
				Description: "The entry credential private key data.",
				Computed:    true,
				Sensitive:   true,
			},
			"public_key": schema.StringAttribute{
				Description: "The entry credential public key data.",
				Computed:    true,
			},
		},
	}
}

func (d *EntryCredentialSSHKeyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *EntryCredentialSSHKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EntryCredentialSSHKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EntryCredentialSSHKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		if !data.Name.IsNull() || !data.Folder.IsNull() {
			resp.Diagnostics.AddWarning("id takes precedence", "When id is provided, name and folder are ignored.")
		}
	}

	entry, err := fetchCredentialEntry(d.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypePrivateKey)
	if err != nil {
		if errors.Is(err, dvls.ErrMultipleEntriesFound) {
			resp.Diagnostics.AddError("multiple entries found", fmt.Sprintf("more than one entry named %q found, use id to target the correct one", data.Name.ValueString()))
			return
		}
		resp.Diagnostics.AddError("unable to read SSH key credential entry", err.Error())
		return
	}

	setEntryCredentialSSHKeyDataModel(entry, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
