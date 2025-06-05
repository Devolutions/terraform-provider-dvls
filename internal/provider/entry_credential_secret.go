package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialSecretFromResourceModel(data *EntryCredentialSecretResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialSecret := dvls.Entry{
		Id:      data.Id.ValueString(),
		VaultId: data.VaultId.ValueString(),
		Name:    data.Name.ValueString(),
		Path:    data.Folder.ValueString(),
		Type:    dvls.EntryCredentialType,
		SubType: dvls.EntryCredentialSubTypeAccessCode,

		Data: dvls.EntryCredentialAccessCodeData{
			Password: data.Secret.ValueString(),
		},

		Description: data.Description.ValueString(),
		Tags:        tags,
	}

	return entryCredentialSecret
}

func setEntryCredentialSecretResourceModel(entryCredentialSecret dvls.Entry, data *EntryCredentialSecretResourceModel) {
	var model EntryCredentialSecretResourceModel

	model.Id = basetypes.NewStringValue(entryCredentialSecret.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialSecret.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialSecret.Name)

	if entryCredentialSecret.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialSecret.Path)
	}

	if entryCredentialSecret.Data != nil {
		data, ok := entryCredentialSecret.GetCredentialAccessCodeData()
		if ok {
			if data.Password != "" {
				model.Secret = basetypes.NewStringValue(data.Password)
			}
		}
	}

	if entryCredentialSecret.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialSecret.Description)
	}

	if entryCredentialSecret.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialSecret.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}

func setEntryCredentialSecretDataModel(entryCredentialSecret dvls.Entry, data *EntryCredentialSecretDataSourceModel) {
	var model EntryCredentialSecretDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialSecret.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialSecret.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialSecret.Name)

	if entryCredentialSecret.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialSecret.Path)
	}

	if entryCredentialSecret.Data != nil {
		data, ok := entryCredentialSecret.GetCredentialAccessCodeData()
		if ok {
			if data.Password != "" {
				model.Secret = basetypes.NewStringValue(data.Password)
			}
		}
	}

	if entryCredentialSecret.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialSecret.Description)
	}

	if entryCredentialSecret.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialSecret.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}
