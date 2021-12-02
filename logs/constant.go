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
	SpaceSep           = " "
	ColonSep           = ":"
	LineSep            = "\n"
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

func GetLogLevel(levelText string) LogLevel {
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

func getLevelColor(level LogLevel) Color {
	switch level {
	case LogLevelDebug:
		return White
	case LogLevelTrace:
		return Yellow
	case LogLevelAccess:
		return Green
	case LogLevelInfo:
		return Blue
	case LogLevelWarn:
		return Cyan
	case LogLevelError:
		return Red
	}

	return Magenta
}
