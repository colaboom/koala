package logs

const (
	LogLevelDebug LogLevel = iota
	LogLevelTrace
	LogLevelAccess
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

const (
	DefaultLogChanSize = 20000
)

type LogLevel int

func getLevelText(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelAccess:
		return "ACCESS"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	}

	return "UNKNOWN"
}

func getLevel(levelText string) LogLevel {
	switch levelText {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "access":
		return LogLevelAccess
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	}

	return LogLevelDebug
}
