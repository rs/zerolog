package zerolog

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Msg("")
		if got, want := out.String(), "{}\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Str("foo", "bar").Msg("")
		if got, want := out.String(), `{"foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})

	t.Run("two-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().
			Str("foo", "bar").
			Int("n", 123).
			Msg("")
		if got, want := out.String(), `{"foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Msg("")
		if got, want := out.String(), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Str("foo", "bar").Msg("")
		if got, want := out.String(), `{"level":"info","foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})

	t.Run("two-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().
			Str("foo", "bar").
			Int("n", 123).
			Msg("")
		if got, want := out.String(), `{"level":"info","foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})
}

func TestWith(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).With().
		Str("foo", "bar").
		AnErr("some_err", nil).
		Err(errors.New("some error")).
		Bool("bool", true).
		Int("int", 1).
		Int8("int8", 2).
		Int16("int16", 3).
		Int32("int32", 4).
		Int64("int64", 5).
		Uint("uint", 6).
		Uint8("uint8", 7).
		Uint16("uint16", 8).
		Uint32("uint32", 9).
		Uint64("uint64", 10).
		Float32("float32", 11).
		Float64("float64", 12).
		Time("time", time.Time{}).
		Logger()
	log.Log().Msg("")
	if got, want := out.String(), `{"foo":"bar","error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11,"float64":12,"time":"0001-01-01T00:00:00Z"}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestFields(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	now := time.Now()
	log.Log().
		Str("foo", "bar").
		AnErr("some_err", nil).
		Err(errors.New("some error")).
		Bool("bool", true).
		Int("int", 1).
		Int8("int8", 2).
		Int16("int16", 3).
		Int32("int32", 4).
		Int64("int64", 5).
		Uint("uint", 6).
		Uint8("uint8", 7).
		Uint16("uint16", 8).
		Uint32("uint32", 9).
		Uint64("uint64", 10).
		Float32("float32", 11).
		Float64("float64", 12).
		Dur("dur", 1*time.Second).
		Time("time", time.Time{}).
		TimeDiff("diff", now, now.Add(-10*time.Second)).
		Msg("")
	if got, want := out.String(), `{"foo":"bar","error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11,"float64":12,"dur":1000,"time":"0001-01-01T00:00:00Z","diff":10000}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestFieldsDisabled(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).Level(InfoLevel)
	now := time.Now()
	log.Debug().
		Str("foo", "bar").
		AnErr("some_err", nil).
		Err(errors.New("some error")).
		Bool("bool", true).
		Int("int", 1).
		Int8("int8", 2).
		Int16("int16", 3).
		Int32("int32", 4).
		Int64("int64", 5).
		Uint("uint", 6).
		Uint8("uint8", 7).
		Uint16("uint16", 8).
		Uint32("uint32", 9).
		Uint64("uint64", 10).
		Float32("float32", 11).
		Float64("float64", 12).
		Dur("dur", 1*time.Second).
		Time("time", time.Time{}).
		TimeDiff("diff", now, now.Add(-10*time.Second)).
		Msg("")
	if got, want := out.String(), ""; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestMsgf(t *testing.T) {
	out := &bytes.Buffer{}
	New(out).Log().Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := out.String(), `{"message":"one two 3.4 5 six"}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestWithAndFieldsCombined(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).With().Str("f1", "val").Str("f2", "val").Logger()
	log.Log().Str("f3", "val").Msg("")
	if got, want := out.String(), `{"f1":"val","f2":"val","f3":"val"}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestLevel(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Info().Msg("test")
		if got, want := out.String(), ""; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})

	t.Run("Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Info().Msg("test")
		if got, want := out.String(), `{"level":"info","message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output: got %q, want %q", got, want)
		}
	})
}

func TestSampling(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).Sample(2)
	log.Log().Int("i", 1).Msg("")
	log.Log().Int("i", 2).Msg("")
	log.Log().Int("i", 3).Msg("")
	log.Log().Int("i", 4).Msg("")
	if got, want := out.String(), "{\"sample\":2,\"i\":2}\n{\"sample\":2,\"i\":4}\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

type levelWriter struct {
	ops []struct {
		l Level
		p string
	}
}

func (lw *levelWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (lw *levelWriter) WriteLevel(lvl Level, p []byte) (int, error) {
	lw.ops = append(lw.ops, struct {
		l Level
		p string
	}{lvl, string(p)})
	return len(p), nil
}

func TestLevelWriter(t *testing.T) {
	lw := &levelWriter{
		ops: []struct {
			l Level
			p string
		}{},
	}
	log := New(lw)
	log.Debug().Msg("1")
	log.Info().Msg("2")
	log.Warn().Msg("3")
	log.Error().Msg("4")
	want := []struct {
		l Level
		p string
	}{
		{DebugLevel, `{"level":"debug","message":"1"}` + "\n"},
		{InfoLevel, `{"level":"info","message":"2"}` + "\n"},
		{WarnLevel, `{"level":"warn","message":"3"}` + "\n"},
		{ErrorLevel, `{"level":"error","message":"4"}` + "\n"},
	}
	if got := lw.ops; !reflect.DeepEqual(got, want) {
		t.Errorf("invalid ops:\ngot:\n%v\nwant:\n%v", got, want)
	}
}

func TestContextTimestamp(t *testing.T) {
	TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	defer func() {
		TimestampFunc = time.Now
	}()
	out := &bytes.Buffer{}
	log := New(out).With().Timestamp().Str("foo", "bar").Logger()
	log.Log().Msg("hello world")

	if got, want := out.String(), `{"time":"2001-02-03T04:05:06Z","foo":"bar","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}

func TestEventTimestamp(t *testing.T) {
	TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	defer func() {
		TimestampFunc = time.Now
	}()
	out := &bytes.Buffer{}
	log := New(out).With().Str("foo", "bar").Logger()
	log.Log().Timestamp().Msg("hello world")

	if got, want := out.String(), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output: got %q, want %q", got, want)
	}
}
