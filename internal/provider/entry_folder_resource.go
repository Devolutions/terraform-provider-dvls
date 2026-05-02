package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EntryFolderResource{}
var _ resource.ResourceWithImportState = &EntryFolderResource{}
var _ resource.ResourceWithModifyPlan = &EntryFolderResource{}

func NewEntryFolderResource() resource.Resource {
	return &EntryFolderResource{}
}

type EntryFolderResource struct {
	client *dvls.Client
}

type EntryFolderResourceModel struct {
	Id           types.String `tfsdk:"id"`
	VaultId      types.String `tfsdk:"vault_id"`
	Name         types.String `tfsdk:"name"`
	ParentFolder types.String `tfsdk:"parent_folder"`
	Description  types.String `tfsdk:"description"`
}

func (r *EntryFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_folder"
}

func (r *EntryFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A DVLS Folder Entry. Folders organize other entries inside a vault. Credential entries that set `folder` require the folder to exist beforehand.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The ID of the folder. This is set by the provider after creation.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_id": schema.StringAttribute{
				Description:   "The ID of the vault.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				// Reject backslash so folderParentFromFullPath stays unambiguous
				// (DVLS GET returns the folder's full path as parent + "\" + name).
				Description: "The name of the folder. Must not contain a backslash.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\\]+$`), "must not contain a backslash"),
				},
			},
			"parent_folder": schema.StringAttribute{
				// DVLS PUT /entry silently ignores `path` changes on folder
				// entries, so the public API has no move operation. Force
				// replacement so users see the destructive plan instead of a
				// no-op update that drifts state from the server.
				MarkdownDescription: "Path of the parent folder. Omit for a folder at the vault root.\n\n" +
					"Changing this forces replacement: DVLS has no folder-move API. " +
					"Any entries still in the folder when it is destroyed are **orphaned** (kept in the vault but unreachable by a `folder = ...` reference). " +
					"Move entries out before changing `parent_folder`.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Description: "The description of the folder.",
				Optional:    true,
			},
		},
	}
}

func (r *EntryFolderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EntryFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EntryFolderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder := newEntryFolderFromResourceModel(plan)

	folderId, err := r.client.Entries.Folder.New(folder)
	if err != nil {
		resp.Diagnostics.AddError("unable to create folder entry", err.Error())
		return
	}

	folder, err = r.client.Entries.Folder.GetById(folder.VaultId, folderId)
	if err != nil {
		resp.Diagnostics.AddError("unable to fetch created folder entry", err.Error())
		return
	}

	setEntryFolderResourceModel(folder, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EntryFolderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder, err := r.client.Entries.Folder.GetById(state.VaultId.ValueString(), state.Id.ValueString())
	if err != nil {
		if dvls.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read folder entry", err.Error())
		return
	}

	setEntryFolderResourceModel(folder, state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EntryFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EntryFolderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folder := newEntryFolderFromResourceModel(plan)

	folder, err := r.client.Entries.Folder.Update(folder)
	if err != nil {
		resp.Diagnostics.AddError("unable to update folder entry", err.Error())
		return
	}

	setEntryFolderResourceModel(folder, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EntryFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EntryFolderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Entries.Folder.DeleteById(state.VaultId.ValueString(), state.Id.ValueString())
	if err != nil {
		if dvls.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to delete folder entry", err.Error())
		return
	}
}

func (r *EntryFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vaultId, entryId, err := parseEntryImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	folder, err := r.client.Entries.Folder.GetById(vaultId, entryId)
	if err != nil {
		resp.Diagnostics.AddError("unable to read folder entry", err.Error())
		return
	}

	if folder.Type != dvls.EntryFolderType {
		resp.Diagnostics.AddError("invalid entry type", "expected a folder entry.")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("vault_id"), vaultId)
	resp.State.SetAttribute(ctx, path.Root("id"), entryId)
}

// ModifyPlan warns at plan time when a folder is being destroyed or replaced
// and still has entries inside. DVLS DELETE on a folder does not cascade to
// its contents — they remain in the vault with a dangling `path` reference,
// which is rarely the user's intent.
func (r *EntryFolderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		return // create — no existing folder to inspect
	}

	var state EntryFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	destroying := req.Plan.Raw.IsNull()
	if !destroying {
		var plan EntryFolderResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
		// In-place updates (rename, description) don't destroy the folder.
		if state.ParentFolder.Equal(plan.ParentFolder) {
			return
		}
	}

	fullPath := folderFullPath(state.ParentFolder.ValueString(), state.Name.ValueString())
	vaultId := state.VaultId.ValueString()

	// go-dvls only exposes type-scoped GetEntries on the same underlying
	// endpoint, so we run both in parallel and surface any errors together.
	// Limitation: Host/Website/Certificate entries are not enumerated (no
	// public method), and sub-folders aren't either (their full path is
	// parent\name\sub, not equal to fullPath). The warning text below makes
	// that explicit.
	var (
		creds, folders []dvls.Entry
		errC, errF     error
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		creds, errC = r.client.Entries.Credential.GetEntriesWithContext(ctx, vaultId, dvls.GetEntriesOptions{Path: &fullPath})
	}()
	go func() {
		defer wg.Done()
		folders, errF = r.client.Entries.Folder.GetEntriesWithContext(ctx, vaultId, dvls.GetEntriesOptions{Path: &fullPath})
	}()
	wg.Wait()
	// Fail loud rather than silently skip the safety check: a transient list
	// error during a destroy/replace would otherwise let the orphaning
	// happen with no warning at all.
	if err := errors.Join(errC, errF); err != nil {
		resp.Diagnostics.AddError("Could not enumerate folder contents", err.Error())
		return
	}

	// Folder.GetEntries(path = "X\Y") returns the folder itself as a self-match
	// because DVLS reports a folder's full path as parent+"\"+name. Filter it out.
	folderId := state.Id.ValueString()
	var children []dvls.Entry
	for _, e := range append(creds, folders...) {
		if e.Id != folderId {
			children = append(children, e)
		}
	}
	if len(children) == 0 {
		return
	}

	names := make([]string, len(children))
	for i, e := range children {
		names[i] = e.Name
	}

	action := "replacement"
	if destroying {
		action = "deletion"
	}
	noun := "entries"
	if len(children) == 1 {
		noun = "entry"
	}
	resp.Diagnostics.AddWarning(
		fmt.Sprintf("Folder %q has %d %s that will be orphaned", fullPath, len(children), noun),
		fmt.Sprintf(
			"Folder %s forces %s and DVLS does not cascade-delete folder contents. "+
				"After apply, the entries below remain in the vault but their `path` will reference a folder that no longer exists, "+
				"making them unreachable by `folder = %q`. Move them out before applying.\n\n"+
				"Direct credential children: %s\n\n"+
				"Note: this check only enumerates direct credential children. Sub-folders, host, website, and certificate entries are NOT listed but will be orphaned the same way.",
			fullPath, action, fullPath, strings.Join(names, ", "),
		),
	)
}
