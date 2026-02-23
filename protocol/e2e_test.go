// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/modern-dev/go-lsp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.lsp.dev/jsonrpc2"
)

// e2eServer is a Server implementation for end-to-end tests.
type e2eServer struct {
	initialized bool
	opened      map[protocol.DocumentURI]string
}

func newE2EServer() *e2eServer {
	return &e2eServer{opened: make(map[protocol.DocumentURI]string)}
}

func (s *e2eServer) Initialize(
	_ context.Context,
	_ *protocol.InitializeParams,
) (*protocol.InitializeResult, error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: new(true),
				Change:    new(protocol.TextDocumentSyncKindFull),
			},
			HoverProvider:              true,
			CompletionProvider:         &protocol.CompletionOptions{},
			DefinitionProvider:         true,
			DocumentSymbolProvider:     true,
			CodeActionProvider:         true,
			DocumentFormattingProvider: true,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "e2e-test-server",
			Version: new("0.0.1-test"),
		},
	}, nil
}

func (s *e2eServer) Initialized(_ context.Context, _ *protocol.InitializedParams) error {
	s.initialized = true
	return nil
}

func (s *e2eServer) Shutdown(_ context.Context) (any, error) { return nil, nil }
func (s *e2eServer) Exit(_ context.Context) error            { return nil }

func (s *e2eServer) DidOpen(_ context.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.opened[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (s *e2eServer) DidChange(_ context.Context, _ *protocol.DidChangeTextDocumentParams) error {
	return nil
}

func (s *e2eServer) DidClose(_ context.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(s.opened, params.TextDocument.URI)
	return nil
}

func (s *e2eServer) DidSave(_ context.Context, _ *protocol.DidSaveTextDocumentParams) error {
	return nil
}

func (s *e2eServer) Hover(
	_ context.Context,
	params *protocol.HoverParams,
) (*protocol.Hover, error) {
	text, ok := s.opened[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind: protocol.MarkupKindMarkdown,
			Value: fmt.Sprintf(
				"Hovering over `%s` at %d:%d",
				text,
				params.Position.Line,
				params.Position.Character,
			),
		},
	}, nil
}

func (s *e2eServer) Completion(_ context.Context, _ *protocol.CompletionParams) (any, error) {
	return &protocol.CompletionList{
		IsIncomplete: false,
		Items: []protocol.CompletionItem{
			{Label: "fmt"},
			{Label: "func"},
		},
	}, nil
}

func (s *e2eServer) Definition(_ context.Context, params *protocol.DefinitionParams) (any, error) {
	return &protocol.Location{
		URI: params.TextDocument.URI,
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 10},
		},
	}, nil
}

func (s *e2eServer) DocumentSymbol(
	_ context.Context,
	_ *protocol.DocumentSymbolParams,
) (any, error) {
	return []protocol.DocumentSymbol{
		{
			Name: "main",
			Kind: protocol.SymbolKindFunction,
			Range: protocol.Range{
				Start: protocol.Position{Line: 2, Character: 0},
				End:   protocol.Position{Line: 4, Character: 1},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{Line: 2, Character: 5},
				End:   protocol.Position{Line: 2, Character: 9},
			},
		},
	}, nil
}

func (s *e2eServer) CodeAction(_ context.Context, _ *protocol.CodeActionParams) ([]any, error) {
	return nil, nil
}

func (s *e2eServer) Formatting(
	_ context.Context,
	_ *protocol.DocumentFormattingParams,
) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *e2eServer) Request(_ context.Context, method string, _ any) (any, error) {
	return map[string]string{"method": method}, nil
}

// --- Stubs for remaining Server interface methods ---

func (s *e2eServer) CancelRequest(context.Context, *protocol.CancelParams) error { return nil }
func (s *e2eServer) Progress(context.Context, *protocol.ProgressParams) error    { return nil }
func (s *e2eServer) SetTrace(context.Context, *protocol.SetTraceParams) error    { return nil }

func (s *e2eServer) IncomingCalls(
	context.Context,
	*protocol.CallHierarchyIncomingCallsParams,
) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, nil
}

func (s *e2eServer) OutgoingCalls(
	context.Context,
	*protocol.CallHierarchyOutgoingCallsParams,
) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, nil
}

func (s *e2eServer) CodeActionResolve(
	_ context.Context,
	p *protocol.CodeAction,
) (*protocol.CodeAction, error) {
	return p, nil
}

func (s *e2eServer) CodeLensResolve(
	_ context.Context,
	p *protocol.CodeLens,
) (*protocol.CodeLens, error) {
	return p, nil
}

func (s *e2eServer) CompletionResolve(
	_ context.Context,
	p *protocol.CompletionItem,
) (*protocol.CompletionItem, error) {
	return p, nil
}

func (s *e2eServer) DocumentLinkResolve(
	_ context.Context,
	p *protocol.DocumentLink,
) (*protocol.DocumentLink, error) {
	return p, nil
}

func (s *e2eServer) InlayHintResolve(
	_ context.Context,
	p *protocol.InlayHint,
) (*protocol.InlayHint, error) {
	return p, nil
}

func (s *e2eServer) NotebookDocumentDidChange(
	context.Context,
	*protocol.DidChangeNotebookDocumentParams,
) error {
	return nil
}

func (s *e2eServer) NotebookDocumentDidClose(
	context.Context,
	*protocol.DidCloseNotebookDocumentParams,
) error {
	return nil
}

func (s *e2eServer) NotebookDocumentDidOpen(
	context.Context,
	*protocol.DidOpenNotebookDocumentParams,
) error {
	return nil
}

func (s *e2eServer) NotebookDocumentDidSave(
	context.Context,
	*protocol.DidSaveNotebookDocumentParams,
) error {
	return nil
}

func (s *e2eServer) CodeLens(
	context.Context,
	*protocol.CodeLensParams,
) ([]protocol.CodeLens, error) {
	return nil, nil
}

func (s *e2eServer) ColorPresentation(
	context.Context,
	*protocol.ColorPresentationParams,
) ([]protocol.ColorPresentation, error) {
	return nil, nil
}

func (s *e2eServer) Declaration(context.Context, *protocol.DeclarationParams) (any, error) {
	return nil, nil
}

func (s *e2eServer) Diagnostic(
	context.Context,
	*protocol.DocumentDiagnosticParams,
) (protocol.DocumentDiagnosticReport, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *e2eServer) DocumentColor(
	context.Context,
	*protocol.DocumentColorParams,
) ([]protocol.ColorInformation, error) {
	return nil, nil
}

func (s *e2eServer) DocumentHighlight(
	context.Context,
	*protocol.DocumentHighlightParams,
) ([]protocol.DocumentHighlight, error) {
	return nil, nil
}

func (s *e2eServer) DocumentLink(
	context.Context,
	*protocol.DocumentLinkParams,
) ([]protocol.DocumentLink, error) {
	return nil, nil
}

func (s *e2eServer) FoldingRanges(
	context.Context,
	*protocol.FoldingRangeParams,
) ([]protocol.FoldingRange, error) {
	return nil, nil
}

func (s *e2eServer) Implementation(context.Context, *protocol.ImplementationParams) (any, error) {
	return nil, nil
}

func (s *e2eServer) InlayHint(
	context.Context,
	*protocol.InlayHintParams,
) ([]protocol.InlayHint, error) {
	return nil, nil
}

func (s *e2eServer) InlineValue(
	context.Context,
	*protocol.InlineValueParams,
) ([]protocol.InlineValue, error) {
	return nil, nil
}

func (s *e2eServer) LinkedEditingRange(
	context.Context,
	*protocol.LinkedEditingRangeParams,
) (*protocol.LinkedEditingRanges, error) {
	return nil, nil
}

func (s *e2eServer) Moniker(context.Context, *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, nil
}

func (s *e2eServer) OnTypeFormatting(
	context.Context,
	*protocol.DocumentOnTypeFormattingParams,
) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *e2eServer) PrepareCallHierarchy(
	context.Context,
	*protocol.CallHierarchyPrepareParams,
) ([]protocol.CallHierarchyItem, error) {
	return nil, nil
}

func (s *e2eServer) PrepareRename(
	context.Context,
	*protocol.PrepareRenameParams,
) (*protocol.PrepareRenameResult, error) {
	return nil, nil
}

func (s *e2eServer) PrepareTypeHierarchy(
	context.Context,
	*protocol.TypeHierarchyPrepareParams,
) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}

func (s *e2eServer) RangeFormatting(
	context.Context,
	*protocol.DocumentRangeFormattingParams,
) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *e2eServer) References(
	context.Context,
	*protocol.ReferenceParams,
) ([]protocol.Location, error) {
	return nil, nil
}

func (s *e2eServer) Rename(
	context.Context,
	*protocol.RenameParams,
) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *e2eServer) SelectionRange(
	context.Context,
	*protocol.SelectionRangeParams,
) ([]protocol.SelectionRange, error) {
	return nil, nil
}

func (s *e2eServer) SemanticTokensFull(
	context.Context,
	*protocol.SemanticTokensParams,
) (*protocol.SemanticTokens, error) {
	return nil, nil
}

func (s *e2eServer) SemanticTokensFullDelta(
	context.Context,
	*protocol.SemanticTokensDeltaParams,
) (any, error) {
	return nil, nil
}

func (s *e2eServer) SemanticTokensRange(
	context.Context,
	*protocol.SemanticTokensRangeParams,
) (*protocol.SemanticTokens, error) {
	return nil, nil
}

func (s *e2eServer) SignatureHelp(
	context.Context,
	*protocol.SignatureHelpParams,
) (*protocol.SignatureHelp, error) {
	return nil, nil
}

func (s *e2eServer) TypeDefinition(context.Context, *protocol.TypeDefinitionParams) (any, error) {
	return nil, nil
}

func (s *e2eServer) WillSave(context.Context, *protocol.WillSaveTextDocumentParams) error {
	return nil
}

func (s *e2eServer) WillSaveWaitUntil(
	context.Context,
	*protocol.WillSaveTextDocumentParams,
) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *e2eServer) Subtypes(
	context.Context,
	*protocol.TypeHierarchySubtypesParams,
) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}

func (s *e2eServer) Supertypes(
	context.Context,
	*protocol.TypeHierarchySupertypesParams,
) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}

func (s *e2eServer) WorkDoneProgressCancel(
	context.Context,
	*protocol.WorkDoneProgressCancelParams,
) error {
	return nil
}

func (s *e2eServer) WorkspaceDiagnostic(
	context.Context,
	*protocol.WorkspaceDiagnosticParams,
) (*protocol.WorkspaceDiagnosticReport, error) {
	return nil, nil
}

func (s *e2eServer) DidChangeConfiguration(
	context.Context,
	*protocol.DidChangeConfigurationParams,
) error {
	return nil
}

func (s *e2eServer) DidChangeWatchedFiles(
	context.Context,
	*protocol.DidChangeWatchedFilesParams,
) error {
	return nil
}

func (s *e2eServer) DidChangeWorkspaceFolders(
	context.Context,
	*protocol.DidChangeWorkspaceFoldersParams,
) error {
	return nil
}

func (s *e2eServer) DidCreateFiles(
	context.Context,
	*protocol.CreateFilesParams,
) error {
	return nil
}

func (s *e2eServer) DidDeleteFiles(
	context.Context,
	*protocol.DeleteFilesParams,
) error {
	return nil
}

func (s *e2eServer) DidRenameFiles(
	context.Context,
	*protocol.RenameFilesParams,
) error {
	return nil
}

func (s *e2eServer) ExecuteCommand(
	context.Context,
	*protocol.ExecuteCommandParams,
) (*protocol.LSPAny, error) {
	return nil, nil
}

func (s *e2eServer) Symbols(context.Context, *protocol.WorkspaceSymbolParams) (any, error) {
	return nil, nil
}

func (s *e2eServer) WillCreateFiles(
	context.Context,
	*protocol.CreateFilesParams,
) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *e2eServer) WillDeleteFiles(
	context.Context,
	*protocol.DeleteFilesParams,
) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *e2eServer) WillRenameFiles(
	context.Context,
	*protocol.RenameFilesParams,
) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *e2eServer) WorkspaceSymbolResolve(
	_ context.Context,
	p *protocol.WorkspaceSymbol,
) (*protocol.WorkspaceSymbol, error) {
	return p, nil
}

var _ protocol.Server = (*e2eServer)(nil)

// setupE2E creates a connected client ↔ server over an in-process pipe.
func setupE2E(t *testing.T) (context.Context, jsonrpc2.Conn, jsonrpc2.Conn, *e2eServer) {
	t.Helper()

	srv := newE2EServer()
	handler := protocol.ServerHandler(srv, nil)

	clientConn, serverConn := net.Pipe()

	serverStream := jsonrpc2.NewStream(serverConn)
	sConn := jsonrpc2.NewConn(serverStream)
	sConn.Go(context.Background(), handler)

	clientStream := jsonrpc2.NewStream(clientConn)
	cConn := jsonrpc2.NewConn(clientStream)
	cConn.Go(context.Background(), jsonrpc2.MethodNotFoundHandler)

	t.Cleanup(func() {
		_ = cConn.Close()
		_ = sConn.Close()
		<-cConn.Done()
		<-sConn.Done()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	return ctx, sConn, cConn, srv
}

func TestE2E_InitializeLifecycle(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	// 1. initialize
	var initResult protocol.InitializeResult
	_, err := clientConn.Call(ctx, "initialize", protocol.InitializeParams{
		ProcessId:    new(int32),
		Capabilities: protocol.ClientCapabilities{},
	}, &initResult)
	require.NoError(t, err)

	require.NotNil(t, initResult.ServerInfo)
	assert.Equal(t, "e2e-test-server", initResult.ServerInfo.Name)
	require.NotNil(t, initResult.ServerInfo.Version)
	assert.Equal(t, "0.0.1-test", *initResult.ServerInfo.Version)
	assert.Equal(t, true, initResult.Capabilities.HoverProvider)
	assert.Equal(t, true, initResult.Capabilities.DefinitionProvider)

	// 2. initialized
	require.NoError(t, clientConn.Notify(ctx, "initialized", protocol.InitializedParams{}))

	// 3. shutdown
	var shutdownResult any
	_, err = clientConn.Call(ctx, "shutdown", nil, &shutdownResult)
	require.NoError(t, err)

	// 4. exit
	require.NoError(t, clientConn.Notify(ctx, "exit", nil))
}

func TestE2E_TextDocumentDidOpen(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	err := clientConn.Notify(ctx, "textDocument/didOpen", protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI: "file:///workspace/main.go", LanguageId: "go", Version: 1, Text: "package main",
		},
	})
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	var hover protocol.Hover
	_, err = clientConn.Call(ctx, "textDocument/hover", protocol.HoverParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///workspace/main.go"},
		Position:     protocol.Position{Line: 0, Character: 8},
	}, &hover)
	require.NoError(t, err)

	contentsMap, ok := hover.Contents.(map[string]any)
	require.True(t, ok, "hover.Contents should be map[string]any, got %T", hover.Contents)

	val, ok := contentsMap["value"].(string)
	require.True(t, ok)
	assert.NotEmpty(t, val)
	assert.Contains(t, val, "package main")
}

func TestE2E_Completion(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	require.NoError(
		t,
		clientConn.Notify(ctx, "textDocument/didOpen", protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI: "file:///workspace/main.go", LanguageId: "go", Version: 1,
				Text: "package main\n\nfunc main() {\n\tf\n}",
			},
		}),
	)

	time.Sleep(50 * time.Millisecond)

	var result json.RawMessage
	_, err := clientConn.Call(ctx, "textDocument/completion", protocol.CompletionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///workspace/main.go"},
		Position:     protocol.Position{Line: 3, Character: 2},
	}, &result)
	require.NoError(t, err)

	var list protocol.CompletionList
	require.NoError(t, json.Unmarshal(result, &list))
	require.Len(t, list.Items, 2)
	assert.Equal(t, "fmt", list.Items[0].Label)
	assert.Equal(t, "func", list.Items[1].Label)
}

func TestE2E_Definition(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	var result json.RawMessage
	_, err := clientConn.Call(ctx, "textDocument/definition", protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///workspace/main.go"},
		Position:     protocol.Position{Line: 5, Character: 10},
	}, &result)
	require.NoError(t, err)

	var loc protocol.Location
	require.NoError(t, json.Unmarshal(result, &loc))
	assert.Equal(t, protocol.DocumentURI("file:///workspace/main.go"), loc.URI)
	assert.Equal(t, uint32(0), loc.Range.Start.Line)
	assert.Equal(t, uint32(10), loc.Range.End.Character)
}

func TestE2E_DocumentSymbol(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	var result json.RawMessage
	_, err := clientConn.Call(ctx, "textDocument/documentSymbol", protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///workspace/main.go"},
	}, &result)
	require.NoError(t, err)

	var symbols []protocol.DocumentSymbol
	require.NoError(t, json.Unmarshal(result, &symbols))
	require.Len(t, symbols, 1)
	assert.Equal(t, "main", symbols[0].Name)
	assert.Equal(t, protocol.SymbolKindFunction, symbols[0].Kind)
}

func TestE2E_CustomRequestCatchAll(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	var result json.RawMessage
	_, err := clientConn.Call(ctx, "custom/myMethod", map[string]string{"hello": "world"}, &result)
	require.NoError(t, err)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(result, &resp))
	assert.Equal(t, "custom/myMethod", resp["method"])
}

func TestE2E_InvalidParams(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	var result any
	_, err := clientConn.Call(ctx, "textDocument/hover", json.RawMessage(`not valid json`), &result)
	assert.Error(t, err)
}

func TestE2E_MultipleRequests(t *testing.T) {
	ctx, _, clientConn, _ := setupE2E(t)

	// Initialize
	var initResult protocol.InitializeResult
	_, err := clientConn.Call(ctx, "initialize", protocol.InitializeParams{
		ProcessId:    new(int32),
		Capabilities: protocol.ClientCapabilities{},
	}, &initResult)
	require.NoError(t, err)

	require.NoError(t, clientConn.Notify(ctx, "initialized", protocol.InitializedParams{}))

	// Open document
	require.NoError(
		t,
		clientConn.Notify(ctx, "textDocument/didOpen", protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI: "file:///test.go", LanguageId: "go", Version: 1, Text: "package test",
			},
		}),
	)

	time.Sleep(50 * time.Millisecond)

	// Fire 10 hover requests sequentially
	for i := range 10 {
		var hover protocol.Hover
		_, err := clientConn.Call(ctx, "textDocument/hover", protocol.HoverParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: "file:///test.go"},
			Position:     protocol.Position{Line: 0, Character: uint32(i)},
		}, &hover)
		require.NoError(t, err, "hover[%d]", i)
	}

	// Shutdown
	var shutdownResult any
	_, err = clientConn.Call(ctx, "shutdown", nil, &shutdownResult)
	require.NoError(t, err)
}
