package zerolog_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func ExampleConsoleWriter() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	log.Info().Str("foo", "bar").Msg("Hello World")
	// Output: <nil> INF Hello World foo=bar
}

func ExampleConsoleWriter_customFormatters() {
	out := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
	out.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%-6s|", i)) }
	out.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s:", i) }
	out.FormatFieldValue = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%s", i)) }
	log := zerolog.New(out)

	log.Info().Str("foo", "bar").Msg("Hello World")
	// Output: <nil> INFO  | Hello World foo:BAR
}

func ExampleNewConsoleWriter() {
	out := zerolog.NewConsoleWriter()
	out.NoColor = true // For testing purposes only
	log := zerolog.New(out)

	log.Debug().Str("foo", "bar").Msg("Hello World")
	// Output: <nil> DBG Hello World foo=bar
}

func ExampleNewConsoleWriter_customFormatters() {
	out := zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			// Customize time format
			w.TimeFormat = time.RFC822
			// Customize level formatting
			w.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("[%-5s]", i)) }
		},
	)
	out.NoColor = true // For testing purposes only

	log := zerolog.New(out)

	log.Info().Str("foo", "bar").Msg("Hello World")
	// Output: <nil> [INFO ] Hello World foo=bar
}

func TestConsoleLogger(t *testing.T) {
	t.Run("Numbers", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := zerolog.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})
		log.Info().
			Float64("float", 1.23).
			Uint64("small", 123).
			Uint64("big", 1152921504606846976).
			Msg("msg")
		if got, want := strings.TrimSpace(buf.String()), "<nil> INF msg big=1152921504606846976 float=1.23 small=123"; got != want {
			t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
		}
	})
}

func TestConsoleWriter(t *testing.T) {
	t.Run("Default field formatter", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true, PartsOrder: []string{"foo"}}

		_, err := w.Write([]byte(`{"foo" : "DEFAULT"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "DEFAULT foo=DEFAULT\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write colorized", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: false}

		_, err := w.Write([]byte(`{"level" : "warn", "message" : "Foobar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "\x1b[2m<nil>\x1b[0m \x1b[31mWRN\x1b[0m Foobar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true}

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		_, err := w.Write([]byte(`{"time" : "` + d + `", "level" : "debug", "message" : "Foobar", "foo" : "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "12:00AM DBG Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write error field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true}

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		evt := `{"time" : "` + d + `", "level" : "error", "message" : "Foobar", "aaa" : "bbb", "error" : "Error"}`
		// t.Log(evt)

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "12:00AM ERR Foobar error=Error aaa=bbb\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write caller field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Cannot get working directory: %s", err)
		}

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		evt := `{"time" : "` + d + `", "level" : "debug", "message" : "Foobar", "foo" : "bar", "caller" : "` + cwd + `/foo/bar.go"}`
		// t.Log(evt)

		_, err = w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "12:00AM DBG foo/bar.go > Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write JSON field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true}

		evt := `{"level" : "debug", "message" : "Foobar", "foo" : [1, 2, 3], "bar" : true}`
		// t.Log(evt)

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "<nil> DBG Foobar bar=true foo=[1,2,3]\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})
}

func TestConsoleWriterConfiguration(t *testing.T) {
	t.Run("Sets TimeFormat", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true, TimeFormat: time.RFC3339}

		d := time.Unix(0, 0).UTC().Format(time.RFC3339)
		evt := `{"time" : "` + d + `", "level" : "info", "message" : "Foobar"}`

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "1970-01-01T00:00:00Z INF Foobar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Sets PartsOrder", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zerolog.ConsoleWriter{Out: buf, NoColor: true, PartsOrder: []string{"message", "level"}}

		evt := `{"level" : "info", "message" : "Foobar"}`
		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "Foobar INF\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})
}

func BenchmarkConsoleWriter(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var msg = []byte(`{"level" : "info", "foo" : "bar", "message" : "HELLO", "time" : "1990-01-01"}`)

	w := zerolog.ConsoleWriter{Out: ioutil.Discard, NoColor: false}

	for i := 0; i < b.N; i++ {
		w.Write(msg)
	}
}
