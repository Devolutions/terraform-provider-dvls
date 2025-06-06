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
var _ resource.Resource = &EntryCredentialConnectionStringResource{}
var _ resource.ResourceWithImportState = &EntryCredentialConnectionStringResource{}

func NewEntryCredentialConnectionStringResource() resource.Resource {
	return &EntryCredentialConnectionStringResource{}
}

// EntryCredentialConnectionStringResource defines the resource implementation.
type EntryCredentialConnectionStringResource struct {
	client *dvls.Client
}

// EntryCredentialConnectionStringResourceModel describes the resource data model.
type EntryCredentialConnectionStringResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Folder      types.String   `tfsdk:"folder"`
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`

	// General
	ConnectionString types.String `tfsdk:"connection_string"`
}

func (r *EntryCredentialConnectionStringResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_connection_string"
}

func (r *EntryCredentialConnectionStringResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Connection String Credential Entry",

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
			"connection_string": schema.StringAttribute{
				Description: "The entry credential connection string.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *EntryCredentialConnectionStringResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryCredentialConnectionStringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialConnectionStringResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialConnectionString := newEntryCredentialConnectionStringFromResourceModel(plan)

	entryCredentialConnectionStringId, err := r.client.Entries.Credential.New(entryCredentialConnectionString)
	if err != nil {
		resp.Diagnostics.AddError("unable to create connection string credential entry", err.Error())
		return
	}

	entryCredentialConnectionString, err = r.client.Entries.Credential.GetById(entryCredentialConnectionString.VaultId, entryCredentialConnectionStringId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created connection string credential entry", err.Error())
		return
	}

	setEntryCredentialConnectionStringResourceModel(entryCredentialConnectionString, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialConnectionStringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialConnectionStringResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialConnectionString := newEntryCredentialConnectionStringFromResourceModel(state)

	entryCredentialConnectionString, err := r.client.Entries.Credential.GetById(entryCredentialConnectionString.VaultId, entryCredentialConnectionString.Id)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read connection string credential entry", err.Error())
		return
	}

	setEntryCredentialConnectionStringResourceModel(entryCredentialConnectionString, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialConnectionStringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialConnectionStringResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialConnectionString := newEntryCredentialConnectionStringFromResourceModel(plan)

	_, err := r.client.Entries.Credential.Update(entryCredentialConnectionString)
	if err != nil {
		resp.Diagnostics.AddError("unable to update connection string credential entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialConnectionStringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialConnectionStringResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialConnectionString := newEntryCredentialConnectionStringFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialConnectionString)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete connection string credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialConnectionStringResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}
