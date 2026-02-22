// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package generate

import (
	"fmt"
	"strings"
	"unicode"
)

type (
	// Generator holds the parsed model and lookup indices used during code generation.
	Generator struct {
		Model *Model

		// Lookup indices built from the model.
		structs  map[string]*Structure
		enums    map[string]*Enumeration
		aliases  map[string]*TypeAlias
		requests map[string]*Request
		notifs   map[string]*Notification

		// namedLiterals tracks anonymous literal types that are promoted to
		// named Go structs. Key is the generated name, value is the literal type.
		namedLiterals map[string]*LiteralType

		// literalCounter disambiguate anonymous literal names.
		literalCounter int
	}

	abbreviation struct {
		mixed string
		upper string
	}
)

var (
	abbreviationPrefixes = []abbreviation{ //nolint:gochecknoglobals
		{"Uri", "URI"},
		{"Id", "ID"},
		{"Json", "JSON"},
		{"Utf", "UTF"},
		{"Lsp", "LSP"},
		{"Url", "URL"},
		{"Html", "HTML"},
		{"Css", "CSS"},
	}
	fieldNameOverrides = map[string]string{ //nolint:gochecknoglobals
		"uri":          "URI",
		"id":           "ID",
		"jsonrpc":      "JSONRPC",
		"documentUri":  "DocumentURI",
		"baseUri":      "BaseURI",
		"rootUri":      "RootURI",
		"resourceUri":  "ResourceURI",
		"oldUri":       "OldURI",
		"newUri":       "NewURI",
		"scopeUri":     "ScopeURI",
		"textDocument": "TextDocument",
	}
)

// NewGenerator creates a Generator from a parsed Model, building all lookup
// indices needed for type resolution.
func NewGenerator(model *Model) *Generator {
	gen := &Generator{ //nolint:exhaustruct
		Model:         model,
		structs:       make(map[string]*Structure, len(model.Structures)),
		enums:         make(map[string]*Enumeration, len(model.Enumerations)),
		aliases:       make(map[string]*TypeAlias, len(model.TypeAliases)),
		requests:      make(map[string]*Request, len(model.Requests)),
		notifs:        make(map[string]*Notification, len(model.Notifications)),
		namedLiterals: make(map[string]*LiteralType),
	}

	for idx := range model.Structures {
		gen.structs[model.Structures[idx].Name] = &model.Structures[idx]
	}

	for idx := range model.Enumerations {
		gen.enums[model.Enumerations[idx].Name] = &model.Enumerations[idx]
	}

	for idx := range model.TypeAliases {
		gen.aliases[model.TypeAliases[idx].Name] = &model.TypeAliases[idx]
	}

	for idx := range model.Requests {
		gen.requests[model.Requests[idx].Method] = &model.Requests[idx]
	}

	for idx := range model.Notifications {
		gen.notifs[model.Notifications[idx].Method] = &model.Notifications[idx]
	}

	return gen
}

// resolveGoType converts an LSP Type into its Go type string representation.
// Anonymous literal types are promoted to named structs and tracked in
// namedLiterals for later emission.
func (g *Generator) resolveGoType(typ *Type) string { //nolint:cyclop
	if typ == nil {
		return "any"
	}

	switch typ.Kind {
	case "base":
		return resolveBaseType(typ.Name)
	case "reference":
		return typ.Name
	case "array":
		return "[]" + g.resolveGoType(typ.Element)
	case "map":
		return "map[" + g.resolveGoType(typ.Key) + "]" + g.resolveGoType(typ.MapValue)
	case "or":
		return g.resolveUnion(typ.Items)
	case "and":
		return "any"
	case "tuple":
		return "any"
	case "literal":
		return g.promoteLiteral(typ.Literal)
	case "stringLiteral":
		return "string"
	case "integerLiteral":
		return "int32"
	case "booleanLiteral":
		return "bool"
	default:
		return "any"
	}
}

// resolveUnion converts an "or" (union) type into a Go type. The logic handles
// common LSP patterns:
//   - T | null → *T (nullable; pointer for structs/primitives, bare for slices/maps/any)
//   - Two non-null types or more → any
func (g *Generator) resolveUnion(items []Type) string {
	nonNull := make([]Type, 0, len(items))

	for _, item := range items {
		if item.Kind == "base" && item.Name == "null" {
			continue
		}

		nonNull = append(nonNull, item)
	}

	hasNull := len(nonNull) < len(items)

	if len(nonNull) == 1 {
		resolved := g.resolveGoType(&nonNull[0])
		if hasNull && needsPointerForNull(resolved) {
			return "*" + resolved
		}

		return resolved
	}

	return "any"
}

// promoteLiteral assigns a name to an anonymous literal type and registers it
// for later emission as a named Go struct.
func (g *Generator) promoteLiteral(lit *LiteralType) string {
	if lit == nil {
		return "any"
	}

	g.literalCounter++
	name := fmt.Sprintf("Literal%d", g.literalCounter)
	g.namedLiterals[name] = lit

	return name
}

// GoFieldName converts an LSP property name (camelCase) to a Go exported field
// name (PascalCase). It handles well-known abbreviation prefixes like "uri",
// "id", "json", etc.
func GoFieldName(lspName string) string {
	if lspName == "" {
		return ""
	}

	if mapped, ok := fieldNameOverrides[lspName]; ok {
		return mapped
	}

	runes := []rune(lspName)
	runes[0] = unicode.ToUpper(runes[0])

	result := string(runes)

	for _, prefix := range abbreviationPrefixes {
		if strings.HasPrefix(result, prefix.mixed) {
			rest := result[len(prefix.mixed):]
			if rest == "" || unicode.IsUpper([]rune(rest)[0]) {
				return prefix.upper + rest
			}
		}
	}

	return result
}

// GoMethodName converts an LSP method name like "textDocument/completion" to a
// Go method name like "Completion". For methods with a slash, the part after
// the last slash is used. For methods without a slash (like "initialize"), the
// name is simply capitalized.
//
// When the short name would collide (e.g. textDocument/diagnostic vs
// workspace/diagnostic), the caller should use GoMethodNameFull instead.
func GoMethodName(method string) string {
	name := method

	if idx := strings.LastIndex(method, "/"); idx >= 0 {
		name = method[idx+1:]
	}

	if name == "" {
		return ""
	}

	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// GoMethodNameFull converts an LSP method name into a fully-qualified Go
// method name by capitalizing each segment: "textDocument/diagnostic" becomes
// "TextDocumentDiagnostic", "$/cancelRequest" becomes "CancelRequest".
func GoMethodNameFull(method string) string {
	method = strings.TrimPrefix(method, "$/")
	parts := strings.Split(method, "/")

	var builder strings.Builder

	builder.Grow(len(method))

	for _, part := range parts {
		if part == "" {
			continue
		}

		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])

		builder.WriteString(string(runes))
	}

	return builder.String()
}

// GoEnumValueName produces a Go constant name from an enumeration name and a
// value name. For example, enum "CompletionItemKind" with value "Text"
// becomes "CompletionItemKindText".
func GoEnumValueName(enumName, valueName string) string {
	if valueName == "" {
		return enumName
	}

	runes := []rune(valueName)
	runes[0] = unicode.ToUpper(runes[0])

	return enumName + string(runes)
}

// JSONTag returns the JSON struct tag for a field, adding omitempty for optional fields.
func JSONTag(lspName string, optional bool) string {
	if optional {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", lspName)
	}

	return fmt.Sprintf("`json:\"%s\"`", lspName)
}

// IsServerMethod reports whether the given request or notification is directed
// at the server (client→server or both directions).
func IsServerMethod(direction string) bool {
	return direction == "clientToServer" || direction == "both"
}

// IsClientMethod reports whether the given request or notification is directed
// at the client (server→client or both directions).
func IsClientMethod(direction string) bool {
	return direction == "serverToClient" || direction == "both"
}

// resolveBaseType maps LSP base type names to Go types.
func resolveBaseType(name string) string { //nolint:cyclop
	switch name {
	case "string", "RegExp":
		return "string"
	case "DocumentUri":
		return "DocumentURI"
	case "URI":
		return "URI"
	case "integer":
		return "int32"
	case "uinteger":
		return "uint32"
	case "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "null":
		return "any"
	case "LSPAny", "LSPObject":
		return "any"
	case "LSPArray":
		return "[]any"
	default:
		return "any"
	}
}

// needsPointerForNull reports whether the Go type needs a pointer wrapper to
// represent a nullable value. Slices, maps, and any already have nil as their
// zero value and don't need wrapping.
func needsPointerForNull(goType string) bool {
	if strings.HasPrefix(goType, "*") ||
		strings.HasPrefix(goType, "[]") ||
		strings.HasPrefix(goType, "map[") {
		return false
	}

	switch goType {
	case "any":
		return false
	default:
		return true
	}
}

// methodConstName converts an LSP method name (e.g. "textDocument/completion",
// "$/cancelRequest") into a Go constant name like "MethodTextDocumentCompletion"
// or "MethodCancelRequest".
func methodConstName(method string) string {
	clean := strings.TrimPrefix(method, "$/")
	parts := strings.Split(clean, "/")

	var builder strings.Builder

	builder.WriteString("Method")

	for _, part := range parts {
		if part == "" {
			continue
		}

		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])

		builder.WriteString(string(runes))
	}

	return builder.String()
}
