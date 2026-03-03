package provider

import (
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func fetchCredentialEntry(client *dvls.Client, vaultId, id, name, folder types.String, subType string) (dvls.Entry, error) {
	if !id.IsNull() && !id.IsUnknown() {
		entry, err := client.Entries.Credential.GetById(vaultId.ValueString(), id.ValueString())
		if err != nil {
			return entry, err
		}
		if entry.Type != dvls.EntryCredentialType || entry.SubType != subType {
			return entry, fmt.Errorf("expected entry type %q with subtype %q, got type %q with subtype %q",
				dvls.EntryCredentialType, subType, entry.Type, entry.SubType)
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
