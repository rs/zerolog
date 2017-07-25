package zerolog

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"
)

var (
	errExample  = errors.New("fail")
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogEmpty(b *testing.B) {
	logger := New(ioutil.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Log().Msg("")
		}
	})
}

func BenchmarkDisabled(b *testing.B) {
	logger := New(ioutil.Discard).Level(Disabled)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	logger := New(ioutil.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkContextFields(b *testing.B) {
	logger := New(ioutil.Discard).With().
		Str("string", "four!").
		Time("time", time.Time{}).
		Int("int", 123).
		Float32("float", -2.203230293249593).
		Logger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkLogFields(b *testing.B) {
	logger := New(ioutil.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().
				Str("string", "four!").
				Time("time", time.Time{}).
				Int("int", 123).
				Float32("float", -2.203230293249593).
				Msg(fakeMessage)
		}
	})
}

func BenchmarkLogFieldType(b *testing.B) {
	types := map[string]func(e *Event) *Event{
		"Int": func(e *Event) *Event {
			return e.Int("int", 1)
		},
		"Float32": func(e *Event) *Event {
			return e.Float32("float", 1)
		},
		"Str": func(e *Event) *Event {
			return e.Str("str", "foo")
		},
		"Err": func(e *Event) *Event {
			return e.Err(errExample)
		},
		"Time": func(e *Event) *Event {
			return e.Time("time", time.Time{})
		},
		"Dur": func(e *Event) *Event {
			return e.Dur("dur", 1*time.Millisecond)
		},
	}
	logger := New(ioutil.Discard)
	b.ResetTimer()
	for name := range types {
		f := types[name]
		b.Run(name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					f(logger.Info()).Msg("")
				}
			})
		})
	}
}
