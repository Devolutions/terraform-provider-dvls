package provider

import (
	"errors"
	"net/http"
	"strings"
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

func TestSingleAssignment(t *testing.T) {
	const assignee = "66666666-7777-8888-9999-000000000000"
	assignment := func(id, assigneeId string) dvls.AdministrativeRoleAssignment {
		return dvls.AdministrativeRoleAssignment{Id: id, AssigneeId: assigneeId}
	}

	tests := []struct {
		name        string
		assignments []dvls.AdministrativeRoleAssignment
		wantId      string
		wantFound   bool
		wantError   string
	}{
		{name: "none", assignments: nil},
		{name: "other assignee only", assignments: []dvls.AdministrativeRoleAssignment{assignment("a1", "someone-else")}},
		{name: "one", assignments: []dvls.AdministrativeRoleAssignment{assignment("a1", "someone-else"), assignment("a2", assignee)}, wantId: "a2", wantFound: true},
		{name: "temporary next to permanent", assignments: []dvls.AdministrativeRoleAssignment{assignment("temporary", assignee), assignment("permanent", assignee)}, wantError: "temporary, permanent"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, found, err := singleAssignment(tc.assignments, assignee)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("error = %v, want one naming %q", err, tc.wantError)
				}
				if found {
					t.Fatal("found = true with an ambiguous match")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tc.wantFound || got.Id != tc.wantId {
				t.Fatalf("got (%q, %v), want (%q, %v)", got.Id, found, tc.wantId, tc.wantFound)
			}
		})
	}
}
