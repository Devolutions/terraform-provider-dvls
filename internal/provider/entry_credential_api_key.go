package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialApiKeyFromResourceModel(data *EntryCredentialApiKeyResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialApiKey := dvls.Entry{
		Id:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		Name:        data.Name.ValueString(),
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeApiKey,
		Description: data.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialApiKeyData{
			ApiId:    data.ApiId.ValueString(),
			ApiKey:   data.ApiKey.ValueString(),
			TenantId: data.TenantId.ValueString(),
		},
	}

	return entryCredentialApiKey
}

func setEntryCredentialApiKeyResourceModel(entryCredentialApiKey dvls.Entry, data *EntryCredentialApiKeyResourceModel) {
	var model EntryCredentialApiKeyResourceModel

	model.Id = basetypes.NewStringValue(entryCredentialApiKey.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialApiKey.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialApiKey.Name)

	if entryCredentialApiKey.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialApiKey.Path)
	}

	if entryCredentialApiKey.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialApiKey.Description)
	}

	if entryCredentialApiKey.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialApiKey.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialApiKey.Data != nil {
		data, ok := entryCredentialApiKey.GetCredentialApiKeyData()
		if ok {
			if data.ApiId != "" {
				model.ApiId = basetypes.NewStringValue(data.ApiId)
			}

			if data.ApiKey != "" {
				model.ApiKey = basetypes.NewStringValue(data.ApiKey)
			}

			if data.TenantId != "" {
				model.TenantId = basetypes.NewStringValue(data.TenantId)
			}
		}
	}

	*data = model
}

func setEntryCredentialApiKeyDataModel(entryCredentialApiKey dvls.Entry, data *EntryCredentialApiKeyDataSourceModel) {
	var model EntryCredentialApiKeyDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialApiKey.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialApiKey.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialApiKey.Name)

	if entryCredentialApiKey.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialApiKey.Path)
	}

	if entryCredentialApiKey.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialApiKey.Description)
	}

	if entryCredentialApiKey.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialApiKey.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialApiKey.Data != nil {
		data, ok := entryCredentialApiKey.GetCredentialApiKeyData()
		if ok {
			if data.ApiId != "" {
				model.ApiId = basetypes.NewStringValue(data.ApiId)
			}

			if data.ApiKey != "" {
				model.ApiKey = basetypes.NewStringValue(data.ApiKey)
			}

			if data.TenantId != "" {
				model.TenantId = basetypes.NewStringValue(data.TenantId)
			}
		}
	}

	*data = model
}
