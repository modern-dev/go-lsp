// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

// Package generate implements the LSP protocol code generator.
// It reads Microsoft's metaModel.json (the machine-readable LSP specification)
// and produces Go source files for the protocol package.
package generate

import "encoding/json"

type (
	// Model is the top-level structure of metaModel.json.
	Model struct {
		MetaData      MetaData       `json:"metaData"`
		Requests      []Request      `json:"requests"`
		Notifications []Notification `json:"notifications"`
		Structures    []Structure    `json:"structures"`
		Enumerations  []Enumeration  `json:"enumerations"`
		TypeAliases   []TypeAlias    `json:"typeAliases"`
	}

	// MetaData contains the LSP version information.
	MetaData struct {
		Version string `json:"version"`
	}

	// Request describes an LSP request (client→server or server→client).
	Request struct {
		Documentation       string `json:"documentation"`
		ErrorData           *Type  `json:"errorData"`
		MessageDirection    string `json:"messageDirection"`
		Method              string `json:"method"`
		Params              *Type  `json:"params"`
		PartialResult       *Type  `json:"partialResult"`
		Proposed            bool   `json:"proposed"`
		RegistrationMethod  string `json:"registrationMethod"`
		RegistrationOptions *Type  `json:"registrationOptions"`
		Result              *Type  `json:"result"`
		Since               string `json:"since"`
	}

	// Notification describes an LSP notification (no response expected).
	Notification struct {
		Documentation    string `json:"documentation"`
		MessageDirection string `json:"messageDirection"`
		Method           string `json:"method"`
		Params           *Type  `json:"params"`
		Proposed         bool   `json:"proposed"`
		Since            string `json:"since"`
	}

	// Structure describes a named LSP type (struct).
	Structure struct {
		Documentation string     `json:"documentation"`
		Extends       []Type     `json:"extends"`
		Mixins        []Type     `json:"mixins"`
		Name          string     `json:"name"`
		Properties    []Property `json:"properties"`
		Proposed      bool       `json:"proposed"`
		Since         string     `json:"since"`
	}

	// Enumeration describes an LSP enum type.
	Enumeration struct {
		Documentation        string             `json:"documentation"`
		Name                 string             `json:"name"`
		Since                string             `json:"since"`
		Proposed             bool               `json:"proposed"`
		SupportsCustomValues bool               `json:"supportsCustomValues"`
		Type                 EnumBaseType       `json:"type"`
		Values               []EnumerationValue `json:"values"`
	}

	// EnumBaseType is the underlying type of enumeration.
	EnumBaseType struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	}

	// EnumerationValue is a single value in an enumeration.
	EnumerationValue struct {
		Documentation string `json:"documentation"`
		Name          string `json:"name"`
		Proposed      bool   `json:"proposed"`
		Since         string `json:"since"`
		Value         any    `json:"value"`
	}

	// TypeAlias describes a named type alias in the LSP spec.
	TypeAlias struct {
		Documentation string `json:"documentation"`
		Name          string `json:"name"`
		Proposed      bool   `json:"proposed"`
		Since         string `json:"since"`
		Type          Type   `json:"type"`
	}

	// Property describes a single property (field) of a Structure or literal type.
	Property struct {
		Documentation string `json:"documentation"`
		Name          string `json:"name"`
		Optional      bool   `json:"optional"`
		Proposed      bool   `json:"proposed"`
		Since         string `json:"since"`
		Type          Type   `json:"type"`
	}

	// LiteralType represents the value of a "literal" kind Type — an anonymous
	// struct with its own properties.
	LiteralType struct {
		Properties []Property `json:"properties"`
	}

	// Type describes an LSP type. The Kind field discriminates which other fields
	// are populated:
	//   - "base"           → Name ("string", "integer", "uinteger", "boolean", "decimal", "null",
	//     "DocumentUri", "URI", "LSPAny", "LSPObject", "LSPArray")
	//   - "reference"      → Name (refers to a Structure, Enumeration, or TypeAlias)
	//   - "array"          → Element
	//   - "map"            → Key + MapValue
	//   - "or"             → Items (union)
	//   - "and"            → Items (intersection)
	//   - "tuple"          → Items
	//   - "literal"        → Literal (anonymous struct with properties)
	//   - "stringLiteral"  → StringValue
	//   - "integerLiteral" → IntValue
	//   - "booleanLiteral" → BoolValue
	//
	// The Value field in metaModel.json is polymorphic — it can be a Type (for map),
	// a LiteralType (for literal), or a primitive (for stringLiteral, etc.).
	// We handle this with a custom UnmarshalJSON.
	Type struct {
		Kind string `json:"kind"`
		Name string `json:"name,omitempty"`

		// Array
		Element *Type `json:"element,omitempty"`

		// Map
		Key      *Type `json:"key,omitempty"`
		MapValue *Type `json:"-"`

		// Union / intersection / tuple
		Items []Type `json:"items,omitempty"`

		// Literal (anonymous struct)
		Literal *LiteralType `json:"-"`

		// Literal values
		StringValue string `json:"-"`
		IntValue    int64  `json:"-"`
		BoolValue   bool   `json:"-"`
	}
)

// UnmarshalJSON handles the polymorphic "value" field based on Kind.
// Uses a single-pass decode into a struct with json.RawMessage for the
// polymorphic value field, avoiding a second full parse.
func (t *Type) UnmarshalJSON(data []byte) error { //nolint:cyclop
	type plain struct {
		Kind    string          `json:"kind"`
		Name    string          `json:"name,omitempty"`
		Element *Type           `json:"element,omitempty"`
		Key     *Type           `json:"key,omitempty"`
		Items   []Type          `json:"items,omitempty"`
		Value   json.RawMessage `json:"value,omitempty"`
	}

	var pln plain
	if err := json.Unmarshal(data, &pln); err != nil { //nolint:noinlineerr
		return err //nolint:wrapcheck
	}

	t.Kind = pln.Kind
	t.Name = pln.Name
	t.Element = pln.Element
	t.Key = pln.Key
	t.Items = pln.Items

	if len(pln.Value) == 0 {
		return nil
	}

	switch t.Kind {
	case "map":
		var mapVal Type
		if err := json.Unmarshal(pln.Value, &mapVal); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}

		t.MapValue = &mapVal

	case "literal":
		var lit LiteralType
		if err := json.Unmarshal(pln.Value, &lit); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}

		t.Literal = &lit

	case "stringLiteral":
		var s string
		if err := json.Unmarshal(pln.Value, &s); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}

		t.StringValue = s

	case "integerLiteral":
		var n int64
		if err := json.Unmarshal(pln.Value, &n); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}

		t.IntValue = n

	case "booleanLiteral":
		var b bool
		if err := json.Unmarshal(pln.Value, &b); err != nil { //nolint:noinlineerr
			return err //nolint:wrapcheck
		}

		t.BoolValue = b
	}

	return nil
}
