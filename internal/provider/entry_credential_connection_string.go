package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialConnectionStringFromResourceModel(data *EntryCredentialConnectionStringResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialConnectionString := dvls.Entry{
		Id:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		Name:        data.Name.ValueString(),
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypeConnectionString,
		Description: data.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialConnectionStringData{
			ConnectionString: data.ConnectionString.ValueString(),
		},
	}

	return entryCredentialConnectionString
}

func setEntryCredentialConnectionStringResourceModel(entryCredentialConnectionString dvls.Entry, data *EntryCredentialConnectionStringResourceModel) {
	var model EntryCredentialConnectionStringResourceModel

	model.Id = basetypes.NewStringValue(entryCredentialConnectionString.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialConnectionString.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialConnectionString.Name)

	if entryCredentialConnectionString.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialConnectionString.Path)
	}

	if entryCredentialConnectionString.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialConnectionString.Description)
	}

	if entryCredentialConnectionString.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialConnectionString.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialConnectionString.Data != nil {
		data, ok := entryCredentialConnectionString.GetCredentialConnectionStringData()
		if ok {
			if data.ConnectionString != "" {
				model.ConnectionString = basetypes.NewStringValue(data.ConnectionString)
			}
		}
	}

	*data = model
}

func setEntryCredentialConnectionStringDataModel(entryCredentialConnectionString dvls.Entry, data *EntryCredentialConnectionStringDataSourceModel) {
	var model EntryCredentialConnectionStringDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialConnectionString.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialConnectionString.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialConnectionString.Name)

	if entryCredentialConnectionString.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialConnectionString.Path)
	}

	if entryCredentialConnectionString.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialConnectionString.Description)
	}

	if entryCredentialConnectionString.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialConnectionString.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialConnectionString.Data != nil {
		data, ok := entryCredentialConnectionString.GetCredentialConnectionStringData()
		if ok {
			if data.ConnectionString != "" {
				model.ConnectionString = basetypes.NewStringValue(data.ConnectionString)
			}
		}
	}

	*data = model
}
