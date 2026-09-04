package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var userAuthenticationTypes = map[dvls.UserAuthenticationType]string{
	dvls.UserAuthenticationBuiltin:      "builtin",
	dvls.UserAuthenticationLocalWindows: "local_windows",
	dvls.UserAuthenticationSqlServer:    "sql_server",
	dvls.UserAuthenticationDomain:       "domain",
	dvls.UserAuthenticationOffice365:    "office_365",
	dvls.UserAuthenticationNone:         "none",
	dvls.UserAuthenticationCloud:        "cloud",
	dvls.UserAuthenticationLegacy:       "legacy",
	dvls.UserAuthenticationAzureAD:      "azure_ad",
	dvls.UserAuthenticationApplication:  "application",
	dvls.UserAuthenticationOkta:         "okta",
	dvls.UserAuthenticationPingOne:      "ping_one",
	dvls.UserAuthenticationContractor:   "contractor",
}

var userGroupTypes = map[dvls.UserGroupType]string{
	dvls.UserGroupTypeActiveDirectory: "active_directory",
	dvls.UserGroupTypeCustom:          "custom",
	dvls.UserGroupTypeOffice365:       "office_365",
	dvls.UserGroupTypeOkta:            "okta",
	dvls.UserGroupTypePingOne:         "ping_one",
}

func exactlyOneOfSentence(attributes []string) string {
	quoted := make([]string, 0, len(attributes))
	for _, attribute := range attributes {
		quoted = append(quoted, "`"+attribute+"`")
	}

	last := len(quoted) - 1
	return fmt.Sprintf("Exactly one of %s or %s must be set.", strings.Join(quoted[:last], ", "), quoted[last])
}

func lookupAttribute(description string, attributes []string) schema.StringAttribute {
	return schema.StringAttribute{
		Description: description + " " + exactlyOneOfSentence(attributes),
		Optional:    true,
		Computed:    true,
	}
}

func principalIdAttribute(kind string, attributes []string) schema.StringAttribute {
	attribute := lookupAttribute(fmt.Sprintf("The ID of the %s.", kind), attributes)
	attribute.Validators = []validator.String{uuidValidator{kind + " id"}}
	return attribute
}

var errPrincipalNotFound = errors.New("no principal matches the given lookup")
var errMultiplePrincipals = errors.New("more than one principal matches the given lookup")

type lookupKey[T any] struct {
	value types.String
	get   func(T) string
}

func findPrincipal[T any](ctx context.Context, list func(context.Context) ([]T, error), keys []lookupKey[T]) (T, error) {
	var zero T

	items, err := list(ctx)
	if err != nil {
		return zero, err
	}

	for _, key := range keys {
		if !key.value.IsNull() && !key.value.IsUnknown() {
			return matchOne(items, key.value.ValueString(), key.get)
		}
	}

	return zero, errPrincipalNotFound
}

func matchOne[T any](items []T, value string, get func(T) string) (T, error) {
	var match T
	count := 0
	for _, item := range items {
		if get(item) == value {
			match = item
			count++
		}
	}

	switch count {
	case 0:
		return match, errPrincipalNotFound
	case 1:
		return match, nil
	default:
		var zero T
		return zero, errMultiplePrincipals
	}
}

func principalErrorDetail(err error, kind string) string {
	switch {
	case errors.Is(err, errPrincipalNotFound):
		return fmt.Sprintf("no %s matches the given lookup", kind)
	case errors.Is(err, errMultiplePrincipals):
		return fmt.Sprintf("more than one %s matches, use id to target the correct one", kind)
	default:
		return err.Error()
	}
}

var userLookupKeys = []string{"id", "name", "full_name"}

func userLookup(id, name, fullName types.String) []lookupKey[dvls.User] {
	return []lookupKey[dvls.User]{
		{id, func(user dvls.User) string { return user.Id }},
		{name, func(user dvls.User) string { return user.Name }},
		{fullName, func(user dvls.User) string { return user.FullName }},
	}
}
