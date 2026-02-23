// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURIFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want DocumentURI
	}{
		{"empty", "", ""},
		{"unix absolute", "/home/user/file.go", "file:///home/user/file.go"},
		{"unix root", "/", "file:///"},
		{"unix nested", "/a/b/c/d.txt", "file:///a/b/c/d.txt"},
		{
			"unix with spaces",
			"/home/user/my project/file.go",
			"file:///home/user/my project/file.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, URIFromPath(tt.path))
		})
	}
}

func TestDocumentURI_Path(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix-only path tests")
	}

	tests := []struct {
		name string
		uri  DocumentURI
		want string
	}{
		{"unix file URI", "file:///home/user/file.go", "/home/user/file.go"},
		{"non-file URI", "https://example.com", "https://example.com"},
		{"unparseable URI", DocumentURI([]byte{0x7f}), string([]byte{0x7f})},
		{"empty", "", ""},
		{"file URI root", "file:///", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.uri.Path())
		})
	}
}

func TestDocumentURI_Filename(t *testing.T) {
	uri := DocumentURI("file:///home/user/file.go")
	assert.Equal(t, uri.Path(), uri.Filename())
}

func TestDocumentURI_IsFile(t *testing.T) {
	tests := []struct {
		name string
		uri  DocumentURI
		want bool
	}{
		{"file URI", "file:///home/user/file.go", true},
		{"https URI", "https://example.com", false},
		{"empty", "", false},
		{"file prefix only", "file://", true},
		{"almost file", "file:/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.uri.IsFile())
		})
	}
}

func TestURIRoundTrip(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("round-trip test for unix paths only")
	}

	paths := []string{
		"/home/user/file.go",
		"/tmp/a.txt",
		"/",
	}

	for _, path := range paths {
		uri := URIFromPath(path)
		got := uri.Path()
		require.Equal(t, path, got, "round-trip failed for path %q via uri %q", path, uri)
	}
}
