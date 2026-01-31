package iface

// Logger defines the common logging interface for all modules
type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	Success(msg string) // logs successful/completed actions
}
