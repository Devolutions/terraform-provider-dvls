package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialSecretFromResourceModel(rm *EntryCredentialSecretResourceModel) dvls.Entry {
	tags := tagsSetToSlice(rm.Tags)

	entryCredentialSecret := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeAccessCode,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialAccessCodeData{
			Password: rm.Secret.ValueString(),
		},
	}

	return entryCredentialSecret
}

func setEntryCredentialSecretResourceModel(entry dvls.Entry, rm *EntryCredentialSecretResourceModel) {
	var model EntryCredentialSecretResourceModel

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
		data, ok := entry.GetCredentialAccessCodeData()
		if ok {
			if data.Password != "" {
				model.Secret = basetypes.NewStringValue(data.Password)
			}
		}
	}

	*rm = model
}

func setEntryCredentialSecretDataModel(entry dvls.Entry, dsm *EntryCredentialSecretDataSourceModel) {
	var model EntryCredentialSecretDataSourceModel

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
		data, ok := entry.GetCredentialAccessCodeData()
		if ok {
			if data.Password != "" {
				model.Secret = basetypes.NewStringValue(data.Password)
			}
		}
	}

	*dsm = model
}
