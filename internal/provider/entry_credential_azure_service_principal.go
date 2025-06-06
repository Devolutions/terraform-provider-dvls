package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialAzureServicePrincipalFromResourceModel(data *EntryCredentialAzureServicePrincipalResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialAzureServicePrincipal := dvls.Entry{
		Id:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		Name:        data.Name.ValueString(),
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeAzureServicePrincipal,
		Description: data.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialAzureServicePrincipalData{
			ClientId:     data.ClientId.ValueString(),
			ClientSecret: data.ClientSecret.ValueString(),
			TenantId:     data.TenantId.ValueString(),
		},
	}

	return entryCredentialAzureServicePrincipal
}

func setEntryCredentialAzureServicePrincipalResourceModel(entryCredentialAzureServicePrincipal dvls.Entry, data *EntryCredentialAzureServicePrincipalResourceModel) {
	var model EntryCredentialAzureServicePrincipalResourceModel

	model.Id = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Name)

	if entryCredentialAzureServicePrincipal.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Path)
	}

	if entryCredentialAzureServicePrincipal.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Description)
	}

	if entryCredentialAzureServicePrincipal.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialAzureServicePrincipal.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialAzureServicePrincipal.Data != nil {
		data, ok := entryCredentialAzureServicePrincipal.GetCredentialAzureServicePrincipalData()
		if ok {
			if data.ClientId != "" {
				model.ClientId = basetypes.NewStringValue(data.ClientId)
			}

			if data.ClientSecret != "" {
				model.ClientSecret = basetypes.NewStringValue(data.ClientSecret)
			}

			if data.TenantId != "" {
				model.TenantId = basetypes.NewStringValue(data.TenantId)
			}
		}
	}

	*data = model
}

func setEntryCredentialAzureServicePrincipalDataModel(entryCredentialAzureServicePrincipal dvls.Entry, data *EntryCredentialAzureServicePrincipalDataSourceModel) {
	var model EntryCredentialAzureServicePrincipalDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Name)

	if entryCredentialAzureServicePrincipal.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Path)
	}

	if entryCredentialAzureServicePrincipal.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialAzureServicePrincipal.Description)
	}

	if entryCredentialAzureServicePrincipal.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialAzureServicePrincipal.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialAzureServicePrincipal.Data != nil {
		data, ok := entryCredentialAzureServicePrincipal.GetCredentialAzureServicePrincipalData()
		if ok {
			if data.ClientId != "" {
				model.ClientId = basetypes.NewStringValue(data.ClientId)
			}

			if data.ClientSecret != "" {
				model.ClientSecret = basetypes.NewStringValue(data.ClientSecret)
			}

			if data.TenantId != "" {
				model.TenantId = basetypes.NewStringValue(data.TenantId)
			}
		}
	}

	*data = model
}
