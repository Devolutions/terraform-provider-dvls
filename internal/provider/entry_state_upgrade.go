package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// emptyTagsToNull normalizes an empty user-provided tags set (`tags = []`) to
// a null value at plan time. DVLS always returns an empty array for entries
// without tags, and `tagsSliceToSet` converts that to a null set; without this
// modifier, a user who explicitly writes `tags = []` would see an
// "inconsistent result after apply" error every plan.
type emptyTagsToNullPlanModifier struct{}

func (m emptyTagsToNullPlanModifier) Description(_ context.Context) string {
	return "Treats `tags = []` as if `tags` were omitted."
}

func (m emptyTagsToNullPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m emptyTagsToNullPlanModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if len(req.PlanValue.Elements()) == 0 {
		resp.PlanValue = types.SetNull(types.StringType)
	}
}

func emptyTagsToNull() planmodifier.Set { return emptyTagsToNullPlanModifier{} }

// DVLS sorts tags server-side, so list ordering carried no meaning.
func tagsListToSet(tags []types.String) types.Set {
	if len(tags) == 0 {
		return basetypes.NewSetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(tags))
	for _, t := range tags {
		elements = append(elements, t)
	}
	return basetypes.NewSetValueMust(types.StringType, elements)
}

// tagsSliceToSet converts the []string DVLS returns on entry reads into the
// types.Set the framework stores. An empty/nil input yields a Null set, not
// an empty set: the API returns [] for entries the user never tagged, and
// matching the user's null plan avoids "inconsistent result after apply".
func tagsSliceToSet(tags []string) types.Set {
	if len(tags) == 0 {
		return basetypes.NewSetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(tags))
	for _, t := range tags {
		elements = append(elements, basetypes.NewStringValue(t))
	}
	return basetypes.NewSetValueMust(types.StringType, elements)
}

// tagsSetToSlice converts the framework's types.Set to the []string DVLS
// expects. Always returns a non-nil slice — DVLS PUT /entry rejects
// "tags":null on every entry type ("PUT requires 'tags' field, use empty
// array if none").
func tagsSetToSlice(tags types.Set) []string {
	if tags.IsNull() || tags.IsUnknown() {
		return []string{}
	}
	elements := tags.Elements()
	out := make([]string, 0, len(elements))
	for _, v := range elements {
		if s, ok := v.(types.String); ok {
			out = append(out, s.ValueString())
		}
	}
	return out
}
