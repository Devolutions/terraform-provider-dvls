package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/Devolutions/go-dvls"
	resourcevalidator "github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Entry resources accept either vault_id (canonical) or vault_name (looked up
// at apply time via dvls.Vaults.GetByName). The helpers below keep the schema
// + validators + resolver consistent across every entry resource.

// vaultIDAttribute is the shared `vault_id` schema entry. Optional+Computed
// so the user can set it directly or have it filled in after a vault_name
// lookup; RequiresReplace because changing the vault is a re-creation.
func vaultIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "The ID of the vault. Exactly one of `vault_id` or `vault_name` must be set.",
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			stringplanmodifier.RequiresReplace(),
		},
	}
}

// vaultNameAttribute is the shared `vault_name` schema entry. The provider
// resolves it to a vault id on Create; subsequent reads compare via the
// stored vault_id only.
func vaultNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description:   "The name of the vault. Exactly one of `vault_id` or `vault_name` must be set.",
		Optional:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
}

// entryVaultRefConfigValidators returns the ExactlyOneOf validator for the
// vault_id / vault_name pair. Use from every entry resource's
// ConfigValidators implementation.
func entryVaultRefConfigValidators() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("vault_id"),
			path.MatchRoot("vault_name"),
		),
	}
}

// resolveVaultId returns the canonical vault id given the user's (possibly
// null) vault_id and vault_name attributes. When vault_id is set it is
// returned directly; otherwise the vault is looked up by name via go-dvls.
func resolveVaultId(ctx context.Context, client *dvls.Client, vaultId, vaultName types.String) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !vaultId.IsNull() && !vaultId.IsUnknown() && vaultId.ValueString() != "" {
		return vaultId.ValueString(), diags
	}
	if vaultName.IsNull() || vaultName.IsUnknown() || vaultName.ValueString() == "" {
		diags.AddError("missing vault reference", "one of vault_id or vault_name must be set")
		return "", diags
	}
	vault, err := client.Vaults.GetByNameWithContext(ctx, vaultName.ValueString())
	if err != nil {
		diags.AddError("unable to look up vault by name", err.Error())
		return "", diags
	}
	return vault.Id, diags
}

func isEntryNotFound(err error) bool {
	return errors.Is(err, dvls.ErrEntryNotFound) || dvls.IsNotFound(err)
}

// entryGone reports whether err, returned by an entry route of the vault vaultId, proves that the entry no longer
// exists. DVLS answers the same 404 for a vault the provider's identity cannot open, so a 404 counts only when the
// content of the vault is readable; otherwise the diagnostics hold an error, since the entry may still exist.
func entryGone(ctx context.Context, client *dvls.Client, vaultId string, err error) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !isEntryNotFound(err) {
		return false, diags
	}

	readable, checkErr := client.Vaults.ContentReadableWithContext(ctx, vaultId)
	if checkErr != nil {
		diags.AddError("unable to tell whether the entry was deleted",
			fmt.Sprintf("DVLS answered 404 for the entry, and checking access to the content of vault %s failed: %s", vaultId, checkErr))
		return false, diags
	}
	if !readable {
		diags.AddError("entry not reachable",
			fmt.Sprintf("DVLS answered 404 for the entry, and the provider's identity cannot read the content of vault %s: "+
				"the vault was deleted or the identity lost access to it, so the entry may still exist. "+
				"Restore the access, or, if the vault is gone, remove the resource from the state with terraform state rm.", vaultId))
		return false, diags
	}

	return true, diags
}
