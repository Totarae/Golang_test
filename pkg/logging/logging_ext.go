package logging

import (
	"fmt"
	"log"
	"strings"
)

type LoggerExt struct {
	*log.Logger
}

// Конструктор
func NewLogger() *LoggerExt {
	return &LoggerExt{log.Default()}
}

// Логируем с конкретными полями
func (l *LoggerExt) LogWithFields(level string, msg string, fields map[string]interface{}) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", level, msg))
	for k, v := range fields {
		sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
	}
	l.Printf(sb.String())
}

func (l *LoggerExt) LogErrorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.LogWithFields("ERROR", msg, nil)
}

func (l *LoggerExt) LogInfo(msg string, fields map[string]interface{}) {
	l.LogWithFields("INFO", msg, fields)
}
