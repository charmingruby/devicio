package observability

const (
	LOG_LEVEL_DEBUG   = "debug"
	LOG_LEVEL_INFO    = "info"
	LOG_LEVEL_WARN    = "warn"
	LOG_LEVEL_ERROR   = "error"
	LOG_LEVEL_DEFAULT = LOG_LEVEL_INFO
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
