// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

import (
	"context"

	"go.lsp.dev/jsonrpc2"
)

// ServerHandler returns a jsonrpc2.Handler that dispatches incoming requests
// and notifications to the given Server implementation.
//
// The logger parameter is used for protocol-level logging.  Pass NopLogger()
// (or nil) to disable logging.
//
// Usage:
//
//	var s protocol.Server = &myServer{}
//	handler := protocol.ServerHandler(s, protocol.NopLogger())
//	conn := jsonrpc2.NewConn(stream)
//	conn.Go(ctx, handler)
func ServerHandler(server Server, logger Logger) jsonrpc2.Handler {
	if logger == nil {
		logger = NopLogger() //nolint:ineffassign,staticcheck,wastedassign
	}

	return func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		return serverDispatch(ctx, server, reply, req)
	}
}
