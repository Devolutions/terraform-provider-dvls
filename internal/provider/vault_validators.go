package provider

import (
	"context"
	"fmt"

	"github.com/Devolutions/go-dvls"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type vaultIdValidator struct{}

func (v vaultIdValidator) Description(_ context.Context) string {
	return "vault must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)"
}

func (v vaultIdValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v vaultIdValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	id := request.ConfigValue.ValueString()

	_, err := uuid.Parse(id)
	if err != nil {
		response.Diagnostics.AddError("vault id is not a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", err.Error())
		return
	}
}

type vaultEnumValidator[K dvls.VaultSecurityLevel | dvls.VaultVisibility | dvls.VaultContentType] struct {
	lookup    map[K]string
	errTitle  string
	errDetail string
}

func (v vaultEnumValidator[K]) Description(_ context.Context) string {
	return v.errDetail
}

func (v vaultEnumValidator[K]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v vaultEnumValidator[K]) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	_, err := lookupMapValue(v.lookup, request.ConfigValue.ValueString())
	if err != nil {
		response.Diagnostics.AddError(v.errTitle, v.errDetail)
		return
	}
}

func newVaultSecurityLevelValidator() vaultEnumValidator[dvls.VaultSecurityLevel] {
	return vaultEnumValidator[dvls.VaultSecurityLevel]{
		lookup:    vaultSecurityLevels,
		errTitle:  "vault security level is invalid",
		errDetail: fmt.Sprintf("valid values are: %s", vaultSecurityLevelValues),
	}
}

func newVaultVisibilityValidator() vaultEnumValidator[dvls.VaultVisibility] {
	return vaultEnumValidator[dvls.VaultVisibility]{
		lookup:    vaultVisibilities,
		errTitle:  "vault visibility is invalid",
		errDetail: fmt.Sprintf("valid values are: %s", vaultVisibilityValues),
	}
}

func newVaultContentTypeValidator() vaultEnumValidator[dvls.VaultContentType] {
	return vaultEnumValidator[dvls.VaultContentType]{
		lookup:    vaultContentTypes,
		errTitle:  "vault content type is invalid",
		errDetail: fmt.Sprintf("valid values are: %s", vaultContentTypeValues),
	}
}
