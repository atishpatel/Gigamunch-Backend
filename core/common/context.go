package common

type contextString string

const (
	// LoggingKey is the context key for the logging client.
	LoggingKey = contextString("logging")
	// DBKey is the context key for the db client.
	// DBKey = contextString("db")
)
