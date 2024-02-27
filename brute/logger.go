package brute

import (
	"errors"

	"go.uber.org/zap"
)

// Predefined errors for specific conditions.
var (
	// errVerboseDebugFalse is returned when attempting to set a logger without enabling verbose or debug mode.
	errVerboseDebugFalse = errors.New("you can't set a logger if verbose and debug is set to false")
)

// SetDebug updates the debug state of a Brute instance.
func (b *Brute) SetDebug(debug bool) {
	b.Debug = debug
}

// SetLogger assigns a structured logger to the Brute instance, returning an error if both verbose and debug modes are disabled.
func (b *Brute) SetLogger(logger zap.Logger) error {
	// Prevent setting a logger if neither verbose nor debug mode is enabled.
	if !b.Verbose && !b.Debug {
		return errVerboseDebugFalse
	}
	b.logger = logger
	return nil
}

// LogError logs an error message with additional context, using the configured logger.
func (b *Brute) LogError(msg string, err error) {
	// Log the error along with any relevant RabbitMQ configuration details.
	b.logger.Error(msg, zap.Error(err), zap.Any("rabbit", b.rabbit))
}

// LogInfo logs an informational message, if verbose mode is enabled.
func (b *Brute) LogInfo(msg string, fields ...zap.Field) {
	if b.Verbose {
		// Log the message with any additional structured fields provided.
		b.logger.Info(msg, fields...)
	}
}

// LogDebug logs a debug message, if debug mode is enabled.
func (b *Brute) LogDebug(msg string, fields ...zap.Field) {
	if b.Debug {
		// Log the message with any additional structured fields provided.
		b.logger.Debug(msg, fields...)
	}
}
