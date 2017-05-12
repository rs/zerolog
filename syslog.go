// +build !windows

package zerolog

import (
	"log/syslog"
)

type syslogWriter struct {
	w *syslog.Writer
}

// SyslogWriter wraps a syslog.Writer and set the right syslog level
// matching the log even level.
func SyslogWriter(w *syslog.Writer) LevelWriter {
	return syslogWriter{w}
}

func (sw syslogWriter) Write(p []byte) (n int, err error) {
	return sw.w.Write(p)
}

// WriteLevel implements LevelWriter interface.
func (sw syslogWriter) WriteLevel(level Level, p []byte) (n int, err error) {
	switch level {
	case DebugLevel:
		err = sw.w.Debug(string(p))
	case InfoLevel:
		err = sw.w.Info(string(p))
	case WarnLevel:
		err = sw.w.Warning(string(p))
	case ErrorLevel:
		err = sw.w.Err(string(p))
	case FatalLevel:
		err = sw.w.Emerg(string(p))
	case PanicLevel:
		err = sw.w.Crit(string(p))
	default:
		panic("invalid level")
	}
	n = len(p)
	return
}
