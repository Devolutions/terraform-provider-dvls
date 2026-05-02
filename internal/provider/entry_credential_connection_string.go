package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialConnectionStringFromResourceModel(rm *EntryCredentialConnectionStringResourceModel) dvls.Entry {
	tags := tagsSetToSlice(rm.Tags)

	entryCredentialConnectionString := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeConnectionString,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialConnectionStringData{
			ConnectionString: rm.ConnectionString.ValueString(),
		},
	}

	return entryCredentialConnectionString
}

func setEntryCredentialConnectionStringResourceModel(entry dvls.Entry, rm *EntryCredentialConnectionStringResourceModel) {
	var model EntryCredentialConnectionStringResourceModel

	model.Id = basetypes.NewStringValue(entry.Id)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.Name = basetypes.NewStringValue(entry.Name)

	if entry.Path != "" {
		model.Folder = basetypes.NewStringValue(entry.Path)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	model.Tags = tagsSliceToSet(entry.Tags)

	if entry.Data != nil {
		data, ok := entry.GetCredentialConnectionStringData()
		if ok {
			if data.ConnectionString != "" {
				model.ConnectionString = basetypes.NewStringValue(data.ConnectionString)
			}
		}
	}

	*rm = model
}

func setEntryCredentialConnectionStringDataModel(entry dvls.Entry, dsm *EntryCredentialConnectionStringDataSourceModel) {
	var model EntryCredentialConnectionStringDataSourceModel

	model.Id = basetypes.NewStringValue(entry.Id)
	model.VaultId = basetypes.NewStringValue(entry.VaultId)
	model.Name = basetypes.NewStringValue(entry.Name)

	if entry.Path != "" {
		model.Folder = basetypes.NewStringValue(entry.Path)
	}

	if entry.Description != "" {
		model.Description = basetypes.NewStringValue(entry.Description)
	}

	model.Tags = tagsSliceToSet(entry.Tags)

	if entry.Data != nil {
		data, ok := entry.GetCredentialConnectionStringData()
		if ok {
			if data.ConnectionString != "" {
				model.ConnectionString = basetypes.NewStringValue(data.ConnectionString)
			}
		}
	}

	*dsm = model
}
