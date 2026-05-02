package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialAzureServicePrincipalFromResourceModel(rm *EntryCredentialAzureServicePrincipalResourceModel) dvls.Entry {
	tags := tagsSetToSlice(rm.Tags)

	entryCredentialAzureServicePrincipal := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeAzureServicePrincipal,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialAzureServicePrincipalData{
			ClientId:     rm.ClientId.ValueString(),
			ClientSecret: rm.ClientSecret.ValueString(),
			TenantId:     rm.TenantId.ValueString(),
		},
	}

	return entryCredentialAzureServicePrincipal
}

func setEntryCredentialAzureServicePrincipalResourceModel(entry dvls.Entry, rm *EntryCredentialAzureServicePrincipalResourceModel) {
	var model EntryCredentialAzureServicePrincipalResourceModel

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
		data, ok := entry.GetCredentialAzureServicePrincipalData()
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

	*rm = model
}

func setEntryCredentialAzureServicePrincipalDataModel(entry dvls.Entry, dsm *EntryCredentialAzureServicePrincipalDataSourceModel) {
	var model EntryCredentialAzureServicePrincipalDataSourceModel

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
		data, ok := entry.GetCredentialAzureServicePrincipalData()
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

	*dsm = model
}
