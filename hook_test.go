package zerolog

import (
	"bytes"
	"io/ioutil"
	"testing"
)

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *Event, level Level, msg string) {
	levelName := level.String()
	if level == NoLevel {
		levelName = "nolevel"
	}
	e.Str("level_name", levelName)
}

type SimpleHook struct{}

func (h SimpleHook) Run(e *Event, level Level, msg string) {
	e.Bool("has_level", level != NoLevel)
	e.Str("test", "logged")
}

type CopyHook struct{}

func (h CopyHook) Run(e *Event, level Level, msg string) {
	hasLevel := level != NoLevel
	e.Bool("copy_has_level", hasLevel)
	if hasLevel {
		e.Str("copy_level", level.String())
	}
	e.Str("copy_msg", msg)
}

type NopHook struct{}

func (h NopHook) Run(e *Event, level Level, msg string) {
}

var (
	levelNameHook LevelNameHook
	simpleHook    SimpleHook
	copyHook      CopyHook
	nopHook       NopHook
)

func TestHook(t *testing.T) {
	t.Run("Message", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Log().Msg("test message")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level_name":"nolevel","message":"test message"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("NoLevel", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Log().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level_name":"nolevel"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Print", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Print("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"debug","level_name":"debug"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Error", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Copy/1", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(copyHook)
		log.Log().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"copy_has_level":false,"copy_msg":""}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Copy/2", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(copyHook)
		log.Info().Msg("a message")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","copy_has_level":true,"copy_level":"info","copy_msg":"a message","message":"a message"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Multi", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).Hook(simpleHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Multi/Message", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).Hook(simpleHook)
		log.Error().Msg("a message")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error","has_level":true,"test":"logged","message":"a message"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/single/pre", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Hook(levelNameHook).Output(out)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/single/post", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Output(out).Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/multi/pre", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Hook(levelNameHook).Hook(simpleHook).Output(out)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/multi/post", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Output(out).Hook(levelNameHook).Hook(simpleHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/mixed", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Hook(levelNameHook).Output(out).Hook(simpleHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/single/pre", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).With().Str("with", "pre").Logger()
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"pre","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/single/post", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).With().Str("with", "post").Logger().Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"post","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/multi/pre", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).Hook(simpleHook).With().Str("with", "pre").Logger()
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"pre","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/multi/post", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).With().Str("with", "post").Logger().Hook(levelNameHook).Hook(simpleHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"post","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/mixed", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).With().Str("with", "mixed").Logger().Hook(simpleHook)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"mixed","level_name":"error","has_level":true,"test":"logged"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("None", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Error().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func BenchmarkHooks(b *testing.B) {
	logger := New(ioutil.Discard)
	b.ResetTimer()
	b.Run("Nop/Single", func(b *testing.B) {
		log := logger.Hook(nopHook)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Log().Msg("")
			}
		})
	})
	b.Run("Nop/Multi", func(b *testing.B) {
		log := logger.Hook(nopHook).Hook(nopHook)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Log().Msg("")
			}
		})
	})
	b.Run("Simple", func(b *testing.B) {
		log := logger.Hook(simpleHook)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Log().Msg("")
			}
		})
	})
}
