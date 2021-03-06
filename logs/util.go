package logs

import (
	"bytes"
	"fmt"
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

func writeField(buffer *bytes.Buffer, field, sep string) {
	buffer.WriteString(field)
	buffer.WriteString(sep)
}

func (l *LogData) Bytes() []byte {
	var buffer bytes.Buffer
	levelText := getLevelText(l.level)

	writeField(&buffer, l.timeStr, SpaceSep)
	writeField(&buffer, levelText, SpaceSep)
	writeField(&buffer, l.serviceName, SpaceSep)

	writeField(&buffer, l.filename, ColonSep)
	writeField(&buffer, fmt.Sprintf("%d", l.lineNo), SpaceSep)
	writeField(&buffer, l.traceId, SpaceSep)

	if l.level == LogLevelAccess && l.fields != nil {
		for _, field := range l.fields.kvs {
			writeField(&buffer, fmt.Sprintf("%v:%v", field.key, field.val), SpaceSep)
		}
	}

	writeField(&buffer, l.message, LineSep)

	return buffer.Bytes()
}