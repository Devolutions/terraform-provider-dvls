package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func folderFullPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "\\" + name
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
	// VaultName is user-input only — never returned by the API. Carry it
	// forward from the existing model so plan/state round-trip stays clean.
	vaultName := rm.VaultName

	var model EntryFolderResourceModel
	model.Id = basetypes.NewStringValue(entry.Id)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.VaultName = vaultName
	model.Name = basetypes.NewStringValue(entry.Name)

	if entry.Path != "" {
		model.ParentFolder = basetypes.NewStringValue(entry.Path)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	*rm = model
}
