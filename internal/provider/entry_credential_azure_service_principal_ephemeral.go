package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialAzureServicePrincipalEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialAzureServicePrincipalEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialAzureServicePrincipalEphemeralResource{}

func NewEntryCredentialAzureServicePrincipalEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialAzureServicePrincipalEphemeralResource{}
}

type EntryCredentialAzureServicePrincipalEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialAzureServicePrincipalEphemeralResourceModel = EntryCredentialAzureServicePrincipalDataSourceModel

func (e *EntryCredentialAzureServicePrincipalEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_azure_service_principal"
}

func (e *EntryCredentialAzureServicePrincipalEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["client_id"] = schema.StringAttribute{
		Description: "The entry credential client ID.",
		Computed:    true,
	}
	attrs["client_secret"] = schema.StringAttribute{
		Description: "The entry credential client secret.",
		Computed:    true,
		Sensitive:   true,
	}
	attrs["tenant_id"] = schema.StringAttribute{
		Description: "The entry credential tenant ID.",
		Computed:    true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS Azure Service Principal Credential Entry, fetched ephemerally so the client secret never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialAzureServicePrincipalEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialAzureServicePrincipalEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypeAzureServicePrincipal)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypeAzureServicePrincipal)
		return
	}

	setEntryCredentialAzureServicePrincipalDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
