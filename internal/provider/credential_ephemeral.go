package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ephemeralResourceBase provides the shared dvls client and Configure method for ephemeral resources.
type ephemeralResourceBase struct {
	client *dvls.Client
}

func (b *ephemeralResourceBase) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dvls.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Ephemeral Resource Configure Type",
			fmt.Sprintf("Expected *dvls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	b.client = client
}

// credentialEphemeralBase adds the id/name mutual-exclusion validator. Cert ephemeral skips it since id is required.
type credentialEphemeralBase struct {
	ephemeralResourceBase
}

func (b *credentialEphemeralBase) ConfigValidators(_ context.Context) []ephemeral.ConfigValidator {
	return []ephemeral.ConfigValidator{
		ephemeralvalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

// credentialEphemeralCommonAttributes returns the attributes shared by every credential ephemeral.
func credentialEphemeralCommonAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the entry.",
			Optional:    true,
			Computed:    true,
			Validators:  []validator.String{entryIdValidator{}},
		},
		"vault_id": schema.StringAttribute{
			Description: "The ID of the vault.",
			Required:    true,
			Validators:  []validator.String{vaultIdValidator{}},
		},
		"name": schema.StringAttribute{
			Description: "The name of the entry.",
			Optional:    true,
			Computed:    true,
		},
		"folder": schema.StringAttribute{
			Description: "The folder path to search in. Returns entries in the specified folder and all sub-folders.",
			Optional:    true,
			Computed:    true,
			Validators:  []validator.String{stringvalidator.AlsoRequires(path.MatchRoot("name"))},
		},
		"description": schema.StringAttribute{
			Description: "The description of the entry.",
			Computed:    true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of tags added to the entry.",
			Computed:    true,
		},
	}
}
