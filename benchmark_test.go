package zerolog

import (
	"errors"
	"io"
	"testing"
	"time"
)

var (
	errExample  = errors.New("fail")
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogEmpty(b *testing.B) {
	logger := New(io.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Log().Msg("")
		}
	})
}

func BenchmarkDisabled(b *testing.B) {
	logger := New(io.Discard).Level(Disabled)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	logger := New(io.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkContextFields(b *testing.B) {
	logger := New(io.Discard).With().
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

func BenchmarkContextAppend(b *testing.B) {
	logger := New(io.Discard).With().
		Str("foo", "bar").
		Logger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.With().Str("bar", "baz")
		}
	})
}

func BenchmarkLogFields(b *testing.B) {
	logger := New(io.Discard)
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

func BenchmarkLogArrayObject(b *testing.B) {
	obj1 := fixtureObj{"a", "b", 2}
	obj2 := fixtureObj{"c", "d", 3}
	obj3 := fixtureObj{"e", "f", 4}
	logger := New(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arr := Arr()
		arr.Object(&obj1)
		arr.Object(&obj2)
		arr.Object(&obj3)
		logger.Info().Array("objects", arr).Msg("test")
	}
}

func BenchmarkLogFieldType(b *testing.B) {
	fixtures := makeFieldFixtures()
	types := map[string]func(e *Event) *Event{
		"Any": func(e *Event) *Event {
			return e.Any("k", fixtures.Interfaces[0])
		},
		"Bool": func(e *Event) *Event {
			return e.Bool("k", fixtures.Bools[0])
		},
		"Bools": func(e *Event) *Event {
			return e.Bools("k", fixtures.Bools)
		},
		"Bytes": func(e *Event) *Event {
			return e.Bytes("k", fixtures.Bytes)
		},
		"Hex": func(e *Event) *Event {
			return e.Hex("k", fixtures.Bytes)
		},
		"Int": func(e *Event) *Event {
			return e.Int("k", fixtures.Ints[0])
		},
		"Ints": func(e *Event) *Event {
			return e.Ints("k", fixtures.Ints)
		},
		"Float32": func(e *Event) *Event {
			return e.Float32("k", fixtures.Floats32[0])
		},
		"Floats32": func(e *Event) *Event {
			return e.Floats32("k", fixtures.Floats32)
		},
		"Float64": func(e *Event) *Event {
			return e.Float64("k", fixtures.Floats64[0])
		},
		"Floats64": func(e *Event) *Event {
			return e.Floats64("k", fixtures.Floats64)
		},
		"Str": func(e *Event) *Event {
			return e.Str("k", fixtures.Strings[0])
		},
		"Strs": func(e *Event) *Event {
			return e.Strs("k", fixtures.Strings)
		},
		"Stringer": func(e *Event) *Event {
			return e.Stringer("k", fixtures.Stringers[0])
		},
		"Stringers": func(e *Event) *Event {
			return e.Stringers("k", fixtures.Stringers)
		},
		"Err": func(e *Event) *Event {
			return e.Err(fixtures.Errs[0])
		},
		"Errs": func(e *Event) *Event {
			return e.Errs("k", fixtures.Errs)
		},
		"Ctx": func(e *Event) *Event {
			return e.Ctx(fixtures.Ctx)
		},
		"Time": func(e *Event) *Event {
			return e.Time("k", fixtures.Times[0])
		},
		"Times": func(e *Event) *Event {
			return e.Times("k", fixtures.Times)
		},
		"Dur": func(e *Event) *Event {
			return e.Dur("k", fixtures.Durations[0])
		},
		"Durs": func(e *Event) *Event {
			return e.Durs("k", fixtures.Durations)
		},
		"Interface": func(e *Event) *Event {
			return e.Interface("k", fixtures.Interfaces[0])
		},
		"Interfaces": func(e *Event) *Event {
			return e.Interface("k", fixtures.Interfaces)
		},
		"Interface(Object)": func(e *Event) *Event {
			return e.Interface("k", fixtures.Objects[0])
		},
		"Interface(Objects)": func(e *Event) *Event {
			return e.Interface("k", fixtures.Objects)
		},
		"Object": func(e *Event) *Event {
			return e.Object("k", fixtures.Objects[0])
		},
		"Objects": func(e *Event) *Event {
			return e.Objects("k", fixtures.Objects)
		},
		"Timestamp": func(e *Event) *Event {
			return e.Timestamp()
		},
		"IPAddr": func(e *Event) *Event {
			return e.IPAddr("k", fixtures.IPAddrs[0])
		},
		"IPAddrs": func(e *Event) *Event {
			return e.IPAddrs("k", fixtures.IPAddrs)
		},
		"IPPrefix": func(e *Event) *Event {
			return e.IPPrefix("k", fixtures.IPPfxs[0])
		},
		"IPPrefixes": func(e *Event) *Event {
			return e.IPPrefixes("k", fixtures.IPPfxs)
		},
		"MACAddr": func(e *Event) *Event {
			return e.MACAddr("k", fixtures.MACAddr)
		},
		"Type": func(e *Event) *Event {
			return e.Type("k", fixtures.Type)
		},
	}

	logger := New(io.Discard)
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

func BenchmarkContextFieldType(b *testing.B) {
	oldFormat := TimeFieldFormat
	TimeFieldFormat = TimeFormatUnix
	defer func() { TimeFieldFormat = oldFormat }()

	fixtures := makeFieldFixtures()
	types := map[string]func(c Context) Context{
		"Any": func(c Context) Context {
			return c.Any("k", fixtures.Interfaces[0])
		},
		"Bool": func(c Context) Context {
			return c.Bool("k", fixtures.Bools[0])
		},
		"Bools": func(c Context) Context {
			return c.Bools("k", fixtures.Bools)
		},
		"Bytes": func(c Context) Context {
			return c.Bytes("k", fixtures.Bytes)
		},
		"Hex": func(c Context) Context {
			return c.Hex("k", fixtures.Bytes)
		},
		"Int": func(c Context) Context {
			return c.Int("k", fixtures.Ints[0])
		},
		"Ints": func(c Context) Context {
			return c.Ints("k", fixtures.Ints)
		},
		"Float32": func(c Context) Context {
			return c.Float32("k", fixtures.Floats32[0])
		},
		"Floats32": func(c Context) Context {
			return c.Floats32("k", fixtures.Floats32)
		},
		"Float64": func(c Context) Context {
			return c.Float64("k", fixtures.Floats64[0])
		},
		"Floats64": func(c Context) Context {
			return c.Floats64("k", fixtures.Floats64)
		},
		"Str": func(c Context) Context {
			return c.Str("k", fixtures.Strings[0])
		},
		"Strs": func(c Context) Context {
			return c.Strs("k", fixtures.Strings)
		},
		"Stringer": func(c Context) Context {
			return c.Stringer("k", fixtures.Stringers[0])
		},
		"Stringers": func(c Context) Context {
			return c.Stringers("k", fixtures.Stringers)
		},
		"Err": func(c Context) Context {
			return c.Err(fixtures.Errs[0])
		},
		"Errs": func(c Context) Context {
			return c.Errs("k", fixtures.Errs)
		},
		"Ctx": func(c Context) Context {
			return c.Ctx(fixtures.Ctx)
		},
		"Time": func(c Context) Context {
			return c.Time("k", fixtures.Times[0])
		},
		"Times": func(c Context) Context {
			return c.Times("k", fixtures.Times)
		},
		"Dur": func(c Context) Context {
			return c.Dur("k", fixtures.Durations[0])
		},
		"Durs": func(c Context) Context {
			return c.Durs("k", fixtures.Durations)
		},
		"Interface": func(c Context) Context {
			return c.Interface("k", fixtures.Interfaces[0])
		},
		"Interfaces": func(c Context) Context {
			return c.Interface("k", fixtures.Interfaces)
		},
		"Interface(Object)": func(c Context) Context {
			return c.Interface("k", fixtures.Objects[0])
		},
		"Interface(Objects)": func(c Context) Context {
			return c.Interface("k", fixtures.Objects)
		},
		"Object": func(c Context) Context {
			return c.Object("k", fixtures.Objects[0])
		},
		"Objects": func(c Context) Context {
			return c.Objects("k", fixtures.Objects)
		},
		"Timestamp": func(c Context) Context {
			return c.Timestamp()
		},
		"IPAddr": func(c Context) Context {
			return c.IPAddr("k", fixtures.IPAddrs[0])
		},
		"IPAddrs": func(c Context) Context {
			return c.IPAddrs("k", fixtures.IPAddrs)
		},
		"IPPrefix": func(c Context) Context {
			return c.IPPrefix("k", fixtures.IPPfxs[0])
		},
		"IPPrefixes": func(c Context) Context {
			return c.IPPrefixes("k", fixtures.IPPfxs)
		},
		"MACAddr": func(c Context) Context {
			return c.MACAddr("k", fixtures.MACAddr)
		},
		"Type": func(c Context) Context {
			return c.Type("k", fixtures.Type)
		},
	}

	logger := New(io.Discard)
	b.ResetTimer()
	for name := range types {
		f := types[name]
		b.Run(name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					l := f(logger.With()).Logger()
					l.Info().Msg("")
				}
			})
		})
	}
}
