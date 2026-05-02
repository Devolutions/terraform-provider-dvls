package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EntryCredentialApiKeyResource{}
var _ resource.ResourceWithImportState = &EntryCredentialApiKeyResource{}
var _ resource.ResourceWithUpgradeState = &EntryCredentialApiKeyResource{}

func NewEntryCredentialApiKeyResource() resource.Resource {
	return &EntryCredentialApiKeyResource{}
}

// EntryCredentialApiKeyResource defines the resource implementation.
type EntryCredentialApiKeyResource struct {
	client *dvls.Client
}

// EntryCredentialApiKeyResourceModel describes the resource data model.
type EntryCredentialApiKeyResourceModel struct {
	Id          types.String `tfsdk:"id"`
	VaultId     types.String `tfsdk:"vault_id"`
	Name        types.String `tfsdk:"name"`
	Folder      types.String `tfsdk:"folder"`
	Description types.String `tfsdk:"description"`
	Tags        types.Set    `tfsdk:"tags"`

	// General
	ApiId    types.String `tfsdk:"api_id"`
	ApiKey   types.String `tfsdk:"api_key"`
	TenantId types.String `tfsdk:"tenant_id"`
}

func (r *EntryCredentialApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_api_key"
}

func (r *EntryCredentialApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "A DVLS API Key Credential Entry",

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
			"tags": schema.SetAttribute{
				ElementType:   types.StringType,
				Description:   "A list of tags to add to the entry.",
				Optional:      true,
				PlanModifiers: []planmodifier.Set{emptyTagsToNull()},
			},
			"api_id": schema.StringAttribute{
				Description: "The entry credential API ID.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The entry credential API key.",
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

func (r *EntryCredentialApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dvls.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *dvls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *EntryCredentialApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialApiKey := newEntryCredentialApiKeyFromResourceModel(plan)

	entryCredentialApiKeyId, err := r.client.Entries.Credential.New(entryCredentialApiKey)
	if err != nil {
		resp.Diagnostics.AddError("unable to create api key credential entry", err.Error())
		return
	}

	entryCredentialApiKey, err = r.client.Entries.Credential.GetById(entryCredentialApiKey.VaultId, entryCredentialApiKeyId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created api key credential entry", err.Error())
		return
	}

	setEntryCredentialApiKeyResourceModel(entryCredentialApiKey, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialApiKey := newEntryCredentialApiKeyFromResourceModel(state)

	entryCredentialApiKey, err := r.client.Entries.Credential.GetById(entryCredentialApiKey.VaultId, entryCredentialApiKey.Id)
	if err != nil {
		if dvls.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read api key credential entry", err.Error())
		return
	}

	setEntryCredentialApiKeyResourceModel(entryCredentialApiKey, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialApiKey := newEntryCredentialApiKeyFromResourceModel(plan)

	entryCredentialApiKey, err := r.client.Entries.Credential.Update(entryCredentialApiKey)
	if err != nil {
		resp.Diagnostics.AddError("unable to update api key credential entry", err.Error())
		return
	}

	setEntryCredentialApiKeyResourceModel(entryCredentialApiKey, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialApiKey := newEntryCredentialApiKeyFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialApiKey)
	if err != nil {
		if dvls.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete api key credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	entryCredentialApiKey, err := r.client.Entries.Credential.GetById(vaultId, entryId)
	if err != nil {
		resp.Diagnostics.AddError("unable to read entry", err.Error())
		return
	}

	if entryCredentialApiKey.Type != dvls.EntryCredentialType ||
		entryCredentialApiKey.SubType != dvls.EntryCredentialSubTypeApiKey {
		resp.Diagnostics.AddError("invalid entry type", "expected an api key credential entry.")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}

func (r *EntryCredentialApiKeyResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id":          schema.StringAttribute{Computed: true},
					"vault_id":    schema.StringAttribute{Required: true},
					"name":        schema.StringAttribute{Required: true},
					"folder":      schema.StringAttribute{Optional: true},
					"description": schema.StringAttribute{Optional: true},
					"tags":        schema.ListAttribute{ElementType: types.StringType, Optional: true},
					"api_id":      schema.StringAttribute{Optional: true},
					"api_key":     schema.StringAttribute{Optional: true, Sensitive: true},
					"tenant_id":   schema.StringAttribute{Optional: true},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				type v0Model struct {
					Id          types.String   `tfsdk:"id"`
					VaultId     types.String   `tfsdk:"vault_id"`
					Name        types.String   `tfsdk:"name"`
					Folder      types.String   `tfsdk:"folder"`
					Description types.String   `tfsdk:"description"`
					Tags        []types.String `tfsdk:"tags"`
					ApiId       types.String   `tfsdk:"api_id"`
					ApiKey      types.String   `tfsdk:"api_key"`
					TenantId    types.String   `tfsdk:"tenant_id"`
				}
				var prior v0Model
				resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
				if resp.Diagnostics.HasError() {
					return
				}
				upgraded := EntryCredentialApiKeyResourceModel{
					Id:          prior.Id,
					VaultId:     prior.VaultId,
					Name:        prior.Name,
					Folder:      prior.Folder,
					Description: prior.Description,
					Tags:        tagsListToSet(prior.Tags),
					ApiId:       prior.ApiId,
					ApiKey:      prior.ApiKey,
					TenantId:    prior.TenantId,
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
			},
		},
	}
}
