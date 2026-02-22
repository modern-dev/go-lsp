// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"context"
	"fmt"

	"go.lsp.dev/jsonrpc2"
)

// LSP error codes, as defined in the LSP specification.
// These extend the JSON-RPC error codes.
const (
	// CodeServerNotInitialized is returned when a request is sent before the
	// server has received the "initialize" request.
	CodeServerNotInitialized int64 = -32002

	// CodeInvalidRequest is returned when the server receives a request that
	// is not valid in the current state.
	CodeInvalidRequest int64 = -32600

	// CodeMethodNotFound is returned when the method is not supported.
	CodeMethodNotFound int64 = -32601

	// CodeInvalidParams is returned when the parameters are invalid.
	CodeInvalidParams int64 = -32602

	// CodeInternalError is returned for internal server errors.
	CodeInternalError int64 = -32603

	// CodeParseError is returned when JSON parsing fails.
	CodeParseError int64 = -32700

	// CodeRequestCancelled is returned when the client cancels a request.
	CodeRequestCancelled int64 = -32800

	// CodeContentModified is returned when content was modified before the
	// request could complete.
	CodeContentModified int64 = -32801
)

// replyParseError sends a parse error reply. This is used by the generated
// dispatch code when JSON unmarshalling of parameters fails.
func replyParseError(ctx context.Context, reply jsonrpc2.Replier, err error) error {
	return reply(ctx, nil, fmt.Errorf("invalid params: %w", err))
}
