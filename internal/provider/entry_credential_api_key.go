package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialApiKeyFromResourceModel(rm *EntryCredentialApiKeyResourceModel) dvls.Entry {
	var tags []string

	for _, v := range rm.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialApiKey := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeApiKey,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialApiKeyData{
			ApiId:    rm.ApiId.ValueString(),
			ApiKey:   rm.ApiKey.ValueString(),
			TenantId: rm.TenantId.ValueString(),
		},
	}

	return entryCredentialApiKey
}

func setEntryCredentialApiKeyResourceModel(entry dvls.Entry, rm *EntryCredentialApiKeyResourceModel) {
	var model EntryCredentialApiKeyResourceModel

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
		data, ok := entry.GetCredentialApiKeyData()
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

	*rm = model
}

func setEntryCredentialApiKeyDataModel(entry dvls.Entry, dsm *EntryCredentialApiKeyDataSourceModel) {
	var model EntryCredentialApiKeyDataSourceModel

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
		data, ok := entry.GetCredentialApiKeyData()
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

	*dsm = model
}
