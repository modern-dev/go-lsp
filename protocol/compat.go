// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

// This file provides backward-compatible aliases for symbol names used by go.lsp.dev/protocol v0.12.0.
// The generated code uses spec-faithful names  (e.g. CodeActionKindQuickFix, SemanticTokenTypesKeyword),
// but existing consumers - reference the shorter names from the old library.
// These aliases let callers simply rename their import path without touching every constant reference.

// ---------------------------------------------------------------------------
// ServerCapabilities sub-types
//
// go.lsp.dev/protocol v0.12.0 hand-wrote these under non-spec names.
// Our generated output uses the spec names from metaModel.json.
// ---------------------------------------------------------------------------

type (
	// ServerCapabilitiesWorkspace is an alias for the generated WorkspaceOptions.
	ServerCapabilitiesWorkspace = WorkspaceOptions

	// ServerCapabilitiesWorkspaceFolders is an alias for the generated WorkspaceFoldersServerCapabilities.
	ServerCapabilitiesWorkspaceFolders = WorkspaceFoldersServerCapabilities

	// ServerCapabilitiesWorkspaceFileOperations is an alias for the generated FileOperationOptions.
	ServerCapabilitiesWorkspaceFileOperations = FileOperationOptions
)

// ---------------------------------------------------------------------------
// Semantic token types — short names
// ---------------------------------------------------------------------------

const (
	SemanticTokenNamespace  = SemanticTokenTypesNamespace //nolint:revive
	SemanticTokenType       = SemanticTokenTypesType
	SemanticTokenClass      = SemanticTokenTypesClass
	SemanticTokenEnum       = SemanticTokenTypesEnum
	SemanticTokenInterface  = SemanticTokenTypesInterface
	SemanticTokenStruct     = SemanticTokenTypesStruct
	SemanticTokenTypeParam  = SemanticTokenTypesTypeParameter
	SemanticTokenParameter  = SemanticTokenTypesParameter
	SemanticTokenVariable   = SemanticTokenTypesVariable
	SemanticTokenProperty   = SemanticTokenTypesProperty
	SemanticTokenEnumMember = SemanticTokenTypesEnumMember
	SemanticTokenEvent      = SemanticTokenTypesEvent
	SemanticTokenFunction   = SemanticTokenTypesFunction
	SemanticTokenMethod     = SemanticTokenTypesMethod
	SemanticTokenMacro      = SemanticTokenTypesMacro
	SemanticTokenKeyword    = SemanticTokenTypesKeyword
	SemanticTokenModifier   = SemanticTokenTypesModifier
	SemanticTokenComment    = SemanticTokenTypesComment
	SemanticTokenString     = SemanticTokenTypesString
	SemanticTokenNumber     = SemanticTokenTypesNumber
	SemanticTokenRegexp     = SemanticTokenTypesRegexp
	SemanticTokenOperator   = SemanticTokenTypesOperator
	SemanticTokenDecorator  = SemanticTokenTypesDecorator
	SemanticTokenLabel      = SemanticTokenTypesLabel
)

// ---------------------------------------------------------------------------
// Semantic token modifiers — short names
// ---------------------------------------------------------------------------

const (
	SemanticTokenModifierDeclaration    = SemanticTokenModifiersDeclaration //nolint:revive
	SemanticTokenModifierDefinition     = SemanticTokenModifiersDefinition
	SemanticTokenModifierReadonly       = SemanticTokenModifiersReadonly
	SemanticTokenModifierStatic         = SemanticTokenModifiersStatic
	SemanticTokenModifierDeprecated     = SemanticTokenModifiersDeprecated
	SemanticTokenModifierAbstract       = SemanticTokenModifiersAbstract
	SemanticTokenModifierAsync          = SemanticTokenModifiersAsync
	SemanticTokenModifierModification   = SemanticTokenModifiersModification
	SemanticTokenModifierDocumentation  = SemanticTokenModifiersDocumentation
	SemanticTokenModifierDefaultLibrary = SemanticTokenModifiersDefaultLibrary
)

// ---------------------------------------------------------------------------
// MarkupKind — short names
// ---------------------------------------------------------------------------

const (
	PlainText = MarkupKindPlainText //nolint:revive
	Markdown  = MarkupKindMarkdown  //nolint:revive
)

// ---------------------------------------------------------------------------
// CodeActionKind — short names
// ---------------------------------------------------------------------------

const (
	QuickFix              = CodeActionKindQuickFix              //nolint:revive
	SourceOrganizeImports = CodeActionKindSourceOrganizeImports //nolint:revive
)

// ---------------------------------------------------------------------------
// FoldingRangeKind — short names
// ---------------------------------------------------------------------------

const (
	CommentFoldingRange = FoldingRangeKindComment //nolint:revive
	ImportsFoldingRange = FoldingRangeKindImports //nolint:revive
	RegionFoldingRange  = FoldingRangeKindRegion
)

// ---------------------------------------------------------------------------
// TextDocumentContentChangeEvent — concrete struct
//
// The LSP spec defines this as a union type, so the generated code emits
// `type TextDocumentContentChangeEvent = any`. Consumers that relied on the
// old go.lsp.dev/protocol struct need a concrete type with Range/Text fields.
// ---------------------------------------------------------------------------

// ContentChangeEvent is a concrete representation of an incremental or full
// text document content change event. Use this when you need to access
// Range and Text fields directly instead of the `any` type alias.
type ContentChangeEvent struct {
	// Range of the document that changed. Nil for full-content replacements.
	Range *Range `json:"range,omitempty"`
	// RangeLength is the optional length of the range being replaced.
	RangeLength uint32 `json:"rangeLength,omitempty"`
	// Text is the new text for the provided range, or the full document.
	Text string `json:"text"`
}
