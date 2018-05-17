package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cReset    = 0
	cBold     = 1
	cRed      = 31
	cGreen    = 32
	cYellow   = 33
	cBlue     = 34
	cMagenta  = 35
	cCyan     = 36
	cGray     = 37
	cDarkGray = 90
)

var consoleBufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 100))
	},
}

// ConsoleWriter reads a JSON object per write operation and output an
// optionally colored human readable version on the Out writer.
type ConsoleWriter struct {
	Out     io.Writer
	NoColor bool
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	p = decodeIfBinaryToBytes(p)
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&event)
	if err != nil {
		return
	}
	buf := consoleBufPool.Get().(*bytes.Buffer)
	defer consoleBufPool.Put(buf)
	lvlColor := cReset
	level := "????"
	if l, ok := event[LevelFieldName].(string); ok {
		if !w.NoColor {
			lvlColor = levelColor(l)
		}
		level = strings.ToUpper(l)[0:4]
	}
	fmt.Fprintf(buf, "%s |%s| %s",
		colorize(formatTime(event[TimestampFieldName]), cDarkGray, !w.NoColor),
		colorize(level, lvlColor, !w.NoColor),
		colorize(event[MessageFieldName], cReset, !w.NoColor))
	fields := make([]string, 0, len(event))
	for field := range event {
		switch field {
		case LevelFieldName, TimestampFieldName, MessageFieldName:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		fmt.Fprintf(buf, " %s=", colorize(field, cCyan, !w.NoColor))
		switch value := event[field].(type) {
		case string:
			if needsQuote(value) {
				buf.WriteString(strconv.Quote(value))
			} else {
				buf.WriteString(value)
			}
		case json.Number:
			fmt.Fprint(buf, value)
		default:
			b, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(buf, "[error: %v]", err)
			} else {
				fmt.Fprint(buf, string(b))
			}
		}
	}
	buf.WriteByte('\n')
	buf.WriteTo(w.Out)
	n = len(p)
	return
}

func formatTime(t interface{}) string {
	switch t := t.(type) {
	case string:
		return t
	case json.Number:
		u, _ := t.Int64()
		return time.Unix(u, 0).Format(time.RFC3339)
	}
	return "<nil>"
}

func colorize(s interface{}, color int, enabled bool) string {
	if !enabled {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

func levelColor(level string) int {
	switch level {
	case "debug":
		return cMagenta
	case "info":
		return cGreen
	case "warn":
		return cYellow
	case "error", "fatal", "panic":
		return cRed
	default:
		return cReset
	}
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}
