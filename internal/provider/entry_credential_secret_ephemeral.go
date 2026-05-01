package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialSecretEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialSecretEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialSecretEphemeralResource{}

func NewEntryCredentialSecretEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialSecretEphemeralResource{}
}

type EntryCredentialSecretEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialSecretEphemeralResourceModel = EntryCredentialSecretDataSourceModel

func (e *EntryCredentialSecretEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_secret"
}

func (e *EntryCredentialSecretEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["secret"] = schema.StringAttribute{
		Description: "The entry credential secret.",
		Computed:    true,
		Sensitive:   true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS Secret Credential Entry, fetched ephemerally so the secret never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialSecretEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialSecretEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypeAccessCode)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypeAccessCode)
		return
	}

	setEntryCredentialSecretDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
