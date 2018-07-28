package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	colorBold = iota + 1
	colorFaint
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

var (
	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}

	consoleDefaultTimeFormat = time.Kitchen
	consoleDefaultFormatter  = func(i interface{}) string { return fmt.Sprintf("%s", i) }
	consoleDefaultPartsOrder = []string{
		TimestampFieldName,
		LevelFieldName,
		CallerFieldName,
		MessageFieldName,
	}

	consoleNoColor    = false
	consoleTimeFormat = consoleDefaultTimeFormat
	consoleFormatters map[string]Formatter
)

func init() {
	consoleFormatters = make(map[string]Formatter)
	initConsoleFormatters()
}

// ConsoleWriter parses the JSON input and writes it in an
// (optionally) colorized, human-friendly format to Out.
type ConsoleWriter struct {
	Out        io.Writer
	NoColor    bool
	TimeFormat string
	PartsOrder []string
}

// NewConsoleWriter creates and initializes a new ConsoleWriter.
//
// Customize the writer by passing a function as an argument:
//
//     zerolog.NewConsoleWriter(
//     	func(w *zerolog.ConsoleWriter) {
//     		// Call methods on w...
//     	},
//     )
//
// or with the SetFormatter() method.
//
func NewConsoleWriter(options ...func(w *ConsoleWriter)) ConsoleWriter {
	w := ConsoleWriter{Out: os.Stdout, TimeFormat: consoleDefaultTimeFormat, PartsOrder: consoleDefaultPartsOrder}

	for _, opt := range options {
		opt(&w)
	}

	return w
}

// The event type describes the JSON "event" object.
type event map[string]interface{}

// Formatter transforms the input into a formatted string.
type Formatter func(interface{}) string

// Formatter returns a formatter by id or the default formatter if none is found.
func (w ConsoleWriter) Formatter(id string) Formatter {
	if f, ok := consoleFormatters[id]; ok {
		return f
	}
	return consoleDefaultFormatter
}

// SetFormatter registers a formatter function f for part id.
func (w ConsoleWriter) SetFormatter(id string, f Formatter) {
	consoleFormatters[id] = f
}

// Reset re-initializes the ConsoleWriter defaults,
// eg. when you change the writer configuration.
func (w ConsoleWriter) Reset() ConsoleWriter {
	consoleFormatters = make(map[string]Formatter)
	initConsoleFormatters()
	return w
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	// NOTE: Be defensive in order to keep backwards compatibility
	if w.PartsOrder == nil {
		w.PartsOrder = consoleDefaultPartsOrder
	}
	if w.TimeFormat == "" && consoleTimeFormat != consoleDefaultTimeFormat {
		consoleTimeFormat = consoleDefaultTimeFormat
	}
	if w.TimeFormat != "" && consoleTimeFormat != w.TimeFormat {
		consoleTimeFormat = w.TimeFormat
	}
	if w.NoColor == false && consoleNoColor != false {
		consoleNoColor = false
	}
	if w.NoColor == true && consoleNoColor != w.NoColor {
		consoleNoColor = w.NoColor
	}

	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer consoleBufPool.Put(buf)

	var evt event
	p = decodeIfBinaryToBytes(p)
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	for _, p := range w.PartsOrder {
		w.writePart(buf, evt, p)
	}

	w.writeFields(evt, buf)

	buf.WriteByte('\n')
	buf.WriteTo(w.Out)
	return len(p), nil
}

// writeFields appends formatted key-value pairs to buf.
func (w ConsoleWriter) writeFields(evt event, buf *bytes.Buffer) {
	var fields = make([]string, 0, len(evt))
	for field := range evt {
		switch field {
		case LevelFieldName, TimestampFieldName, MessageFieldName, CallerFieldName:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	if len(fields) > 0 {
		buf.WriteByte(' ')
	}

	// Move the "error" field to front
	ei := sort.Search(len(fields), func(i int) bool { return fields[i] >= ErrorFieldName })
	if ei < len(fields) && fields[ei] == ErrorFieldName {
		fields[ei] = ""
		fields = append([]string{ErrorFieldName}, fields...)
		var xfields = make([]string, 0, len(fields))
		for _, field := range fields {
			if field == "" { // Skip empty fields
				continue
			}
			xfields = append(xfields, field)
		}
		fields = xfields
	}

	for i, field := range fields {
		var fn Formatter
		var fv Formatter
		if _, ok := consoleFormatters[field+"_field_name"]; ok {
			fn = w.Formatter(field + "_field_name")
			fv = w.Formatter(field + "_field_value")
		} else {
			fn = w.Formatter("field_name")
			fv = w.Formatter("field_value")
		}

		buf.WriteString(fn(field))

		switch fValue := evt[field].(type) {
		case string:
			if needsQuote(fValue) {
				buf.WriteString(fv(strconv.Quote(fValue)))
			} else {
				buf.WriteString(fv(fValue))
			}
		case json.Number:
			buf.WriteString(fv(fValue))
		default:
			b, err := json.Marshal(fValue)
			if err != nil {
				fmt.Fprintf(buf, colorize("[error: %v]", colorRed, w.NoColor), err)
			} else {
				fmt.Fprint(buf, fv(b))
			}
		}

		if i < len(fields)-1 { // Skip space for last field
			buf.WriteByte(' ')
		}
	}
}

// writePart appends a formatted part to buf.
func (w ConsoleWriter) writePart(buf *bytes.Buffer, evt event, p string) {
	var s = w.Formatter(p)(evt[p])
	if len(s) > 0 {
		buf.WriteString(s)
		if p != w.PartsOrder[len(w.PartsOrder)-1] { // Skip space for last part
			buf.WriteByte(' ')
		}
	}
}

// needsQuote returns true when the string s should be quoted in output.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

// colorize returns the string s wrapped in ANSI escape code, unless disabled is true.
func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// initConsoleFormatters sets the default formatters.
func initConsoleFormatters() {
	// Timestamp
	//
	consoleFormatters[TimestampFieldName] = func(i interface{}) string {
		t := "<nil>"
		if tt, ok := i.(string); ok {
			ts, err := time.Parse(time.RFC3339, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Format(consoleTimeFormat)
			}
		}
		return colorize(t, colorFaint, consoleNoColor)
	}

	// Level
	//
	consoleFormatters[LevelFieldName] = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = colorize("DBG", colorYellow, consoleNoColor)
			case "info":
				l = colorize("INF", colorGreen, consoleNoColor)
			case "warn":
				l = colorize("WRN", colorRed, consoleNoColor)
			case "error":
				l = colorize(colorize("ERR", colorRed, consoleNoColor), colorBold, consoleNoColor)
			case "fatal":
				l = colorize(colorize("FTL", colorRed, consoleNoColor), colorBold, consoleNoColor)
			case "panic":
				l = colorize(colorize("PNC", colorRed, consoleNoColor), colorBold, consoleNoColor)
			default:
				l = colorize("???", colorBold, consoleNoColor)
			}
		} else {
			l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
		}
		return l
	}

	// Caller
	//
	consoleFormatters[CallerFieldName] = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
			c = colorize(c, colorBold, consoleNoColor) + colorize(" >", colorFaint, consoleNoColor)
		}
		return c
	}

	// Message
	//
	consoleFormatters[MessageFieldName] = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	// Field name
	//
	consoleFormatters["field_name"] = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorFaint, consoleNoColor)
	}

	// Field value
	//
	consoleFormatters["field_value"] = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	// Errors
	//
	consoleFormatters["error_field_name"] = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorRed, consoleNoColor)
	}
	consoleFormatters["error_field_value"] = func(i interface{}) string {
		// return colorize(colorize(fmt.Sprintf("%s", i), colorRed, consoleNoColor), colorBold, consoleNoColor)
		return colorize(fmt.Sprintf("%s", i), colorRed, consoleNoColor)
	}
}
