package provider

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Devolutions/go-dvls"
)

func TestIsAdministrativeRoleAssignmentFilterGone(t *testing.T) {
	requestError := func(status int, body string) error {
		return errors.Join(errors.New("error while fetching role assignments"), &dvls.RequestError{StatusCode: status, Body: []byte(body), Err: errors.New("unexpected status code")})
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "unknown role", err: requestError(http.StatusBadRequest, `{"errors":[{"field":"roleId","message":"roleId does not match an existing role."}]}`), want: true},
		{name: "unknown scope", err: requestError(http.StatusBadRequest, `{"errors":[{"field":"scopeId","message":"x"}]}`), want: true},
		{name: "unknown assignee", err: requestError(http.StatusBadRequest, `{"errors":[{"field":"assigneeId","message":"x"}]}`), want: true},
		{name: "another field", err: requestError(http.StatusBadRequest, `{"errors":[{"field":"scopeType","message":"x"}]}`), want: false},
		{name: "gone and another field", err: requestError(http.StatusBadRequest, `{"errors":[{"field":"roleId","message":"x"},{"field":"scopeType","message":"x"}]}`), want: false},
		{name: "no field", err: requestError(http.StatusBadRequest, `{"message":"bad request"}`), want: false},
		{name: "not json", err: requestError(http.StatusBadRequest, "oops"), want: false},
		{name: "forbidden", err: requestError(http.StatusForbidden, `{"errors":[{"field":"roleId","message":"x"}]}`), want: false},
		{name: "not a request error", err: errors.New("roleId"), want: false},
		{name: "nil", err: nil, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isAdministrativeRoleAssignmentFilterGone(tc.err); got != tc.want {
				t.Fatalf("isAdministrativeRoleAssignmentFilterGone() = %v, want %v", got, tc.want)
			}
		})
	}
}
