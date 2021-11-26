package logs

import (
	"runtime"
	"time"
)

type LogData struct {
	curTime     time.Time
	message     string
	timeStr     string
	level       LogLevel
	filename    string
	lineNo      int
	traceId     string
	serviceName string
	fields      *LogField
}

func GetLineInfo() (fileName string, lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	return
}