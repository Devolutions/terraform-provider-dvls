package provider

import (
	"errors"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var errInvalidEntryType = errors.New("the entry does not match the expected credential type")

func fetchCredentialEntry(client *dvls.Client, vaultId, id, name, folder types.String, subType string) (dvls.Entry, error) {
	if !id.IsNull() && !id.IsUnknown() {
		entry, err := client.Entries.Credential.GetById(vaultId.ValueString(), id.ValueString())
		if err != nil {
			return entry, err
		}
		if entry.Type != dvls.EntryCredentialType || entry.SubType != subType {
			return entry, errInvalidEntryType
		}
		return entry, nil
	}

	var folderPath *string
	if !folder.IsNull() && !folder.IsUnknown() {
		v := folder.ValueString()
		folderPath = &v
	}

	return client.Entries.Credential.GetByName(vaultId.ValueString(), name.ValueString(), subType, dvls.GetByNameOptions{Path: folderPath})
}
