// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeAliases(t *testing.T) {
	var w ServerCapabilitiesWorkspace
	var wf ServerCapabilitiesWorkspaceFolders
	var fo ServerCapabilitiesWorkspaceFileOperations
	_ = w
	_ = wf
	_ = fo
}

func TestSemanticTokenTypeAliases(t *testing.T) {
	tests := []struct {
		alias    SemanticTokenTypes
		original SemanticTokenTypes
	}{
		{SemanticTokenNamespace, SemanticTokenTypesNamespace},
		{SemanticTokenType, SemanticTokenTypesType},
		{SemanticTokenClass, SemanticTokenTypesClass},
		{SemanticTokenEnum, SemanticTokenTypesEnum},
		{SemanticTokenInterface, SemanticTokenTypesInterface},
		{SemanticTokenStruct, SemanticTokenTypesStruct},
		{SemanticTokenTypeParam, SemanticTokenTypesTypeParameter},
		{SemanticTokenParameter, SemanticTokenTypesParameter},
		{SemanticTokenVariable, SemanticTokenTypesVariable},
		{SemanticTokenProperty, SemanticTokenTypesProperty},
		{SemanticTokenEnumMember, SemanticTokenTypesEnumMember},
		{SemanticTokenEvent, SemanticTokenTypesEvent},
		{SemanticTokenFunction, SemanticTokenTypesFunction},
		{SemanticTokenMethod, SemanticTokenTypesMethod},
		{SemanticTokenMacro, SemanticTokenTypesMacro},
		{SemanticTokenKeyword, SemanticTokenTypesKeyword},
		{SemanticTokenModifier, SemanticTokenTypesModifier},
		{SemanticTokenComment, SemanticTokenTypesComment},
		{SemanticTokenString, SemanticTokenTypesString},
		{SemanticTokenNumber, SemanticTokenTypesNumber},
		{SemanticTokenRegexp, SemanticTokenTypesRegexp},
		{SemanticTokenOperator, SemanticTokenTypesOperator},
		{SemanticTokenDecorator, SemanticTokenTypesDecorator},
		{SemanticTokenLabel, SemanticTokenTypesLabel},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.original, tt.alias)
	}
}

func TestSemanticTokenModifierAliases(t *testing.T) {
	tests := []struct {
		alias    SemanticTokenModifiers
		original SemanticTokenModifiers
	}{
		{SemanticTokenModifierDeclaration, SemanticTokenModifiersDeclaration},
		{SemanticTokenModifierDefinition, SemanticTokenModifiersDefinition},
		{SemanticTokenModifierReadonly, SemanticTokenModifiersReadonly},
		{SemanticTokenModifierStatic, SemanticTokenModifiersStatic},
		{SemanticTokenModifierDeprecated, SemanticTokenModifiersDeprecated},
		{SemanticTokenModifierAbstract, SemanticTokenModifiersAbstract},
		{SemanticTokenModifierAsync, SemanticTokenModifiersAsync},
		{SemanticTokenModifierModification, SemanticTokenModifiersModification},
		{SemanticTokenModifierDocumentation, SemanticTokenModifiersDocumentation},
		{SemanticTokenModifierDefaultLibrary, SemanticTokenModifiersDefaultLibrary},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.original, tt.alias)
	}
}

func TestMarkupKindAliases(t *testing.T) {
	assert.Equal(t, MarkupKindPlainText, PlainText)
	assert.Equal(t, MarkupKindMarkdown, Markdown)
}

func TestCodeActionKindAliases(t *testing.T) {
	assert.Equal(t, CodeActionKindQuickFix, QuickFix)
	assert.Equal(t, CodeActionKindSourceOrganizeImports, SourceOrganizeImports)
}

func TestFoldingRangeKindAliases(t *testing.T) {
	assert.Equal(t, FoldingRangeKindComment, CommentFoldingRange)
	assert.Equal(t, FoldingRangeKindImports, ImportsFoldingRange)
	assert.Equal(t, FoldingRangeKindRegion, RegionFoldingRange)
}
