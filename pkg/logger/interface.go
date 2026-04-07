package logger

// Logger is the interface for logging operations
// This abstraction allows us to switch logger implementations without changing business logic
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, fields ...Field)

	// Info logs an info message
	Info(msg string, fields ...Field)

	// Warn logs a warning message
	Warn(msg string, fields ...Field)

	// Error logs an error message
	Error(msg string, fields ...Field)

	// Fatal logs a fatal message and exits
	Fatal(msg string, fields ...Field)

	// Sync flushes any buffered log entries
	Sync() error
}

// Field represents a log field (key-value pair)
type Field interface {
	Key() string
	Value() interface{}
}
