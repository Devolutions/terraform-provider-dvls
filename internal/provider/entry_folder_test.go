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

func TestFolderParentFromFullPath(t *testing.T) {
	tests := []struct {
		name, fullPath, child string
		wantParent            string
		wantOK                bool
	}{
		{name: "root folder, path equals name", fullPath: "X", child: "X", wantParent: "", wantOK: false},
		{name: "nested under one parent", fullPath: `A\X`, child: "X", wantParent: "A", wantOK: true},
		{name: "nested under two parents", fullPath: `A\B\X`, child: "X", wantParent: `A\B`, wantOK: true},
		{name: "name-suffix collision: parent ends in name", fullPath: "AX", child: "X", wantParent: "", wantOK: false},
		{name: "name appears inside parent path", fullPath: `X\Y`, child: "Y", wantParent: "X", wantOK: true},
		{name: "empty path", fullPath: "", child: "X", wantParent: "", wantOK: false},
		{name: "path is just a backslash + name", fullPath: `\X`, child: "X", wantParent: "", wantOK: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotParent, gotOK := folderParentFromFullPath(tc.fullPath, tc.child)
			if gotParent != tc.wantParent || gotOK != tc.wantOK {
				t.Fatalf("folderParentFromFullPath(%q, %q) = (%q, %v), want (%q, %v)",
					tc.fullPath, tc.child, gotParent, gotOK, tc.wantParent, tc.wantOK)
			}
		})
	}
}
