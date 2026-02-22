// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

package protocol

type (
	// Logger defines a leveled logging interface for the protocol package.
	// The interface mirrors zap.Logger's method names (Debug, Info, Warn, Error)
	// but uses variadic any instead of zap.Field, so no zap dependency is needed.
	// A thin adapter can bridge *zap.Logger if structured fields are desired.
	// Callers that do not need logging can pass NopLogger().
	Logger interface {
		Debug(msg string, fields ...any)
		Info(msg string, fields ...any)
		Warn(msg string, fields ...any)
		Error(msg string, fields ...any)
	}

	// nopLogger is a Logger that silently discards all log output.
	nopLogger struct{}
)

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

// NopLogger returns a Logger that silently discards all log output.
// Use this when no logging is desired (equivalent to zap.NewNop()).
func NopLogger() Logger { //nolint:ireturn
	return nopLogger{}
}
