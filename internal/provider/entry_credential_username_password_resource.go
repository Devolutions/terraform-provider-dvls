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
var _ resource.Resource = &EntryCredentialUsernamePasswordResource{}
var _ resource.ResourceWithImportState = &EntryCredentialUsernamePasswordResource{}

func NewEntryCredentialUsernamePasswordResource() resource.Resource {
	return &EntryCredentialUsernamePasswordResource{}
}

// EntryCredentialUsernamePasswordResource defines the resource implementation.
type EntryCredentialUsernamePasswordResource struct {
	client *dvls.Client
}

// EntryCredentialUsernamePasswordResourceModel describes the resource data model.
type EntryCredentialUsernamePasswordResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	VaultId     types.String   `tfsdk:"vault_id"`
	Name        types.String   `tfsdk:"name"`
	Folder      types.String   `tfsdk:"folder"`
	Description types.String   `tfsdk:"description"`
	Tags        []types.String `tfsdk:"tags"`

	// General
	Username types.String `tfsdk:"username"`
	Domain   types.String `tfsdk:"domain"`
	Password types.String `tfsdk:"password"`
}

func (r *EntryCredentialUsernamePasswordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_credential_username_password"
}

func (r *EntryCredentialUsernamePasswordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Username and Password Credential Entry",

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
			"username": schema.StringAttribute{
				Description: "The entry credential username.",
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The entry credential domain.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The entry credential password.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *EntryCredentialUsernamePasswordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryCredentialUsernamePasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryCredentialUsernamePasswordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialUsernamePassword := newEntryCredentialUsernamePasswordFromResourceModel(plan)

	entryCredentialUsernamePasswordId, err := r.client.Entries.Credential.New(entryCredentialUsernamePassword)
	if err != nil {
		resp.Diagnostics.AddError("unable to create username password credential entry", err.Error())
		return
	}

	entryCredentialUsernamePassword, err = r.client.Entries.Credential.GetById(entryCredentialUsernamePassword.VaultId, entryCredentialUsernamePasswordId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created username password credential entry", err.Error())
		return
	}

	setEntryCredentialUsernamePasswordResourceModel(entryCredentialUsernamePassword, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialUsernamePasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryCredentialUsernamePasswordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialUsernamePassword := newEntryCredentialUsernamePasswordFromResourceModel(state)

	entryCredentialUsernamePassword, err := r.client.Entries.Credential.GetById(entryCredentialUsernamePassword.VaultId, entryCredentialUsernamePassword.Id)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read username password credential entry", err.Error())
		return
	}

	setEntryCredentialUsernamePasswordResourceModel(entryCredentialUsernamePassword, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryCredentialUsernamePasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryCredentialUsernamePasswordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialUsernamePassword := newEntryCredentialUsernamePasswordFromResourceModel(plan)

	_, err := r.client.Entries.Credential.Update(entryCredentialUsernamePassword)
	if err != nil {
		resp.Diagnostics.AddError("unable to update username password credential entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryCredentialUsernamePasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryCredentialUsernamePasswordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entryCredentialUsernamePassword := newEntryCredentialUsernamePasswordFromResourceModel(state)

	err := r.client.Entries.Credential.Delete(entryCredentialUsernamePassword)
	if err != nil {
		if strings.Contains(err.Error(), dvls.SaveResultNotFound.String()) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete username password credential entry", err.Error())
		return
	}
}

func (r *EntryCredentialUsernamePasswordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	entryCredentialUsernamePassword, err := r.client.Entries.Credential.GetById(vaultId, entryId)
	if err != nil {
		resp.Diagnostics.AddError("unable to read entry", err.Error())
		return
	}

	if entryCredentialUsernamePassword.Type != dvls.EntryCredentialType ||
		entryCredentialUsernamePassword.SubType != dvls.EntryCredentialSubTypeDefault {
		resp.Diagnostics.AddError("invalid entry type", "expected an username password credential entry.")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}
