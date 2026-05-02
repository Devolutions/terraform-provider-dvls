package provider

import (
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DVLS asymmetry: a folder's `path` field is the parent path on POST, but the
// folder's own full path (parent + "\" + name) on GET. The helpers below hide
// that on either side of the wire.

func folderFullPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "\\" + name
}

func folderParentFromFullPath(fullPath, name string) (parent string, ok bool) {
	suffix := "\\" + name
	if !strings.HasSuffix(fullPath, suffix) {
		return "", false
	}
	return strings.TrimSuffix(fullPath, suffix), true
}

func newEntryFolderFromResourceModel(rm *EntryFolderResourceModel) dvls.Entry {
	// Tags and Data must both be non-nil: PUT /entry rejects "tags":null and
	// expects a `data` object even though the folder type has no data fields.
	return dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.ParentFolder.ValueString(),
		Type:        dvls.EntryFolderType,
		SubType:     dvls.EntryFolderSubTypeFolder,
		Description: rm.Description.ValueString(),
		Tags:        []string{},
		Data:        &dvls.EntryFolderData{},
	}
}

func setEntryFolderResourceModel(entry dvls.Entry, rm *EntryFolderResourceModel) {
	var model EntryFolderResourceModel

	model.Id = basetypes.NewStringValue(entry.Id)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.Name = basetypes.NewStringValue(entry.Name)

	if parent, ok := folderParentFromFullPath(entry.Path, entry.Name); ok && parent != "" {
		model.ParentFolder = basetypes.NewStringValue(parent)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	*rm = model
}
