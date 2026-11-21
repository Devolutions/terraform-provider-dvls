package provider

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Devolutions/go-dvls"
)

func TestEntryGone(t *testing.T) {
	const vaultId = "11111111-2222-3333-4444-555555555555"

	notFound := &dvls.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("unexpected status code 404")}

	tests := []struct {
		name       string
		err        error
		listStatus int
		wantGone   bool
		wantError  bool
		wantListed bool
	}{
		{name: "readable vault", err: notFound, listStatus: http.StatusOK, wantGone: true, wantListed: true},
		{name: "name lookup miss", err: dvls.ErrEntryNotFound, listStatus: http.StatusOK, wantGone: true, wantListed: true},
		{name: "vault not found or not opened", err: notFound, listStatus: http.StatusNotFound, wantError: true, wantListed: true},
		{name: "vault content forbidden", err: notFound, listStatus: http.StatusForbidden, wantError: true, wantListed: true},
		{name: "vault check fails", err: notFound, listStatus: http.StatusInternalServerError, wantError: true, wantListed: true},
		{name: "not a 404", err: &dvls.RequestError{StatusCode: http.StatusForbidden, Err: errors.New("unexpected status code 403")}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			listed := false
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/vault/"+vaultId+"/entry", func(w http.ResponseWriter, r *http.Request) {
				listed = true
				if r.Method != http.MethodGet {
					t.Errorf("method = %s, want GET", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.listStatus)
				_, _ = w.Write([]byte(`{"currentPage":1,"pageSize":1,"totalCount":0,"totalPage":0,"data":[]}`))
			})
			server := httptest.NewServer(mux)
			t.Cleanup(server.Close)

			client, err := dvls.NewClientWithApiKey("test-api-key", server.URL)
			if err != nil {
				t.Fatal(err)
			}

			gone, diags := entryGone(context.Background(), &client, vaultId, tc.err)
			if gone != tc.wantGone {
				t.Fatalf("gone = %v, want %v", gone, tc.wantGone)
			}
			if diags.HasError() != tc.wantError {
				t.Fatalf("diagnostics error = %v, want %v: %v", diags.HasError(), tc.wantError, diags)
			}
			if listed != tc.wantListed {
				t.Fatalf("vault content checked = %v, want %v", listed, tc.wantListed)
			}
		})
	}
}
