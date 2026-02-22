// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package generate

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type (
	// GeneratedOutput holds the generated Go source files.
	GeneratedOutput struct {
		Types  []byte // types_gen.go
		Server []byte // server_gen.go
		Client []byte // client_gen.go
	}

	// methodInfo describes a single method on the Server or Client interface.
	methodInfo struct {
		method    string // LSP method name, e.g. "textDocument/completion"
		goName    string // Go method name, e.g. "Completion"
		signature string // Go method signature
		doc       string
		isRequest bool

		paramsType string // Go type for params, empty if none
		resultType string // Go type for result, empty if notification
	}
)

// Generate produces all generated source files from the loaded model.
func (g *Generator) Generate() (*GeneratedOutput, error) {
	out := &GeneratedOutput{} //nolint:exhaustruct

	var err error

	out.Types, err = g.generateTypes()
	if err != nil {
		return nil, fmt.Errorf("generate types: %w", err)
	}

	out.Server, err = g.generateServer()
	if err != nil {
		return nil, fmt.Errorf("generate server: %w", err)
	}

	out.Client, err = g.generateClient()
	if err != nil {
		return nil, fmt.Errorf("generate client: %w", err)
	}

	return out, nil
}

// generateTypes emits types_gen.go containing all structures, enumerations,
// type aliases, and promoted literal types.
func (g *Generator) generateTypes() ([]byte, error) { //nolint:gocognit,cyclop,funlen,unparam
	var buf bytes.Buffer

	buf.Grow(256 * 1024) //nolint:mnd
	g.writeHeader(&buf, "protocol", "encoding/json")

	for _, strc := range g.Model.Structures {
		if strc.Proposed {
			continue
		}

		writeDoc(&buf, strc.Documentation, strc.Name)

		_, _ = fmt.Fprintf(&buf, "type %s struct {\n", strc.Name)
		props := g.collectProperties(&strc)

		for _, prop := range props {
			if prop.Proposed {
				continue
			}

			writeFieldDoc(&buf, prop.Documentation)

			goType := optionalType(g.resolveGoType(&prop.Type), prop.Optional)
			_, _ = fmt.Fprintf(
				&buf,
				"\t%s %s %s\n",
				GoFieldName(prop.Name),
				goType,
				JSONTag(prop.Name, prop.Optional),
			)
		}

		_, _ = fmt.Fprintf(&buf, "}\n\n")
	}

	for _, enum := range g.Model.Enumerations {
		if enum.Proposed {
			continue
		}

		goType := resolveEnumBaseType(enum.Type)

		writeDoc(&buf, enum.Documentation, enum.Name)

		_, _ = fmt.Fprintf(&buf, "type %s %s\n\n", enum.Name, goType)
		_, _ = fmt.Fprintf(&buf, "const (\n")

		for _, val := range enum.Values {
			if val.Proposed {
				continue
			}

			writeFieldDoc(&buf, val.Documentation)

			constName := GoEnumValueName(enum.Name, val.Name)

			if goType == "string" {
				_, _ = fmt.Fprintf(&buf, "\t%s %s = %q\n", constName, enum.Name, val.Value)
			} else {
				_, _ = fmt.Fprintf(
					&buf,
					"\t%s %s = %v\n",
					constName,
					enum.Name,
					formatNumericValue(val.Value),
				)
			}
		}

		_, _ = fmt.Fprintf(&buf, ")\n\n")
	}

	for _, alias := range g.Model.TypeAliases {
		if alias.Proposed {
			continue
		}

		writeDoc(&buf, alias.Documentation, alias.Name)
		goType := g.resolveGoType(&alias.Type)
		_, _ = fmt.Fprintf(&buf, "type %s = %s\n\n", alias.Name, goType)
	}

	if len(g.namedLiterals) > 0 {
		names := make([]string, 0, len(g.namedLiterals))
		for name := range g.namedLiterals {
			names = append(names, name)
		}

		slices.Sort(names)

		for _, name := range names {
			lit := g.namedLiterals[name]
			_, _ = fmt.Fprintf(&buf, "type %s struct {\n", name)

			for _, prop := range lit.Properties {
				if prop.Proposed {
					continue
				}

				writeFieldDoc(&buf, prop.Documentation)
				goType := optionalType(g.resolveGoType(&prop.Type), prop.Optional)
				_, _ = fmt.Fprintf(
					&buf,
					"\t%s %s %s\n",
					GoFieldName(prop.Name),
					goType,
					JSONTag(prop.Name, prop.Optional),
				)
			}

			_, _ = fmt.Fprintf(&buf, "}\n\n")
		}
	}

	buf.WriteString("// Ensure json import is used.\nvar _ = json.RawMessage{}\n")

	return buf.Bytes(), nil
}

// generateServer emits server_gen.go containing the Server interface and the
// dispatch function (serverDispatch).
func (g *Generator) generateServer() ([]byte, error) { //nolint:funlen,unparam
	var buf bytes.Buffer

	buf.Grow(40 * 1024) //nolint:mnd

	g.writeHeader(&buf, "protocol",
		"context",
		"encoding/json",
		"go.lsp.dev/jsonrpc2",
	)

	// Emit method name constants for all server methods.
	serverMethods := g.collectServerMethods()
	clientMethods := g.collectClientMethods()

	buf.WriteString("// LSP method name constants.\n")
	buf.WriteString("const (\n")

	emitted := make(map[string]bool)

	for _, m := range serverMethods {
		constName := methodConstName(m.method)
		if constName != "" && !emitted[constName] {
			emitted[constName] = true
			_, _ = fmt.Fprintf(&buf, "\t%s = %q\n", constName, m.method)
		}
	}

	for _, m := range clientMethods {
		constName := methodConstName(m.method)
		if constName != "" && !emitted[constName] {
			emitted[constName] = true
			_, _ = fmt.Fprintf(&buf, "\t%s = %q\n", constName, m.method)
		}
	}

	buf.WriteString(")\n\n")

	buf.WriteString("// Server defines the interface for an LSP server.\n")
	buf.WriteString("// All methods correspond to LSP requests and notifications\n")
	buf.WriteString("// directed from client to server.\n")
	buf.WriteString("type Server interface {\n")

	for _, m := range serverMethods {
		writeMethodDoc(&buf, m.doc, m.goName, m.method)
		_, _ = fmt.Fprintf(&buf, "\t%s\n", m.signature)
	}

	buf.WriteString("\n")
	buf.WriteString("\t// Request is a catch-all handler for any LSP method not covered by the\n")
	buf.WriteString("\t// interface above.  The method string is the raw LSP method name and\n")
	buf.WriteString("\t// params is the already-decoded parameter value.\n")
	buf.WriteString("\tRequest(ctx context.Context, method string, params any) (any, error)\n")
	buf.WriteString("}\n\n")

	buf.WriteString(
		"// serverDispatch dispatches a JSON-RPC request to the appropriate Server method.\n",
	)
	buf.WriteString(
		"func serverDispatch(ctx context.Context, server Server, reply jsonrpc2.Replier, req jsonrpc2.Request) error {\n",
	)
	buf.WriteString("\tswitch req.Method() {\n")

	for _, meth := range serverMethods {
		_, _ = fmt.Fprintf(&buf, "\tcase %q:\n", meth.method)

		if meth.isRequest {
			writeRequestDispatch(&buf, &meth)
		} else {
			writeNotificationDispatch(&buf, &meth)
		}
	}

	buf.WriteString("\tdefault:\n")
	buf.WriteString("\t\tvar params any\n")
	buf.WriteString("\t\tif req.Params() != nil {\n")
	buf.WriteString("\t\t\tif err := json.Unmarshal(req.Params(), &params); err != nil {\n")
	buf.WriteString("\t\t\t\treturn replyParseError(ctx, reply, err)\n")
	buf.WriteString("\t\t\t}\n")
	buf.WriteString("\t\t}\n")
	buf.WriteString("\t\tresp, err := server.Request(ctx, req.Method(), params)\n")
	buf.WriteString("\t\treturn reply(ctx, resp, err)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	return buf.Bytes(), nil
}

// generateClient emits client_gen.go containing the Client interface and the
// clientDispatcher implementation.
func (g *Generator) generateClient() ([]byte, error) { //nolint:unparam
	var buf bytes.Buffer

	buf.Grow(10 * 1024) //nolint:mnd

	g.writeHeader(&buf, "protocol",
		"context",
		"go.lsp.dev/jsonrpc2",
	)

	buf.WriteString("// Client defines the interface for an LSP client.\n")
	buf.WriteString("// All methods correspond to LSP requests and notifications\n")
	buf.WriteString("// directed from server to client.\n")
	buf.WriteString("type Client interface {\n")

	clientMethods := g.collectClientMethods()
	for _, m := range clientMethods {
		writeMethodDoc(&buf, m.doc, m.goName, m.method)
		_, _ = fmt.Fprintf(&buf, "\t%s\n", m.signature)
	}

	buf.WriteString("}\n\n")

	buf.WriteString("type clientDispatcher struct {\n")
	buf.WriteString("\tconn jsonrpc2.Conn\n")
	buf.WriteString("\tlogger Logger\n")
	buf.WriteString("}\n\n")

	buf.WriteString(
		"// ClientDispatcher returns a Client that dispatches LSP requests/notifications\n",
	)
	buf.WriteString("// across the given jsonrpc2 connection.\n")
	buf.WriteString("//\n")
	buf.WriteString(
		"// The logger parameter is used for protocol-level logging. Pass NopLogger()\n",
	)
	buf.WriteString("// (or nil) to disable logging.\n")
	buf.WriteString("func ClientDispatcher(conn jsonrpc2.Conn, logger Logger) Client {\n")
	buf.WriteString("\tif logger == nil {\n")
	buf.WriteString("\t\tlogger = NopLogger()\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn &clientDispatcher{conn: conn, logger: logger}\n")
	buf.WriteString("}\n\n")

	for _, m := range clientMethods {
		writeClientMethod(&buf, &m)
	}

	return buf.Bytes(), nil
}

// collectServerMethods returns all methods that belong on the Server interface
// (clientToServer and both directions), sorted by method name.
func (g *Generator) collectServerMethods() []methodInfo {
	var methods []methodInfo

	for _, r := range g.Model.Requests {
		if r.Proposed || !IsServerMethod(r.MessageDirection) {
			continue
		}

		methods = append(methods, g.buildRequestMethod(&r))
	}

	for _, n := range g.Model.Notifications {
		if n.Proposed || !IsServerMethod(n.MessageDirection) {
			continue
		}

		methods = append(methods, g.buildNotificationMethod(&n))
	}

	disambiguateMethods(methods)

	slices.SortFunc(methods, func(a, b methodInfo) int {
		return cmp.Compare(a.method, b.method)
	})

	return methods
}

// collectClientMethods returns all methods that belong on the Client interface
// (serverToClient and both directions), sorted by method name.
func (g *Generator) collectClientMethods() []methodInfo {
	var methods []methodInfo

	for _, r := range g.Model.Requests {
		if r.Proposed || !IsClientMethod(r.MessageDirection) {
			continue
		}

		methods = append(methods, g.buildRequestMethod(&r))
	}

	for _, n := range g.Model.Notifications {
		if n.Proposed || !IsClientMethod(n.MessageDirection) {
			continue
		}

		methods = append(methods, g.buildNotificationMethod(&n))
	}

	disambiguateMethods(methods)

	slices.SortFunc(methods, func(a, b methodInfo) int {
		return cmp.Compare(a.method, b.method)
	})

	return methods
}

// disambiguateMethods detects Go name collisions and switches colliding entries
// to their fully-qualified names, unless a preferred name is specified in
// methodNameOverrides. Overridden methods are pinned to their override name
// and never renamed by the collision resolver.
func disambiguateMethods(methods []methodInfo) {
	pinned := make(map[int]bool, len(methods))

	// Apply overrides first: some methods keep legacy short names for
	// backward compatibility with go.lsp.dev/protocol v0.12.0.
	for idx := range methods {
		if override, ok := methodNameOverrides[methods[idx].method]; ok {
			methods[idx].signature = strings.Replace(
				methods[idx].signature,
				methods[idx].goName+"(",
				override+"(",
				1,
			)
			methods[idx].goName = override
			pinned[idx] = true
		}
	}

	counts := make(map[string]int, len(methods))
	for _, m := range methods {
		counts[m.goName]++
	}

	for idx := range methods {
		if pinned[idx] {
			continue
		}

		if counts[methods[idx].goName] > 1 {
			fullName := GoMethodNameFull(methods[idx].method)
			methods[idx].signature = strings.Replace(
				methods[idx].signature,
				methods[idx].goName+"(",
				fullName+"(",
				1,
			)
			methods[idx].goName = fullName
		}
	}
}

// methodNameOverrides maps LSP method strings to preferred Go method names.
// These ensure backward compatibility with the names used by
// go.lsp.dev/protocol v0.12.0 when the default short name would collide
// with newer LSP 3.17 methods (e.g. notebookDocument/* vs textDocument/*).
var methodNameOverrides = map[string]string{ //nolint:gosec,gochecknoglobals
	// textDocument/ notifications — keep the v0.12.0 short names.
	// These collide with notebookDocument/* methods added in LSP 3.17.
	"textDocument/didOpen":   "DidOpen",
	"textDocument/didClose":  "DidClose",
	"textDocument/didChange": "DidChange",
	"textDocument/didSave":   "DidSave",

	// Resolve methods — v0.12.0 used TypeResolve names.
	// These collide with each other under the short name "Resolve".
	"completionItem/resolve":  "CompletionResolve",
	"codeLens/resolve":        "CodeLensResolve",
	"documentLink/resolve":    "DocumentLinkResolve",
	"codeAction/resolve":      "CodeActionResolve",
	"inlayHint/resolve":       "InlayHintResolve",
	"workspaceSymbol/resolve": "WorkspaceSymbolResolve",

	// textDocument/diagnostic vs workspace/diagnostic — keep short name
	// for the common textDocument variant.
	"textDocument/diagnostic": "Diagnostic",

	// textDocument/foldingRange — old protocol used plural "FoldingRanges".
	"textDocument/foldingRange": "FoldingRanges",

	// workspace/symbol — old protocol used plural "Symbols".
	"workspace/symbol": "Symbols",

	// Semantic tokens — old protocol prefixed with SemanticTokens.
	"textDocument/semanticTokens/full":       "SemanticTokensFull",
	"textDocument/semanticTokens/full/delta": "SemanticTokensFullDelta",
	"textDocument/semanticTokens/range":      "SemanticTokensRange",

	// window/workDoneProgress/cancel — old protocol named this WorkDoneProgressCancel.
	"window/workDoneProgress/cancel": "WorkDoneProgressCancel",
}

func (g *Generator) buildRequestMethod(req *Request) methodInfo {
	goName := GoMethodName(req.Method)
	paramsType := g.resolveMethodType(req.Params)
	resultType := g.resolveMethodType(req.Result)

	var sig string

	switch {
	case paramsType != "" && resultType != "":
		sig = fmt.Sprintf(
			"%s(ctx context.Context, params %s) (%s, error)",
			goName,
			paramsType,
			resultType,
		)
	case paramsType != "":
		sig = fmt.Sprintf("%s(ctx context.Context, params %s) error", goName, paramsType)
	case resultType != "":
		sig = fmt.Sprintf("%s(ctx context.Context) (%s, error)", goName, resultType)
	default:
		sig = goName + "(ctx context.Context) error"
	}

	return methodInfo{
		method:     req.Method,
		goName:     goName,
		signature:  sig,
		doc:        req.Documentation,
		isRequest:  true,
		paramsType: paramsType,
		resultType: resultType,
	}
}

func (g *Generator) buildNotificationMethod(notif *Notification) methodInfo {
	goName := GoMethodName(notif.Method)
	paramsType := g.resolveMethodType(notif.Params)

	var sig string
	if paramsType != "" {
		sig = fmt.Sprintf("%s(ctx context.Context, params %s) error", goName, paramsType)
	} else {
		sig = goName + "(ctx context.Context) error"
	}

	return methodInfo{ //nolint:exhaustruct
		method:     notif.Method,
		goName:     goName,
		signature:  sig,
		doc:        notif.Documentation,
		isRequest:  false,
		paramsType: paramsType,
	}
}

// resolveMethodType resolves a method parameter or result Type to its Go
// representation. Struct types are returned as pointers.
func (g *Generator) resolveMethodType(t *Type) string {
	if t == nil {
		return ""
	}

	resolved := g.resolveGoType(t)
	if resolved == "any" {
		return "any"
	}

	if _, ok := g.structs[resolved]; ok {
		return "*" + resolved
	}

	return resolved
}

// collectProperties gathers all properties for a structure, including inherited
// ones from Extends and Mixins. Uses a visited set to prevent infinite
// recursion if the spec ever contains cycles.
func (g *Generator) collectProperties(s *Structure) []Property {
	visited := make(map[string]bool)

	return g.collectPropertiesImpl(s, visited)
}

func (g *Generator) collectPropertiesImpl( //nolint:gocognit,cyclop
	structure *Structure,
	visited map[string]bool,
) []Property {
	if visited[structure.Name] {
		return nil
	}

	visited[structure.Name] = true
	seen := make(map[string]bool)

	var result []Property

	for _, p := range structure.Properties {
		if !seen[p.Name] {
			seen[p.Name] = true

			result = append(result, p)
		}
	}

	for _, ext := range structure.Extends {
		if ext.Kind == "reference" {
			if base, ok := g.structs[ext.Name]; ok {
				for _, p := range g.collectPropertiesImpl(base, visited) {
					if !seen[p.Name] {
						seen[p.Name] = true

						result = append(result, p)
					}
				}
			}
		}
	}

	for _, mixin := range structure.Mixins {
		if mixin.Kind == "reference" {
			if base, ok := g.structs[mixin.Name]; ok {
				for _, p := range g.collectPropertiesImpl(base, visited) {
					if !seen[p.Name] {
						seen[p.Name] = true

						result = append(result, p)
					}
				}
			}
		}
	}

	return result
}

// optionalType wraps goType in a pointer if the field is optional and the
// type doesn't already represent nil natively.
func optionalType(goType string, optional bool) string {
	if !optional {
		return goType
	}

	if needsPointerForNull(goType) {
		return "*" + goType
	}

	return goType
}

// resolveEnumBaseType maps the enumeration's base type to a Go type.
func resolveEnumBaseType(t EnumBaseType) string {
	switch t.Name {
	case "string":
		return "string"
	case "integer":
		return "int32"
	case "uinteger":
		return "uint32"
	default:
		return "string"
	}
}

// formatNumericValue formats a numeric value from metaModel.json (which JSON
// decodes as float64) to an integer string.
func formatNumericValue(v any) string {
	switch val := v.(type) {
	case float64:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// writeHeader writes the standard file header with package declaration, code
// generation notice, and imports.
func (g *Generator) writeHeader(buf *bytes.Buffer, pkg string, imports ...string) {
	_, _ = fmt.Fprintf(buf, "// Copyright %d Bohdan Shtepan.\n", time.Now().Year())
	buf.WriteString("// Licensed under the MIT License.\n\n")
	buf.WriteString("// Code generated by go-lsp/cmd/generate; DO NOT EDIT.\n")
	_, _ = fmt.Fprintf(buf, "// LSP version: %s\n\n", g.Model.MetaData.Version)
	_, _ = fmt.Fprintf(buf, "package %s\n\n", pkg)

	if len(imports) > 0 {
		buf.WriteString("import (\n")

		for _, imp := range imports {
			_, _ = fmt.Fprintf(buf, "\t%q\n", imp)
		}

		buf.WriteString(")\n\n")
	}
}

func writeDoc(buf *bytes.Buffer, doc, name string) {
	if doc != "" {
		for line := range strings.SplitSeq(strings.TrimSpace(doc), "\n") {
			_, _ = fmt.Fprintf(buf, "// %s\n", strings.TrimSpace(line))
		}
	} else {
		_, _ = fmt.Fprintf(buf, "// %s is an LSP type.\n", name)
	}
}

func writeFieldDoc(buf *bytes.Buffer, doc string) {
	if doc == "" {
		return
	}

	lines := strings.SplitSeq(strings.TrimSpace(doc), "\n")
	for line := range lines {
		_, _ = fmt.Fprintf(buf, "\t// %s\n", strings.TrimSpace(line))
	}
}

func writeMethodDoc(buf *bytes.Buffer, doc, goName, method string) {
	if doc != "" {
		for line := range strings.SplitSeq(strings.TrimSpace(doc), "\n") {
			_, _ = fmt.Fprintf(buf, "\t// %s\n", strings.TrimSpace(line))
		}
	} else {
		_, _ = fmt.Fprintf(buf, "\t// %s handles the %q method.\n", goName, method)
	}
}

// writeRequestDispatch writes the dispatch case for a request (expects a response).
func writeRequestDispatch(buf *bytes.Buffer, info *methodInfo) {
	if info.paramsType != "" {
		bareType := strings.TrimPrefix(info.paramsType, "*")
		_, _ = fmt.Fprintf(buf, "\t\tvar params %s\n", bareType)
		buf.WriteString("\t\tif err := json.Unmarshal(req.Params(), &params); err != nil {\n")
		buf.WriteString("\t\t\treturn replyParseError(ctx, reply, err)\n")
		buf.WriteString("\t\t}\n")
	}

	switch {
	case info.paramsType != "" && info.resultType != "":
		_, _ = fmt.Fprintf(buf, "\t\tresult, err := server.%s(ctx, &params)\n", info.goName)
		buf.WriteString("\t\treturn reply(ctx, result, err)\n")
	case info.paramsType != "":
		_, _ = fmt.Fprintf(buf, "\t\terr := server.%s(ctx, &params)\n", info.goName)
		buf.WriteString("\t\treturn reply(ctx, nil, err)\n")
	case info.resultType != "":
		_, _ = fmt.Fprintf(buf, "\t\tresult, err := server.%s(ctx)\n", info.goName)
		buf.WriteString("\t\treturn reply(ctx, result, err)\n")
	default:
		_, _ = fmt.Fprintf(buf, "\t\terr := server.%s(ctx)\n", info.goName)
		buf.WriteString("\t\treturn reply(ctx, nil, err)\n")
	}
}

// writeNotificationDispatch writes the dispatch case for a notification (no response).
func writeNotificationDispatch(buf *bytes.Buffer, info *methodInfo) {
	if info.paramsType != "" {
		bareType := strings.TrimPrefix(info.paramsType, "*")
		_, _ = fmt.Fprintf(buf, "\t\tvar params %s\n", bareType)
		buf.WriteString("\t\tif err := json.Unmarshal(req.Params(), &params); err != nil {\n")
		buf.WriteString("\t\t\treturn replyParseError(ctx, reply, err)\n")
		buf.WriteString("\t\t}\n")
		_, _ = fmt.Fprintf(buf, "\t\treturn server.%s(ctx, &params)\n", info.goName)
	} else {
		_, _ = fmt.Fprintf(buf, "\t\treturn server.%s(ctx)\n", info.goName)
	}
}

// writeClientMethod writes a single clientDispatcher method implementation.
func writeClientMethod(buf *bytes.Buffer, info *methodInfo) {
	_, _ = fmt.Fprintf(buf, "func (c *clientDispatcher) %s {\n", info.signature)

	if info.isRequest { //nolint:nestif
		if info.resultType != "" {
			bareResult := strings.TrimPrefix(info.resultType, "*")
			isPtr := strings.HasPrefix(info.resultType, "*")

			_, _ = fmt.Fprintf(buf, "\tvar result %s\n", bareResult)

			if info.paramsType != "" {
				_, _ = fmt.Fprintf(
					buf,
					"\t_, err := c.conn.Call(ctx, %q, params, &result)\n",
					info.method,
				)
			} else {
				_, _ = fmt.Fprintf(
					buf,
					"\t_, err := c.conn.Call(ctx, %q, nil, &result)\n",
					info.method,
				)
			}

			buf.WriteString("\tif err != nil {\n")

			if isPtr {
				buf.WriteString("\t\treturn nil, err\n")
			} else {
				_, _ = fmt.Fprintf(buf, "\t\tvar zero %s\n", bareResult)
				buf.WriteString("\t\treturn zero, err\n")
			}

			buf.WriteString("\t}\n")

			if isPtr {
				buf.WriteString("\treturn &result, nil\n")
			} else {
				buf.WriteString("\treturn result, nil\n")
			}
		} else {
			if info.paramsType != "" {
				_, _ = fmt.Fprintf(
					buf,
					"\t_, err := c.conn.Call(ctx, %q, params, nil)\n",
					info.method,
				)
			} else {
				_, _ = fmt.Fprintf(buf, "\t_, err := c.conn.Call(ctx, %q, nil, nil)\n", info.method)
			}

			buf.WriteString("\treturn err\n")
		}
	} else {
		if info.paramsType != "" {
			_, _ = fmt.Fprintf(buf, "\treturn c.conn.Notify(ctx, %q, params)\n", info.method)
		} else {
			_, _ = fmt.Fprintf(buf, "\treturn c.conn.Notify(ctx, %q, nil)\n", info.method)
		}
	}

	buf.WriteString("}\n\n")
}
