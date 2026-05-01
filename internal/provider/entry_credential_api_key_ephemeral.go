package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialApiKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialApiKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialApiKeyEphemeralResource{}

func NewEntryCredentialApiKeyEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialApiKeyEphemeralResource{}
}

type EntryCredentialApiKeyEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialApiKeyEphemeralResourceModel = EntryCredentialApiKeyDataSourceModel

func (e *EntryCredentialApiKeyEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_api_key"
}

func (e *EntryCredentialApiKeyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["api_id"] = schema.StringAttribute{
		Description: "The entry credential API ID.",
		Computed:    true,
	}
	attrs["api_key"] = schema.StringAttribute{
		Description: "The entry credential API key.",
		Computed:    true,
		Sensitive:   true,
	}
	attrs["tenant_id"] = schema.StringAttribute{
		Description: "The entry credential tenant ID.",
		Computed:    true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS API Key Credential Entry, fetched ephemerally so the API key never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialApiKeyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialApiKeyEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypeApiKey)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypeApiKey)
		return
	}

	setEntryCredentialApiKeyDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
