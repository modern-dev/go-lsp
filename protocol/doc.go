// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

// Package protocol implements Go types and interfaces for the Language Server
// Protocol (LSP). Most of the code in this package is generated from Microsoft's
// metaModel.json specification.
//
// Generated files (DO NOT EDIT):
//   - types_gen.go  — structures, enumerations, type aliases
//   - server_gen.go — Server interface, method constants, dispatch
//   - client_gen.go — Client interface, ClientDispatcher
//
// Hand-written files:
//   - doc.go      — this file
//   - uri.go      — DocumentURI / URI types and helpers
//   - errors.go   — LSP error codes and helpers
//   - handler.go  — ServerHandler (adapts Server to jsonrpc2.Handler)
//   - logger.go   — Logger interface and NopLogger
//   - json.go     — JSON codec abstraction (Marshal / Unmarshal)
//   - compat.go   — backward-compatible aliases for go.lsp.dev/protocol v0.12.0
package protocol

//go:generate go run github.com/modern-dev/go-lsp/cmd/generate -o .
