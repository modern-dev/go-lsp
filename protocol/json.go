// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import "encoding/json"

// Marshal serializes v to JSON.
// This is a thin wrapper around the standard library's json.Marshal.
// It exists so that the protocol package has a single point where the
// JSON encoder can be swapped (e.g. for a faster third-party encoder)
// without changing callers.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v) //nolint:wrapcheck
}

// Unmarshal deserializes data into v.
// This is a thin wrapper around the standard library's json.Unmarshal.
// It exists so that the protocol package has a single point where the
// JSON decoder can be swapped without changing callers.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v) //nolint:wrapcheck
}
