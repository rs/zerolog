package zerolog

import (
	"bytes"
	"reflect"
	"testing"
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
		{WarnLevel, `{"level":"warning","message":"3"}` + "\n"},
		{ErrorLevel, `{"level":"error","message":"4"}` + "\n"},
	}
	if got := lw.ops; !reflect.DeepEqual(got, want) {
		t.Errorf("invalid ops:\ngot:\n%v\nwant:\n%v", got, want)
	}
}
