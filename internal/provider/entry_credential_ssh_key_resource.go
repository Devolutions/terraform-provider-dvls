package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EntryCredentialSSHKeyResource{}
var _ resource.ResourceWithImportState = &EntryCredentialSSHKeyResource{}

func NewEntryCredentialSSHKeyResource() resource.Resource {
	return &EntryCredentialSSHKeyResource{}
}

// EntryCredentialSSHKeyResource defines the resource implementation.
type EntryCredentialSSHKeyResource struct {
	client *dvls.Client
}

// EntryCredentialSSHKeyResourceModel describes the resource data model.
type EntryCredentialSSHKeyResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Folder      types.String   `tfsdk:"folder"`
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`

	// General
	Password       types.String `tfsdk:"password"`
	Passphrase     types.String `tfsdk:"passphrase"`
	PrivateKeyData types.String `tfsdk:"private_key_data"`
	PublicKey      types.String `tfsdk:"public_key"`
}

func (r *EntryCredentialSSHKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_ssh_key"
}

func (r *EntryCredentialSSHKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS SSH Key Credential Entry",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The ID of the entry. This is set by the provider after creation.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description:   "The ID of the vault.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description: "The name of the entry.",
				Required:    true,
			},
			"folder": schema.StringAttribute{
				Description: "The folder path where the entry is created.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the entry.",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "A list of tags to add to the entry.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The entry credential password.",
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase": schema.StringAttribute{
				Description: "The entry credential passphrase.",
				Optional:    true,
				Sensitive:   true,
			},
			"private_key_data": schema.StringAttribute{
				Description: "The entry credential private key.",
				Optional:    true,
				Sensitive:   true,
			},
			"public_key": schema.StringAttribute{
				Description: "The entry credential public key.",
				Optional:    true,
			},
		},
	}
}

func (r *EntryCredentialSSHKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *EntryCredentialSSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialSSHKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSSHKey := newEntryCredentialSSHKeyFromResourceModel(plan)

	entryCredentialSSHKeyId, err := r.client.Entries.Credential.New(entryCredentialSSHKey)
	if err != nil {
		resp.Diagnostics.AddError("unable to create SSH key credential entry", err.Error())
		return
	}

	entryCredentialSSHKey, err = r.client.Entries.Credential.GetById(entryCredentialSSHKey.VaultId, entryCredentialSSHKeyId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created SSH key credential entry", err.Error())
		return
	}

	setEntryCredentialSSHKeyResourceModel(entryCredentialSSHKey, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialSSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialSSHKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSSHKey := newEntryCredentialSSHKeyFromResourceModel(state)

	entryCredentialSSHKey, err := r.client.Entries.Credential.GetById(entryCredentialSSHKey.VaultId, entryCredentialSSHKey.Id)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read SSH key credential entry", err.Error())
		return
	}

	setEntryCredentialSSHKeyResourceModel(entryCredentialSSHKey, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialSSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialSSHKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSSHKey := newEntryCredentialSSHKeyFromResourceModel(plan)

	_, err := r.client.Entries.Credential.Update(entryCredentialSSHKey)
	if err != nil {
		resp.Diagnostics.AddError("unable to update SSH key credential entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialSSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialSSHKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSSHKey := newEntryCredentialSSHKeyFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialSSHKey)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete SSH key credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialSSHKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	entryCredentialSSHKey, err := r.client.Entries.Credential.GetById(vaultId, entryId)
	if err != nil {
		resp.Diagnostics.AddError("unable to read entry", err.Error())
		return
	}

	if entryCredentialSSHKey.Type != dvls.EntryCredentialType ||
		entryCredentialSSHKey.SubType != dvls.EntryCredentialSubTypePrivateKey {
		resp.Diagnostics.AddError("invalid entry type", "expected a SSH key credential entry.")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}
