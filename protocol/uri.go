// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

// This file defines DocumentURI and URI, the two URI base types in the LSP
// specification. They are handwritten (not generated) because they carry
// semantic meaning beyond a plain string and benefit from helper methods.
// See https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#uri

import (
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	// DocumentURI represents the URI of a client editor document.
	// Over the wire it is transferred as a string, but this named type guarantees
	// that the contents can be parsed as a valid URI and provides helper methods
	// for path conversion.
	DocumentURI string

	// URI is a generic URI as defined by the LSP specification. Unlike
	// DocumentURI, it may refer to any scheme (http, https, etc.), not
	// just file URIs.
	URI string
)

// URIFromPath creates a DocumentURI from a filesystem path.
//
//	URIFromPath("/home/user/file.go") => "file:///home/user/file.go"
//	URIFromPath("C:\\Users\\file.go") => "file:///C:/Users/file.go"  (Windows)
func URIFromPath(path string) DocumentURI {
	if path == "" {
		return ""
	}

	// Normalize to forward slashes.
	path = filepath.ToSlash(path)

	// On Windows, paths like "C:/..." need a leading slash in the URI.
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	return DocumentURI("file://" + path)
}

// Path converts a DocumentURI to a filesystem path.
//
// If the URI is not a file URI or cannot be parsed, it returns the raw URI
// string unchanged.
func (u DocumentURI) Path() string {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return string(u)
	}

	if parsed.Scheme != "file" {
		return string(u)
	}

	path := parsed.Path
	if parsed.Host != "" {
		// UNC path: file://server/share/...
		path = "//" + parsed.Host + path
	}

	// On Windows, /C:/foo → C:\foo
	if runtime.GOOS == "windows" && len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}

	return filepath.FromSlash(path)
}

// Filename is an alias for Path. It returns the file path corresponding to
// the given DocumentURI.
func (u DocumentURI) Filename() string {
	return u.Path()
}

// IsFile reports whether the URI has a "file" scheme.
func (u DocumentURI) IsFile() bool {
	return strings.HasPrefix(string(u), "file://")
}
