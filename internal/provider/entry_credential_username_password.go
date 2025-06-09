package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialUsernamePasswordFromResourceModel(rm *EntryCredentialUsernamePasswordResourceModel) dvls.Entry {
	var tags []string

	for _, v := range rm.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialUsernamePassword := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeDefault,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialDefaultData{
			Username: rm.Username.ValueString(),
			Domain:   rm.Domain.ValueString(),
			Password: rm.Password.ValueString(),
		},
	}

	return entryCredentialUsernamePassword
}

func setEntryCredentialUsernamePasswordResourceModel(entry dvls.Entry, rm *EntryCredentialUsernamePasswordResourceModel) {
	var model EntryCredentialUsernamePasswordResourceModel

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
		data, ok := entry.GetCredentialDefaultData()
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

	*rm = model
}

func setEntryCredentialUsernamePasswordDataModel(entry dvls.Entry, dsm *EntryCredentialUsernamePasswordDataSourceModel) {
	var model EntryCredentialUsernamePasswordDataSourceModel

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
		data, ok := entry.GetCredentialDefaultData()
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

	*dsm = model
}
