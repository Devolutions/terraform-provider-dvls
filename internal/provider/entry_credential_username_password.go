package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialUsernamePasswordFromResourceModel(data *EntryCredentialUsernamePasswordResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialUsernamePassword := dvls.Entry{
		Id:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		Name:        data.Name.ValueString(),
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeDefault,
		Description: data.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialDefaultData{
			Username: data.Username.ValueString(),
			Domain:   data.Domain.ValueString(),
			Password: data.Password.ValueString(),
		},
	}

	return entryCredentialUsernamePassword
}

func setEntryCredentialUsernamePasswordResourceModel(entryCredentialUsernamePassword dvls.Entry, data *EntryCredentialUsernamePasswordResourceModel) {
	var model EntryCredentialUsernamePasswordResourceModel

	model.Id = basetypes.NewStringValue(entryCredentialUsernamePassword.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialUsernamePassword.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialUsernamePassword.Name)

	if entryCredentialUsernamePassword.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialUsernamePassword.Path)
	}

	if entryCredentialUsernamePassword.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialUsernamePassword.Description)
	}

	if entryCredentialUsernamePassword.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialUsernamePassword.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialUsernamePassword.Data != nil {
		data, ok := entryCredentialUsernamePassword.GetCredentialDefaultData()
		if ok {
			if data.Username != "" {
				model.Username = basetypes.NewStringValue(data.Username)
			}

			if data.Domain != "" {
				model.Domain = basetypes.NewStringValue(data.Domain)
			}

			if data.Password != "" {
				model.Password = basetypes.NewStringValue(data.Password)
			}
		}
	}

	*data = model
}

func setEntryCredentialUsernamePasswordDataModel(entryCredentialUsernamePassword dvls.Entry, data *EntryCredentialUsernamePasswordDataSourceModel) {
	var model EntryCredentialUsernamePasswordDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialUsernamePassword.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialUsernamePassword.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialUsernamePassword.Name)

	if entryCredentialUsernamePassword.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialUsernamePassword.Path)
	}

	if entryCredentialUsernamePassword.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialUsernamePassword.Description)
	}

	if entryCredentialUsernamePassword.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialUsernamePassword.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialUsernamePassword.Data != nil {
		data, ok := entryCredentialUsernamePassword.GetCredentialDefaultData()
		if ok {
			if data.Username != "" {
				model.Username = basetypes.NewStringValue(data.Username)
			}

			if data.Domain != "" {
				model.Domain = basetypes.NewStringValue(data.Domain)
			}

			if data.Password != "" {
				model.Password = basetypes.NewStringValue(data.Password)
			}
		}
	}

	*data = model
}
