package provider

import "testing"

func TestFolderFullPath(t *testing.T) {
	tests := []struct {
		name, parent, child string
		want                string
	}{
		{name: "root folder", parent: "", child: "X", want: "X"},
		{name: "single nesting", parent: "A", child: "X", want: `A\X`},
		{name: "two levels", parent: `A\B`, child: "X", want: `A\B\X`},
		{name: "empty child stays empty parent", parent: "", child: "", want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := folderFullPath(tc.parent, tc.child); got != tc.want {
				t.Fatalf("folderFullPath(%q, %q) = %q, want %q", tc.parent, tc.child, got, tc.want)
			}
		})
	}
}
