package provider

import (
	"context"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

var _ ephemeral.EphemeralResource = &EntryCredentialUsernamePasswordEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &EntryCredentialUsernamePasswordEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigValidators = &EntryCredentialUsernamePasswordEphemeralResource{}

func NewEntryCredentialUsernamePasswordEphemeralResource() ephemeral.EphemeralResource {
	return &EntryCredentialUsernamePasswordEphemeralResource{}
}

type EntryCredentialUsernamePasswordEphemeralResource struct {
	credentialEphemeralBase
}

type EntryCredentialUsernamePasswordEphemeralResourceModel = EntryCredentialUsernamePasswordDataSourceModel

func (e *EntryCredentialUsernamePasswordEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_username_password"
}

func (e *EntryCredentialUsernamePasswordEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	attrs := credentialEphemeralCommonAttributes()
	attrs["username"] = schema.StringAttribute{
		Description: "The entry credential username.",
		Computed:    true,
	}
	attrs["domain"] = schema.StringAttribute{
		Description: "The entry credential domain.",
		Computed:    true,
	}
	attrs["password"] = schema.StringAttribute{
		Description: "The entry credential password.",
		Computed:    true,
		Sensitive:   true,
	}
	resp.Schema = schema.Schema{
		Description: "A DVLS Username and Password Credential Entry, fetched ephemerally so the password never lands in Terraform state.",
		Attributes:  attrs,
	}
}

func (e *EntryCredentialUsernamePasswordEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *EntryCredentialUsernamePasswordEphemeralResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := fetchCredentialEntry(e.client, data.VaultId, data.Id, data.Name, data.Folder, dvls.EntryCredentialSubTypeDefault)
	if err != nil {
		appendCredentialFetchError(&resp.Diagnostics, err, data.Name, dvls.EntryCredentialSubTypeDefault)
		return
	}

	setEntryCredentialUsernamePasswordDataModel(entry, data)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
