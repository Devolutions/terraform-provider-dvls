package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialConnectionStringEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialConnectionStringEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialConnectionStringEphemeralResource{}

func NewEntryCredentialConnectionStringEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialConnectionStringEphemeralResource{}
}

type EntryCredentialConnectionStringEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialConnectionStringEphemeralResourceModel = EntryCredentialConnectionStringDataSourceModel

func (e *EntryCredentialConnectionStringEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_connection_string"
}

func (e *EntryCredentialConnectionStringEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["connection_string"] = schema.StringAttribute{
		Description: "The entry credential connection string.",
		Computed:    true,
		Sensitive:   true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS Connection String Credential Entry, fetched ephemerally so the connection string never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialConnectionStringEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialConnectionStringEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypeConnectionString)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypeConnectionString)
		return
	}

	setEntryCredentialConnectionStringDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
