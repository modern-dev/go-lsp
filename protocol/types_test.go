// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypesJSONRoundTrip_Position(t *testing.T) {
	orig := Position{Line: 42, Character: 7}
	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got Position
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, orig, got)
}

func TestTypesJSONRoundTrip_Range(t *testing.T) {
	orig := Range{
		Start: Position{Line: 1, Character: 0},
		End:   Position{Line: 1, Character: 10},
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got Range
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, orig, got)
}

func TestTypesJSONRoundTrip_Location(t *testing.T) {
	orig := Location{
		URI: "file:///test.go",
		Range: Range{
			Start: Position{Line: 5, Character: 0},
			End:   Position{Line: 5, Character: 15},
		},
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got Location
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, orig, got)
}

func TestTypesJSONRoundTrip_TextDocumentItem(t *testing.T) {
	orig := TextDocumentItem{
		URI:        "file:///hello.go",
		LanguageId: "go",
		Version:    3,
		Text:       "package main\n\nfunc main() {}",
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got TextDocumentItem
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, orig.URI, got.URI)
	assert.Equal(t, orig.LanguageId, got.LanguageId)
	assert.Equal(t, orig.Version, got.Version)
	assert.Equal(t, orig.Text, got.Text)
}

func TestTypesJSONRoundTrip_InitializeResult(t *testing.T) {
	orig := InitializeResult{
		Capabilities: ServerCapabilities{},
		ServerInfo: &ServerInfo{
			Name:    "test-server",
			Version: new("1.0.0"),
		},
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got InitializeResult
	require.NoError(t, json.Unmarshal(data, &got))
	require.NotNil(t, got.ServerInfo)
	assert.Equal(t, "test-server", got.ServerInfo.Name)
	require.NotNil(t, got.ServerInfo.Version)
	assert.Equal(t, "1.0.0", *got.ServerInfo.Version)
}

func TestTypesJSONRoundTrip_Hover(t *testing.T) {
	orig := Hover{
		Contents: MarkupContent{
			Kind:  MarkupKindMarkdown,
			Value: "# Hello\nWorld",
		},
		Range: &Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 0, Character: 5},
		},
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got Hover
	require.NoError(t, json.Unmarshal(data, &got))
	require.NotNil(t, got.Range)
	assert.Equal(t, uint32(0), got.Range.Start.Line)
	assert.Equal(t, uint32(5), got.Range.End.Character)

	contentsMap, ok := got.Contents.(map[string]any)
	require.True(t, ok, "Contents should be map[string]any, got %T", got.Contents)
	assert.Equal(t, string(MarkupKindMarkdown), contentsMap["kind"])
}

func TestTypesJSONRoundTrip_Diagnostic(t *testing.T) {
	orig := Diagnostic{
		Range: Range{
			Start: Position{Line: 10, Character: 0},
			End:   Position{Line: 10, Character: 20},
		},
		Severity: new(DiagnosticSeverityError),
		Message:  "undefined variable",
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got Diagnostic
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "undefined variable", got.Message)
	require.NotNil(t, got.Severity)
	assert.Equal(t, DiagnosticSeverityError, *got.Severity)
}

func TestTypesJSONRoundTrip_CompletionItem(t *testing.T) {
	orig := CompletionItem{
		Label: "myFunc",
		Kind:  new(CompletionItemKindFunction),
	}

	data, err := json.Marshal(orig)
	require.NoError(t, err)

	var got CompletionItem
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "myFunc", got.Label)
	require.NotNil(t, got.Kind)
	assert.Equal(t, CompletionItemKindFunction, *got.Kind)
}

func TestMethodConstants(t *testing.T) {
	assert.Equal(t, "initialize", MethodInitialize)
	assert.Equal(t, "initialized", MethodInitialized)
	assert.Equal(t, "shutdown", MethodShutdown)
	assert.Equal(t, "exit", MethodExit)
	assert.Equal(t, "textDocument/didOpen", MethodTextDocumentDidOpen)
	assert.Equal(t, "textDocument/didChange", MethodTextDocumentDidChange)
	assert.Equal(t, "textDocument/didClose", MethodTextDocumentDidClose)
	assert.Equal(t, "textDocument/didSave", MethodTextDocumentDidSave)
	assert.Equal(t, "textDocument/hover", MethodTextDocumentHover)
	assert.Equal(t, "textDocument/completion", MethodTextDocumentCompletion)
	assert.Equal(t, "textDocument/definition", MethodTextDocumentDefinition)
	assert.Equal(t, "textDocument/references", MethodTextDocumentReferences)
	assert.Equal(t, "textDocument/codeAction", MethodTextDocumentCodeAction)
	assert.Equal(t, "textDocument/codeLens", MethodTextDocumentCodeLens)
	assert.Equal(t, "textDocument/formatting", MethodTextDocumentFormatting)
	assert.Equal(t, "textDocument/rename", MethodTextDocumentRename)
	assert.Equal(t, "textDocument/signatureHelp", MethodTextDocumentSignatureHelp)
	assert.Equal(t, "textDocument/documentSymbol", MethodTextDocumentDocumentSymbol)
	assert.Equal(t, "textDocument/foldingRange", MethodTextDocumentFoldingRange)
	assert.Equal(t, "textDocument/documentLink", MethodTextDocumentDocumentLink)
	assert.Equal(t, "textDocument/documentHighlight", MethodTextDocumentDocumentHighlight)
	assert.Equal(t, "textDocument/semanticTokens/full", MethodTextDocumentSemanticTokensFull)
	assert.Equal(t, "textDocument/inlayHint", MethodTextDocumentInlayHint)
	assert.Equal(t, "workspace/symbol", MethodWorkspaceSymbol)
	assert.Equal(t, "workspace/executeCommand", MethodWorkspaceExecuteCommand)
}

func TestEnumerationValues(t *testing.T) {
	assert.Equal(t, DiagnosticSeverity(1), DiagnosticSeverityError)
	assert.Equal(t, DiagnosticSeverity(2), DiagnosticSeverityWarning)
	assert.Equal(t, DiagnosticSeverity(3), DiagnosticSeverityInformation)
	assert.Equal(t, DiagnosticSeverity(4), DiagnosticSeverityHint)

	assert.Equal(t, TextDocumentSyncKind(0), TextDocumentSyncKindNone)
	assert.Equal(t, TextDocumentSyncKind(1), TextDocumentSyncKindFull)
	assert.Equal(t, TextDocumentSyncKind(2), TextDocumentSyncKindIncremental)

	assert.Equal(t, CompletionItemKind(3), CompletionItemKindFunction)
}

func TestOptionalPointerFields(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		si := ServerInfo{Name: "srv", Version: new("2.0.0")}

		data, err := json.Marshal(si)
		require.NoError(t, err)

		var got ServerInfo
		require.NoError(t, json.Unmarshal(data, &got))
		require.NotNil(t, got.Version)
		assert.Equal(t, "2.0.0", *got.Version)
	})

	t.Run("omitted", func(t *testing.T) {
		si := ServerInfo{Name: "srv"}
		data, _ := json.Marshal(si)

		var got ServerInfo
		_ = json.Unmarshal(data, &got)
		assert.Nil(t, got.Version)
	})
}
