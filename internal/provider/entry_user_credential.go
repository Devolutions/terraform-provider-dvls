package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryUserCredentialFromResourceModel(data *EntryUserCredentialResourceModel) dvls.EntryUserCredential {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryusercredential := dvls.EntryUserCredential{
		ID:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		EntryName:   data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Credentials: dvls.EntryCredentials{Username: data.Username.ValueString(), Password: data.Password.ValueString()},
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryTypeCredential,
		SubType:     dvls.EntrySubTypeDefault,
		Tags:        tags,
	}
	return entryusercredential
}

func setEntryUserCredentialResourceModel(entryusercredential dvls.EntryUserCredential, data *EntryUserCredentialResourceModel) {
	var model EntryUserCredentialResourceModel

	model.Id = basetypes.NewStringValue(entryusercredential.ID)
	model.VaultId = basetypes.NewStringValue(entryusercredential.VaultId)
	model.Name = basetypes.NewStringValue(entryusercredential.EntryName)

	if entryusercredential.Credentials.Password != "" {
		model.Password = basetypes.NewStringValue(entryusercredential.Credentials.Password)
	}

	if entryusercredential.Description != "" {
		model.Description = basetypes.NewStringValue(entryusercredential.Description)
	}

	if entryusercredential.Credentials.Username != "" {
		model.Username = basetypes.NewStringValue(entryusercredential.Credentials.Username)
	}

	if entryusercredential.Path != "" {
		model.Folder = basetypes.NewStringValue(entryusercredential.Path)
	}

	if entryusercredential.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryusercredential.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}

func setEntryUserCredentialDataModel(entryusercredential dvls.EntryUserCredential, data *EntryUserCredentialDataSourceModel) {
	var model EntryUserCredentialDataSourceModel

	model.Id = basetypes.NewStringValue(entryusercredential.ID)
	model.VaultId = basetypes.NewStringValue(entryusercredential.VaultId)
	model.Name = basetypes.NewStringValue(entryusercredential.EntryName)

	if entryusercredential.Credentials.Password != "" {
		model.Password = basetypes.NewStringValue(entryusercredential.Credentials.Password)
	}

	if entryusercredential.Description != "" {
		model.Description = basetypes.NewStringValue(entryusercredential.Description)
	}

	if entryusercredential.Credentials.Username != "" {
		model.Username = basetypes.NewStringValue(entryusercredential.Credentials.Username)
	}

	if entryusercredential.Path != "" {
		model.Folder = basetypes.NewStringValue(entryusercredential.Path)
	}

	if entryusercredential.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryusercredential.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	*data = model
}
