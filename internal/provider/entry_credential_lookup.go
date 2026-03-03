package provider

import (
	"errors"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var idOrNameConfigValidators = []datasource.ConfigValidator{
	datasourcevalidator.AtLeastOneOf(
		path.MatchRoot("id"),
		path.MatchRoot("name"),
	),
}

func getCredentialEntry(
	client *dvls.Client,
	diagnostics *diag.Diagnostics,
	id, vaultId, name, folder types.String,
	subType string,
	entryTypeName string,
) (dvls.Entry, bool) {
	if !id.IsNull() && !id.IsUnknown() {
		if !name.IsNull() || !folder.IsNull() {
			diagnostics.AddWarning("id takes precedence", "When id is provided, name and folder are ignored.")
		}

		entry, err := client.Entries.Credential.GetById(vaultId.ValueString(), id.ValueString())
		if err != nil {
			diagnostics.AddError(fmt.Sprintf("unable to read %s credential entry", entryTypeName), err.Error())
			return dvls.Entry{}, false
		}

		if entry.Type != dvls.EntryCredentialType || entry.SubType != subType {
			diagnostics.AddError("invalid entry type", fmt.Sprintf("expected a %s credential entry.", entryTypeName))
			return dvls.Entry{}, false
		}

		return entry, true
	}

	var folderPath *string
	if !folder.IsNull() && !folder.IsUnknown() {
		v := folder.ValueString()
		folderPath = &v
	}

	entry, err := client.Entries.Credential.GetByName(
		vaultId.ValueString(),
		name.ValueString(),
		subType,
		dvls.GetByNameOptions{Path: folderPath},
	)
	if err != nil {
		if errors.Is(err, dvls.ErrMultipleEntriesFound) {
			diagnostics.AddError(
				"multiple entries found",
				fmt.Sprintf("more than one entry named %q found, use id to target the correct one", name.ValueString()),
			)
			return dvls.Entry{}, false
		}

		diagnostics.AddError(fmt.Sprintf("unable to read %s credential entry", entryTypeName), err.Error())
		return dvls.Entry{}, false
	}

	return entry, true
}
