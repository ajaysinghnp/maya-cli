package metadata

type Logger interface {
	Info(string)
	Debug(string)
	Warn(string)
	Error(string)
}
