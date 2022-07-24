package zerolog

import (
	"bytes"
	"fmt"
	"testing"
)

func TestContext_LogLevel(t *testing.T) {
	levels := []Level{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
		NoLevel,
		Disabled,
	}

	for _, l := range levels {
		t.Run(l.String(), func(t *testing.T) {
			out := &bytes.Buffer{}
			log := New(out).With().Logger().Level(l)
			log.UpdateContext(func(c Context) Context {
				return c.LogLevel()
			})
			log.Log().Msg("test")

			if l == Disabled {
				if got := decodeIfBinaryToString(out.Bytes()); got != `` {
					t.Errorf("invalid log output:\ngot:  %v\nwant: ", got)
				}
				return
			}
			if got, want := decodeIfBinaryToString(out.Bytes()), fmt.Sprintf(`{"level":"%s","message":"test"}`+"\n", l); got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})
	}
}
