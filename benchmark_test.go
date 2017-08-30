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

func BenchmarkContextAppend(b *testing.B) {
	logger := New(ioutil.Discard).With().
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

type obj struct {
	Pub  string
	Tag  string `json:"tag"`
	priv int
}

func (o obj) MarshalZerologObject(e *Event) {
	e.Str("Pub", o.Pub).
		Str("Tag", o.Tag).
		Int("priv", o.priv)
}

func BenchmarkLogFieldType(b *testing.B) {
	bools := []bool{true, false, true, false, true, false, true, false, true, false}
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	floats := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	strings := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	durations := []time.Duration{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	times := []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	interfaces := []struct {
		Pub  string
		Tag  string `json:"tag"`
		priv int
	}{
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
		{"a", "a", 0},
	}
	objects := []obj{
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
		obj{"a", "a", 0},
	}
	errs := []error{errors.New("a"), errors.New("b"), errors.New("c"), errors.New("d"), errors.New("e")}
	types := map[string]func(e *Event) *Event{
		"Bool": func(e *Event) *Event {
			return e.Bool("k", bools[0])
		},
		"Bools": func(e *Event) *Event {
			return e.Bools("k", bools)
		},
		"Int": func(e *Event) *Event {
			return e.Int("k", ints[0])
		},
		"Ints": func(e *Event) *Event {
			return e.Ints("k", ints)
		},
		"Float": func(e *Event) *Event {
			return e.Float64("k", floats[0])
		},
		"Floats": func(e *Event) *Event {
			return e.Floats64("k", floats)
		},
		"Str": func(e *Event) *Event {
			return e.Str("k", strings[0])
		},
		"Strs": func(e *Event) *Event {
			return e.Strs("k", strings)
		},
		"Err": func(e *Event) *Event {
			return e.Err(errs[0])
		},
		"Errs": func(e *Event) *Event {
			return e.Errs("k", errs)
		},
		"Time": func(e *Event) *Event {
			return e.Time("k", times[0])
		},
		"Times": func(e *Event) *Event {
			return e.Times("k", times)
		},
		"Dur": func(e *Event) *Event {
			return e.Dur("k", durations[0])
		},
		"Durs": func(e *Event) *Event {
			return e.Durs("k", durations)
		},
		"Interface": func(e *Event) *Event {
			return e.Interface("k", interfaces[0])
		},
		"Interfaces": func(e *Event) *Event {
			return e.Interface("k", interfaces)
		},
		"Interface(Object)": func(e *Event) *Event {
			return e.Interface("k", objects[0])
		},
		"Interface(Objects)": func(e *Event) *Event {
			return e.Interface("k", objects)
		},
		"Object": func(e *Event) *Event {
			return e.Object("k", objects[0])
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
