// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.lsp.dev/jsonrpc2"
)

func TestServerHandlerNilLogger(t *testing.T) {
	h := ServerHandler(&stubServer{}, nil)
	require.NotNil(t, h)
}

func TestServerHandlerWithLogger(t *testing.T) {
	h := ServerHandler(&stubServer{}, NopLogger())
	require.NotNil(t, h)
}

func TestServerDispatchInitialize(t *testing.T) {
	srv := &stubServer{}
	h := ServerHandler(srv, nil)

	params := InitializeParams{ProcessId: new(int32)}
	raw, _ := json.Marshal(params)
	req, _ := jsonrpc2.NewCall(jsonrpc2.NewNumberID(1), "initialize", json.RawMessage(raw))

	var replied bool
	var replyResult any
	replier := func(ctx context.Context, result any, err error) error {
		replied = true
		replyResult = result
		return nil
	}

	require.NoError(t, h(context.Background(), replier, req))
	assert.True(t, replied, "replier should have been called")
	assert.True(t, srv.initializeCalled, "Initialize should have been called")
	assert.NotNil(t, replyResult)
}

func TestServerDispatchDidOpen(t *testing.T) {
	srv := &stubServer{}
	h := ServerHandler(srv, nil)

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI: "file:///test.go", LanguageId: "go", Version: 1, Text: "package main",
		},
	}
	raw, _ := json.Marshal(params)
	notif, _ := jsonrpc2.NewNotification("textDocument/didOpen", json.RawMessage(raw))

	nopReplier := func(ctx context.Context, result any, err error) error { return nil }
	require.NoError(t, h(context.Background(), nopReplier, notif))
	assert.True(t, srv.didOpenCalled)
}

func TestServerDispatchHover(t *testing.T) {
	srv := &stubServer{}
	h := ServerHandler(srv, nil)

	params := HoverParams{
		TextDocument: TextDocumentIdentifier{URI: "file:///test.go"},
		Position:     Position{Line: 1, Character: 5},
	}
	raw, _ := json.Marshal(params)
	req, _ := jsonrpc2.NewCall(jsonrpc2.NewNumberID(2), "textDocument/hover", json.RawMessage(raw))

	var replyResult any
	replier := func(ctx context.Context, result any, err error) error {
		replyResult = result
		return nil
	}

	require.NoError(t, h(context.Background(), replier, req))
	assert.True(t, srv.hoverCalled)

	hover, ok := replyResult.(*Hover)
	require.True(t, ok, "reply should be *Hover, got %T", replyResult)
	assert.Equal(t, "hello", hover.Contents)
}

func TestServerDispatchUnknownMethod(t *testing.T) {
	srv := &stubServer{}
	h := ServerHandler(srv, nil)

	req, _ := jsonrpc2.NewCall(
		jsonrpc2.NewNumberID(3),
		"custom/method",
		json.RawMessage(`{"key":"value"}`),
	)

	var replied bool
	replier := func(ctx context.Context, result any, err error) error {
		replied = true
		return nil
	}

	require.NoError(t, h(context.Background(), replier, req))
	assert.True(t, replied)
	assert.True(t, srv.requestCalled, "Request catch-all should have been called")
	assert.Equal(t, "custom/method", srv.requestMethod)
}

func TestServerDispatchInvalidParams(t *testing.T) {
	h := ServerHandler(&stubServer{}, nil)

	req, _ := jsonrpc2.NewCall(jsonrpc2.NewNumberID(4), "initialize", json.RawMessage(`not json`))

	var replyErr error
	replier := func(ctx context.Context, result any, err error) error {
		replyErr = err
		return nil
	}

	_ = h(context.Background(), replier, req)
	assert.Error(t, replyErr, "should reply with parse error for invalid params")
}

func TestServerDispatchShutdown(t *testing.T) {
	srv := &stubServer{}
	h := ServerHandler(srv, nil)

	req, _ := jsonrpc2.NewCall(jsonrpc2.NewNumberID(5), "shutdown", nil)

	var replied bool
	replier := func(ctx context.Context, result any, err error) error {
		replied = true
		return nil
	}

	require.NoError(t, h(context.Background(), replier, req))
	assert.True(t, replied)
	assert.True(t, srv.shutdownCalled)
}
