package observability

import "log/slog"

const (
	LOG_LEVEL_DEBUG   = "debug"
	LOG_LEVEL_INFO    = "info"
	LOG_LEVEL_WARN    = "warn"
	LOG_LEVEL_ERROR   = "error"
	LOG_LEVEL_DEFAULT = LOG_LEVEL_INFO
)

type Logger = slog.Logger
