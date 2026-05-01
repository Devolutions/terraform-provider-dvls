package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialSSHKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialSSHKeyEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialSSHKeyEphemeralResource{}

func NewEntryCredentialSSHKeyEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialSSHKeyEphemeralResource{}
}

type EntryCredentialSSHKeyEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialSSHKeyEphemeralResourceModel = EntryCredentialSSHKeyDataSourceModel

func (e *EntryCredentialSSHKeyEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_ssh_key"
}

func (e *EntryCredentialSSHKeyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["username"] = schema.StringAttribute{
		Description: "The entry credential username.",
		Computed:    true,
	}
	attrs["password"] = schema.StringAttribute{
		Description: "The entry credential password.",
		Computed:    true,
		Sensitive:   true,
	}
	attrs["passphrase"] = schema.StringAttribute{
		Description: "The entry credential passphrase.",
		Computed:    true,
		Sensitive:   true,
	}
	attrs["private_key_data"] = schema.StringAttribute{
		Description: "The entry credential private key data.",
		Computed:    true,
		Sensitive:   true,
	}
	attrs["public_key"] = schema.StringAttribute{
		Description: "The entry credential public key data.",
		Computed:    true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS SSH Key Credential Entry, fetched ephemerally so the private key never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialSSHKeyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialSSHKeyEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypePrivateKey)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypePrivateKey)
		return
	}

	setEntryCredentialSSHKeyDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
