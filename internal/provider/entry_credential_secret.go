package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialSecretFromResourceModel(rm *EntryCredentialSecretResourceModel) dvls.Entry {
	var tags []string

	for _, v := range rm.Tags {
		tags = append(tags, v.ValueString())
	}

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

	if entry.Tags != nil {
		var tagsBase []types.String

		for _, v := range entry.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

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

func setEntryCredentialSecretDataModel(rm dvls.Entry, dsm *EntryCredentialSecretDataSourceModel) {
	var model EntryCredentialSecretDataSourceModel

	model.Id = basetypes.NewStringValue(rm.Id)
	model.VaultId = basetypes.NewStringValue(rm.VaultId)
	model.Name = basetypes.NewStringValue(rm.Name)

	if rm.Path != "" {
		model.Folder = basetypes.NewStringValue(rm.Path)
	}

	if rm.Description != "" {
		model.Description = basetypes.NewStringValue(rm.Description)
	}

	if rm.Tags != nil {
		var tagsBase []types.String

		for _, v := range rm.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if rm.Data != nil {
		data, ok := rm.GetCredentialAccessCodeData()
		if ok {
			if data.Password != "" {
				model.Secret = basetypes.NewStringValue(data.Password)
			}
		}
	}

	*dsm = model
}
