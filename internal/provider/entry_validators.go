package provider

import (
	"fmt"
	"strings"
)

func parseEntryImportId(id string) (vaultId, entryId string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected <vault_id>/<entry_id>", id)
	}

	return parts[0], parts[1], nil
}
