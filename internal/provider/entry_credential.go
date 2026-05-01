package provider

import (
	"errors"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// credentialSubTypeLabels maps subtype constants to human-readable error labels.
var credentialSubTypeLabels = map[string]string{
	dvls.EntryCredentialSubTypeAccessCode:            "secret",
	dvls.EntryCredentialSubTypeApiKey:                "api key",
	dvls.EntryCredentialSubTypeAzureServicePrincipal: "azure service principal",
	dvls.EntryCredentialSubTypeConnectionString:      "connection string",
	dvls.EntryCredentialSubTypeDefault:               "username password",
	dvls.EntryCredentialSubTypePrivateKey:            "ssh key",
}

func appendCredentialFetchError(diags *diag.Diagnostics, err error, name types.String, subType string) {
	if errors.Is(err, dvls.ErrMultipleEntriesFound) {
		diags.AddError(
			"multiple entries found",
			fmt.Sprintf("more than one entry named %q found, use id to target the correct one", name.ValueString()),
		)
		return
	}
	label, ok := credentialSubTypeLabels[subType]
	if !ok {
		label = "credential"
	}
	diags.AddError("unable to read "+label+" credential entry", err.Error())
}
