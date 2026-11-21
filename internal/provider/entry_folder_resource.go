package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
var _ resource.ResourceWithConfigValidators = &EntryFolderResource{}

func NewEntryFolderResource() resource.Resource {
	return &EntryFolderResource{}
}

type EntryFolderResource struct {
	client *dvls.Client
}

type EntryFolderResourceModel struct {
	Id           types.String `tfsdk:"id"`
	VaultId      types.String `tfsdk:"vault_id"`
	VaultName    types.String `tfsdk:"vault_name"`
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
			"vault_id":   vaultIDAttribute(),
			"vault_name": vaultNameAttribute(),
			"name": schema.StringAttribute{
				// A backslash is the folder path separator.
				Description: "The name of the folder. Must not contain a backslash.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\\]+$`), "must not contain a backslash"),
				},
			},
			"parent_folder": schema.StringAttribute{
				MarkdownDescription: "Path of the parent folder, e.g. `A\\B`. Omit for a folder at the vault root.\n\n" +
					"Changing this moves the folder, with its content, in place. The new parent folder must already exist: " +
					"declare it as its own `dvls_entry_folder` and reference it, the provider does not create it.",
				Optional: true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the folder.",
				Optional:    true,
			},
		},
	}
}

func (r *EntryFolderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return entryVaultRefConfigValidators()
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

	vaultId, diags := resolveVaultId(ctx, r.client, plan.VaultId, plan.VaultName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.VaultId = types.StringValue(vaultId)

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
		if gone, diags := entryGone(ctx, r.client, state.VaultId.ValueString(), err); gone || diags.HasError() {
			resp.Diagnostics.Append(diags...)
			if gone {
				resp.State.RemoveResource(ctx)
			}
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

	var state EntryFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.ParentFolder.Equal(state.ParentFolder) && plan.ParentFolder.ValueString() != "" {
		exists, err := r.folderExists(ctx, plan.VaultId.ValueString(), plan.ParentFolder.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to look up the new parent folder", err.Error())
			return
		}
		if !exists {
			resp.Diagnostics.AddAttributeError(path.Root("parent_folder"), "parent folder not found",
				fmt.Sprintf("The folder %q does not exist in the vault, so %q cannot be moved into it. "+
					"Declare %q as its own dvls_entry_folder resource and reference it from parent_folder, "+
					"so that it is created before the move; the provider does not create parent folders.",
					plan.ParentFolder.ValueString(), plan.Name.ValueString(), plan.ParentFolder.ValueString()))
			return
		}
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
		if gone, diags := entryGone(ctx, r.client, state.VaultId.ValueString(), err); gone || diags.HasError() {
			resp.Diagnostics.Append(diags...)
			if gone {
				resp.State.RemoveResource(ctx)
			}
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

func (r *EntryFolderResource) folderExists(ctx context.Context, vaultId, fullPath string) (bool, error) {
	parent, name := "", fullPath
	if i := strings.LastIndex(fullPath, "\\"); i >= 0 {
		parent, name = fullPath[:i], fullPath[i+1:]
	}

	folders, err := r.client.Entries.Folder.GetEntriesWithContext(ctx, vaultId, dvls.GetEntriesOptions{Path: &parent})
	if err != nil {
		return false, err
	}

	for _, folder := range folders {
		if strings.EqualFold(folder.Name, name) && strings.EqualFold(folder.Path, parent) {
			return true, nil
		}
	}

	return false, nil
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
		// Skip when the change is an in-place update (rename, move, description).
		if state.VaultId.Equal(plan.VaultId) &&
			state.VaultName.Equal(plan.VaultName) {
			return
		}
	}

	fullPath := folderFullPath(state.ParentFolder.ValueString(), state.Name.ValueString())
	vaultId := state.VaultId.ValueString()

	children, err := r.client.Entries.Credential.GetEntriesWithContext(ctx, vaultId, dvls.GetEntriesOptions{Path: &fullPath})
	if err != nil {
		resp.Diagnostics.AddError("Could not enumerate folder contents", err.Error())
		return
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
