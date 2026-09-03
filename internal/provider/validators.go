package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type uuidValidator struct {
	name string
}

func (v uuidValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s must be a valid UUID (ex.: 00000000-0000-0000-0000-000000000000)", v.name)
}

func (v uuidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uuidValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	_, err := uuid.Parse(request.ConfigValue.ValueString())
	if err != nil {
		response.Diagnostics.AddAttributeError(request.Path, v.Description(ctx), err.Error())
	}
}

func exactlyOneOfConfigValidators(attributes ...string) []datasource.ConfigValidator {
	expressions := make([]path.Expression, 0, len(attributes))
	for _, attribute := range attributes {
		expressions = append(expressions, path.MatchRoot(attribute))
	}

	return []datasource.ConfigValidator{datasourcevalidator.ExactlyOneOf(expressions...)}
}
