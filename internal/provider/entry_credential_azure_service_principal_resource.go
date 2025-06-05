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
var _ resource.Resource = &EntryCredentialAzureServicePrincipalResource{}
var _ resource.ResourceWithImportState = &EntryCredentialAzureServicePrincipalResource{}

func NewEntryCredentialAzureServicePrincipalResource() resource.Resource {
	return &EntryCredentialAzureServicePrincipalResource{}
}

// EntryCredentialAzureServicePrincipalResource defines the resource implementation.
type EntryCredentialAzureServicePrincipalResource struct {
	client *dvls.Client
}

// EntryCredentialAzureServicePrincipalResourceModel describes the resource data model.
type EntryCredentialAzureServicePrincipalResourceModel struct {
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

func (r *EntryCredentialAzureServicePrincipalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_azure_service_principal"
}

func (r *EntryCredentialAzureServicePrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Azure Service Principal Credential Entry",

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
			"client_id": schema.StringAttribute{
				Description: "The entry credential client ID.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The entry credential client secret.",
				Optional:    true,
				Sensitive:   true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The entry credential tenant ID.",
				Optional:    true,
			},
		},
	}
}

func (r *EntryCredentialAzureServicePrincipalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryCredentialAzureServicePrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialAzureServicePrincipalResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialAzureServicePrincipal := newEntryCredentialAzureServicePrincipalFromResourceModel(plan)

	entryCredentialAzureServicePrincipalId, err := r.client.Entries.Credential.New(entryCredentialAzureServicePrincipal)
	if err != nil {
		resp.Diagnostics.AddError("unable to create azure service principal credential entry", err.Error())
		return
	}

	entryCredentialAzureServicePrincipal, err = r.client.Entries.Credential.GetById(entryCredentialAzureServicePrincipal.VaultId, entryCredentialAzureServicePrincipalId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created azure service principal credential entry", err.Error())
		return
	}

	setEntryCredentialAzureServicePrincipalResourceModel(entryCredentialAzureServicePrincipal, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialAzureServicePrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialAzureServicePrincipalResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialAzureServicePrincipal := newEntryCredentialAzureServicePrincipalFromResourceModel(state)

	entryCredentialAzureServicePrincipal, err := r.client.Entries.Credential.GetById(entryCredentialAzureServicePrincipal.VaultId, entryCredentialAzureServicePrincipal.Id)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read azure service principal credential entry", err.Error())
		return
	}

	setEntryCredentialAzureServicePrincipalResourceModel(entryCredentialAzureServicePrincipal, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialAzureServicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialAzureServicePrincipalResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialAzureServicePrincipal := newEntryCredentialAzureServicePrincipalFromResourceModel(plan)

	_, err := r.client.Entries.Credential.Update(entryCredentialAzureServicePrincipal)
	if err != nil {
		resp.Diagnostics.AddError("unable to update azure service principal credential entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialAzureServicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialAzureServicePrincipalResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialAzureServicePrincipal := newEntryCredentialAzureServicePrincipalFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialAzureServicePrincipal)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete azure service principal credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialAzureServicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}
