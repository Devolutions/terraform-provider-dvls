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
var _ resource.Resource = &EntryCredentialSecretResource{}
var _ resource.ResourceWithImportState = &EntryCredentialSecretResource{}

func NewEntryCredentialSecretResource() resource.Resource {
	return &EntryCredentialSecretResource{}
}

// EntryCredentialSecretResource defines the resource implementation.
type EntryCredentialSecretResource struct {
	client *dvls.Client
}

// EntryCredentialSecretResourceModel describes the resource data model.
type EntryCredentialSecretResourceModel struct {
	Id      types.String `tfsdk:"id"`
	VaultId types.String `tfsdk:"vault_id"`
	Name    types.String `tfsdk:"name"`
	Folder  types.String `tfsdk:"folder"`

	// General
	Secret types.String `tfsdk:"secret"`

	// More
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`
}

func (r *EntryCredentialSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_secret"
}

func (r *EntryCredentialSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Secret Credential Entry",

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
			"secret": schema.StringAttribute{
				Description: "The entry credential secret.",
				Optional:    true,
				Sensitive:   true,
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
		},
	}
}

func (r *EntryCredentialSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryCredentialSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialSecretResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSecret := newEntryCredentialSecretFromResourceModel(plan)

	entryCredentialSecretId, err := r.client.Entries.Credential.New(entryCredentialSecret)
	if err != nil {
		resp.Diagnostics.AddError("unable to create secret credential entry", err.Error())
		return
	}

	entryCredentialSecret, err = r.client.Entries.Credential.GetById(entryCredentialSecret.VaultId, entryCredentialSecretId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created secret credential entry", err.Error())
		return
	}

	setEntryCredentialSecretResourceModel(entryCredentialSecret, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialSecretResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSecret := newEntryCredentialSecretFromResourceModel(state)

	entryCredentialSecret, err := r.client.Entries.Credential.GetById(entryCredentialSecret.VaultId, entryCredentialSecret.Id)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read secret credential entry", err.Error())
		return
	}

	setEntryCredentialSecretResourceModel(entryCredentialSecret, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialSecretResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSecret := newEntryCredentialSecretFromResourceModel(plan)

	_, err := r.client.Entries.Credential.Update(entryCredentialSecret)
	if err != nil {
		resp.Diagnostics.AddError("unable to update secret credential entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialSecretResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialSecret := newEntryCredentialSecretFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialSecret)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete secret credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}
