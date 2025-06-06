package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialSSHKeyFromResourceModel(data *EntryCredentialSSHKeyResourceModel) dvls.Entry {
	var tags []string

	for _, v := range data.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialSSHKey := dvls.Entry{
		Id:          data.Id.ValueString(),
		VaultId:     data.VaultId.ValueString(),
		Name:        data.Name.ValueString(),
		Path:        data.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypePrivateKey,
		Description: data.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialPrivateKeyData{
			OverridePassword: data.Password.ValueString(),
			Passphrase:       data.Passphrase.ValueString(),
			PrivateKey:       data.PrivateKeyData.ValueString(),
			PublicKey:        data.PublicKey.ValueString(),
		},
	}

	return entryCredentialSSHKey
}

func setEntryCredentialSSHKeyResourceModel(entryCredentialSSHKey dvls.Entry, data *EntryCredentialSSHKeyResourceModel) {
	var model EntryCredentialSSHKeyResourceModel

	// VALIDATE IF THE ENTRY IS THE CORRECT SUBTYPE

	model.Id = basetypes.NewStringValue(entryCredentialSSHKey.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialSSHKey.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialSSHKey.Name)

	if entryCredentialSSHKey.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialSSHKey.Path)
	}

	if entryCredentialSSHKey.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialSSHKey.Description)
	}

	if entryCredentialSSHKey.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialSSHKey.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialSSHKey.Data != nil {
		data, ok := entryCredentialSSHKey.GetCredentialPrivateKeyData()
		if ok {
			if data.OverridePassword != "" {
				model.Password = basetypes.NewStringValue(data.OverridePassword)
			}

			if data.Passphrase != "" {
				model.Passphrase = basetypes.NewStringValue(data.Passphrase)
			}

			if data.PrivateKey != "" {
				model.PrivateKeyData = basetypes.NewStringValue(data.PrivateKey)
			}

			if data.PublicKey != "" {
				model.PublicKey = basetypes.NewStringValue(data.PublicKey)
			}
		}
	}

	*data = model
}

func setEntryCredentialSSHKeyDataModel(entryCredentialSSHKey dvls.Entry, data *EntryCredentialSSHKeyDataSourceModel) {
	var model EntryCredentialSSHKeyDataSourceModel

	model.Id = basetypes.NewStringValue(entryCredentialSSHKey.Id)
	model.VaultId = basetypes.NewStringValue(entryCredentialSSHKey.VaultId)
	model.Name = basetypes.NewStringValue(entryCredentialSSHKey.Name)

	if entryCredentialSSHKey.Path != "" {
		model.Folder = basetypes.NewStringValue(entryCredentialSSHKey.Path)
	}

	if entryCredentialSSHKey.Description != "" {
		model.Description = basetypes.NewStringValue(entryCredentialSSHKey.Description)
	}

	if entryCredentialSSHKey.Tags != nil {
		var tagsBase []types.String

		for _, v := range entryCredentialSSHKey.Tags {
			tagsBase = append(tagsBase, basetypes.NewStringValue(v))
		}

		model.Tags = tagsBase
	}

	if entryCredentialSSHKey.Data != nil {
		data, ok := entryCredentialSSHKey.GetCredentialPrivateKeyData()
		if ok {
			if data.OverridePassword != "" {
				model.Password = basetypes.NewStringValue(data.OverridePassword)
			}

			if data.Passphrase != "" {
				model.Passphrase = basetypes.NewStringValue(data.Passphrase)
			}

			if data.PrivateKey != "" {
				model.PrivateKeyData = basetypes.NewStringValue(data.PrivateKey)
			}

			if data.PublicKey != "" {
				model.PublicKey = basetypes.NewStringValue(data.PublicKey)
			}
		}
	}

	*data = model
}
