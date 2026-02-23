// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"context"
	"fmt"
)

// stubServer is a minimal Server implementation for testing.
// It tracks which methods were called and returns canned responses.
type stubServer struct {
	initializeCalled bool
	didOpenCalled    bool
	hoverCalled      bool
	shutdownCalled   bool
	requestCalled    bool
	requestMethod    string
}

func (s *stubServer) CancelRequest(_ context.Context, _ *CancelParams) error { return nil }
func (s *stubServer) Progress(_ context.Context, _ *ProgressParams) error    { return nil }
func (s *stubServer) SetTrace(_ context.Context, _ *SetTraceParams) error    { return nil }

func (s *stubServer) IncomingCalls(
	_ context.Context,
	_ *CallHierarchyIncomingCallsParams,
) ([]CallHierarchyIncomingCall, error) {
	return nil, nil
}

func (s *stubServer) OutgoingCalls(
	_ context.Context,
	_ *CallHierarchyOutgoingCallsParams,
) ([]CallHierarchyOutgoingCall, error) {
	return nil, nil
}

func (s *stubServer) CodeActionResolve(_ context.Context, params *CodeAction) (*CodeAction, error) {
	return params, nil
}

func (s *stubServer) CodeLensResolve(_ context.Context, params *CodeLens) (*CodeLens, error) {
	return params, nil
}

func (s *stubServer) CompletionResolve(
	_ context.Context,
	params *CompletionItem,
) (*CompletionItem, error) {
	return params, nil
}

func (s *stubServer) DocumentLinkResolve(
	_ context.Context,
	params *DocumentLink,
) (*DocumentLink, error) {
	return params, nil
}
func (s *stubServer) Exit(_ context.Context) error { return nil }
func (s *stubServer) Initialize(_ context.Context, _ *InitializeParams) (*InitializeResult, error) {
	s.initializeCalled = true
	return &InitializeResult{
		Capabilities: ServerCapabilities{},
		ServerInfo: &ServerInfo{
			Name:    "stub-server",
			Version: new("0.1.0-test"),
		},
	}, nil
}
func (s *stubServer) Initialized(_ context.Context, _ *InitializedParams) error { return nil }
func (s *stubServer) InlayHintResolve(_ context.Context, params *InlayHint) (*InlayHint, error) {
	return params, nil
}

func (s *stubServer) NotebookDocumentDidChange(
	_ context.Context,
	_ *DidChangeNotebookDocumentParams,
) error {
	return nil
}

func (s *stubServer) NotebookDocumentDidClose(
	_ context.Context,
	_ *DidCloseNotebookDocumentParams,
) error {
	return nil
}

func (s *stubServer) NotebookDocumentDidOpen(
	_ context.Context,
	_ *DidOpenNotebookDocumentParams,
) error {
	return nil
}

func (s *stubServer) NotebookDocumentDidSave(
	_ context.Context,
	_ *DidSaveNotebookDocumentParams,
) error {
	return nil
}

func (s *stubServer) Shutdown(_ context.Context) (any, error) {
	s.shutdownCalled = true
	return nil, nil
}

func (s *stubServer) CodeAction(_ context.Context, _ *CodeActionParams) ([]any, error) {
	return nil, nil
}

func (s *stubServer) CodeLens(_ context.Context, _ *CodeLensParams) ([]CodeLens, error) {
	return nil, nil
}

func (s *stubServer) ColorPresentation(
	_ context.Context,
	_ *ColorPresentationParams,
) ([]ColorPresentation, error) {
	return nil, nil
}

func (s *stubServer) Completion(_ context.Context, _ *CompletionParams) (any, error) {
	return nil, nil
}

func (s *stubServer) Declaration(_ context.Context, _ *DeclarationParams) (any, error) {
	return nil, nil
}

func (s *stubServer) Definition(_ context.Context, _ *DefinitionParams) (any, error) {
	return nil, nil
}

func (s *stubServer) Diagnostic(
	_ context.Context,
	_ *DocumentDiagnosticParams,
) (DocumentDiagnosticReport, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *stubServer) DidChange(_ context.Context, _ *DidChangeTextDocumentParams) error {
	return nil
}

func (s *stubServer) DidClose(_ context.Context, _ *DidCloseTextDocumentParams) error {
	return nil
}

func (s *stubServer) DidOpen(_ context.Context, _ *DidOpenTextDocumentParams) error {
	s.didOpenCalled = true
	return nil
}

func (s *stubServer) DidSave(_ context.Context, _ *DidSaveTextDocumentParams) error {
	return nil
}

func (s *stubServer) DocumentColor(
	_ context.Context,
	_ *DocumentColorParams,
) ([]ColorInformation, error) {
	return nil, nil
}

func (s *stubServer) DocumentHighlight(
	_ context.Context,
	_ *DocumentHighlightParams,
) ([]DocumentHighlight, error) {
	return nil, nil
}

func (s *stubServer) DocumentLink(
	_ context.Context,
	_ *DocumentLinkParams,
) ([]DocumentLink, error) {
	return nil, nil
}

func (s *stubServer) DocumentSymbol(_ context.Context, _ *DocumentSymbolParams) (any, error) {
	return nil, nil
}

func (s *stubServer) FoldingRanges(
	_ context.Context,
	_ *FoldingRangeParams,
) ([]FoldingRange, error) {
	return nil, nil
}

func (s *stubServer) Formatting(
	_ context.Context,
	_ *DocumentFormattingParams,
) ([]TextEdit, error) {
	return nil, nil
}

func (s *stubServer) Hover(_ context.Context, params *HoverParams) (*Hover, error) {
	s.hoverCalled = true
	return &Hover{
		Contents: "hello",
		Range: &Range{
			Start: params.Position,
			End:   params.Position,
		},
	}, nil
}

func (s *stubServer) Implementation(_ context.Context, _ *ImplementationParams) (any, error) {
	return nil, nil
}

func (s *stubServer) InlayHint(_ context.Context, _ *InlayHintParams) ([]InlayHint, error) {
	return nil, nil
}

func (s *stubServer) InlineValue(_ context.Context, _ *InlineValueParams) ([]InlineValue, error) {
	return nil, nil
}

func (s *stubServer) LinkedEditingRange(
	_ context.Context,
	_ *LinkedEditingRangeParams,
) (*LinkedEditingRanges, error) {
	return nil, nil
}

func (s *stubServer) Moniker(_ context.Context, _ *MonikerParams) ([]Moniker, error) {
	return nil, nil
}

func (s *stubServer) OnTypeFormatting(
	_ context.Context,
	_ *DocumentOnTypeFormattingParams,
) ([]TextEdit, error) {
	return nil, nil
}

func (s *stubServer) PrepareCallHierarchy(
	_ context.Context,
	_ *CallHierarchyPrepareParams,
) ([]CallHierarchyItem, error) {
	return nil, nil
}

func (s *stubServer) PrepareRename(
	_ context.Context,
	_ *PrepareRenameParams,
) (*PrepareRenameResult, error) {
	return nil, nil
}

func (s *stubServer) PrepareTypeHierarchy(
	_ context.Context,
	_ *TypeHierarchyPrepareParams,
) ([]TypeHierarchyItem, error) {
	return nil, nil
}

func (s *stubServer) RangeFormatting(
	_ context.Context,
	_ *DocumentRangeFormattingParams,
) ([]TextEdit, error) {
	return nil, nil
}

func (s *stubServer) References(_ context.Context, _ *ReferenceParams) ([]Location, error) {
	return nil, nil
}

func (s *stubServer) Rename(_ context.Context, _ *RenameParams) (*WorkspaceEdit, error) {
	return nil, nil
}

func (s *stubServer) SelectionRange(
	_ context.Context,
	_ *SelectionRangeParams,
) ([]SelectionRange, error) {
	return nil, nil
}

func (s *stubServer) SemanticTokensFull(
	_ context.Context,
	_ *SemanticTokensParams,
) (*SemanticTokens, error) {
	return nil, nil
}

func (s *stubServer) SemanticTokensFullDelta(
	_ context.Context,
	_ *SemanticTokensDeltaParams,
) (any, error) {
	return nil, nil
}

func (s *stubServer) SemanticTokensRange(
	_ context.Context,
	_ *SemanticTokensRangeParams,
) (*SemanticTokens, error) {
	return nil, nil
}

func (s *stubServer) SignatureHelp(
	_ context.Context,
	_ *SignatureHelpParams,
) (*SignatureHelp, error) {
	return nil, nil
}

func (s *stubServer) TypeDefinition(_ context.Context, _ *TypeDefinitionParams) (any, error) {
	return nil, nil
}

func (s *stubServer) WillSave(_ context.Context, _ *WillSaveTextDocumentParams) error {
	return nil
}

func (s *stubServer) WillSaveWaitUntil(
	_ context.Context,
	_ *WillSaveTextDocumentParams,
) ([]TextEdit, error) {
	return nil, nil
}

func (s *stubServer) Subtypes(
	_ context.Context,
	_ *TypeHierarchySubtypesParams,
) ([]TypeHierarchyItem, error) {
	return nil, nil
}

func (s *stubServer) Supertypes(
	_ context.Context,
	_ *TypeHierarchySupertypesParams,
) ([]TypeHierarchyItem, error) {
	return nil, nil
}

func (s *stubServer) WorkDoneProgressCancel(
	_ context.Context,
	_ *WorkDoneProgressCancelParams,
) error {
	return nil
}

func (s *stubServer) WorkspaceDiagnostic(
	_ context.Context,
	_ *WorkspaceDiagnosticParams,
) (*WorkspaceDiagnosticReport, error) {
	return nil, nil
}

func (s *stubServer) DidChangeConfiguration(
	_ context.Context,
	_ *DidChangeConfigurationParams,
) error {
	return nil
}

func (s *stubServer) DidChangeWatchedFiles(
	_ context.Context,
	_ *DidChangeWatchedFilesParams,
) error {
	return nil
}

func (s *stubServer) DidChangeWorkspaceFolders(
	_ context.Context,
	_ *DidChangeWorkspaceFoldersParams,
) error {
	return nil
}

func (s *stubServer) DidCreateFiles(_ context.Context, _ *CreateFilesParams) error {
	return nil
}

func (s *stubServer) DidDeleteFiles(_ context.Context, _ *DeleteFilesParams) error {
	return nil
}

func (s *stubServer) DidRenameFiles(_ context.Context, _ *RenameFilesParams) error {
	return nil
}

func (s *stubServer) ExecuteCommand(_ context.Context, _ *ExecuteCommandParams) (*LSPAny, error) {
	return nil, nil
}

func (s *stubServer) Symbols(_ context.Context, _ *WorkspaceSymbolParams) (any, error) {
	return nil, nil
}

func (s *stubServer) WillCreateFiles(
	_ context.Context,
	_ *CreateFilesParams,
) (*WorkspaceEdit, error) {
	return nil, nil
}

func (s *stubServer) WillDeleteFiles(
	_ context.Context,
	_ *DeleteFilesParams,
) (*WorkspaceEdit, error) {
	return nil, nil
}

func (s *stubServer) WillRenameFiles(
	_ context.Context,
	_ *RenameFilesParams,
) (*WorkspaceEdit, error) {
	return nil, nil
}

func (s *stubServer) WorkspaceSymbolResolve(
	_ context.Context,
	params *WorkspaceSymbol,
) (*WorkspaceSymbol, error) {
	return params, nil
}

func (s *stubServer) Request(_ context.Context, method string, _ any) (any, error) {
	s.requestCalled = true
	s.requestMethod = method
	return map[string]string{"echo": method}, nil
}

// Verify stubServer implements Server at compile time.
var _ Server = (*stubServer)(nil)
