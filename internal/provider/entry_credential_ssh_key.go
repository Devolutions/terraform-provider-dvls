package provider

import (
	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newEntryCredentialSSHKeyFromResourceModel(rm *EntryCredentialSSHKeyResourceModel) dvls.Entry {
	var tags []string

	for _, v := range rm.Tags {
		tags = append(tags, v.ValueString())
	}

	entryCredentialSSHKey := dvls.Entry{
		Id:          rm.Id.ValueString(),
		VaultId:     rm.VaultId.ValueString(),
		Name:        rm.Name.ValueString(),
		Path:        rm.Folder.ValueString(),
		Type:        dvls.EntryCredentialType,
		SubType:     dvls.EntryCredentialSubTypePrivateKey,
		Description: rm.Description.ValueString(),
		Tags:        tags,
		Data: dvls.EntryCredentialPrivateKeyData{
			Username:   rm.Username.ValueString(),
			Password:   rm.Password.ValueString(),
			Passphrase: rm.Passphrase.ValueString(),
			PrivateKey: rm.PrivateKeyData.ValueString(),
			PublicKey:  rm.PublicKey.ValueString(),
		},
	}

	return entryCredentialSSHKey
}

func setEntryCredentialSSHKeyResourceModel(entry dvls.Entry, rm *EntryCredentialSSHKeyResourceModel) {
	var model EntryCredentialSSHKeyResourceModel

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
		data, ok := entry.GetCredentialPrivateKeyData()
		if ok {
			if data.Username != "" {
				model.Username = basetypes.NewStringValue(data.Username)
			}

			if data.Password != "" {
				model.Password = basetypes.NewStringValue(data.Password)
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

	*rm = model
}

func setEntryCredentialSSHKeyDataModel(entry dvls.Entry, dsm *EntryCredentialSSHKeyDataSourceModel) {
	var model EntryCredentialSSHKeyDataSourceModel

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
		data, ok := entry.GetCredentialPrivateKeyData()
		if ok {
			if data.Username != "" {
				model.Username = basetypes.NewStringValue(data.Username)
			}

			if data.Password != "" {
				model.Password = basetypes.NewStringValue(data.Password)
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

	*dsm = model
}
