package provider

import (
	"context"

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
