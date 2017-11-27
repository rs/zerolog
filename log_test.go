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
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Str("foo", "bar").Msg("")
		if got, want := out.String(), `{"foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
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
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Msg("")
		if got, want := out.String(), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Str("foo", "bar").Msg("")
		if got, want := out.String(), `{"level":"info","foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
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
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
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
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMap(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Fields(map[string]interface{}{
		"nil":     nil,
		"string":  "foo",
		"bytes":   []byte("bar"),
		"error":   errors.New("some error"),
		"bool":    true,
		"int":     int(1),
		"int8":    int8(2),
		"int16":   int16(3),
		"int32":   int32(4),
		"int64":   int64(5),
		"uint":    uint(6),
		"uint8":   uint8(7),
		"uint16":  uint16(8),
		"uint32":  uint32(9),
		"uint64":  uint64(10),
		"float32": float32(11),
		"float64": float64(12),
		"dur":     1 * time.Second,
		"time":    time.Time{},
	}).Msg("")
	if got, want := out.String(), `{"bool":true,"bytes":"bar","dur":1000,"error":"some error","float32":11,"float64":12,"int":1,"int16":3,"int32":4,"int64":5,"int8":2,"nil":null,"string":"foo","time":"0001-01-01T00:00:00Z","uint":6,"uint16":8,"uint32":9,"uint64":10,"uint8":7}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFields(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	now := time.Now()
	log.Log().
		Str("string", "foo").
		Bytes("bytes", []byte("bar")).
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
	if got, want := out.String(), `{"string":"foo","bytes":"bar","error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11,"float64":12,"dur":1000,"time":"0001-01-01T00:00:00Z","diff":10000}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArrayEmpty(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().
		Strs("string", []string{}).
		Errs("err", []error{}).
		Bools("bool", []bool{}).
		Ints("int", []int{}).
		Ints8("int8", []int8{}).
		Ints16("int16", []int16{}).
		Ints32("int32", []int32{}).
		Ints64("int64", []int64{}).
		Uints("uint", []uint{}).
		Uints8("uint8", []uint8{}).
		Uints16("uint16", []uint16{}).
		Uints32("uint32", []uint32{}).
		Uints64("uint64", []uint64{}).
		Floats32("float32", []float32{}).
		Floats64("float64", []float64{}).
		Durs("dur", []time.Duration{}).
		Times("time", []time.Time{}).
		Msg("")
	if got, want := out.String(), `{"string":[],"err":[],"bool":[],"int":[],"int8":[],"int16":[],"int32":[],"int64":[],"uint":[],"uint8":[],"uint16":[],"uint32":[],"uint64":[],"float32":[],"float64":[],"dur":[],"time":[]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArraySingleElement(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().
		Strs("string", []string{"foo"}).
		Errs("err", []error{errors.New("some error")}).
		Bools("bool", []bool{true}).
		Ints("int", []int{1}).
		Ints8("int8", []int8{2}).
		Ints16("int16", []int16{3}).
		Ints32("int32", []int32{4}).
		Ints64("int64", []int64{5}).
		Uints("uint", []uint{6}).
		Uints8("uint8", []uint8{7}).
		Uints16("uint16", []uint16{8}).
		Uints32("uint32", []uint32{9}).
		Uints64("uint64", []uint64{10}).
		Floats32("float32", []float32{11}).
		Floats64("float64", []float64{12}).
		Durs("dur", []time.Duration{1 * time.Second}).
		Times("time", []time.Time{time.Time{}}).
		Msg("")
	if got, want := out.String(), `{"string":["foo"],"err":["some error"],"bool":[true],"int":[1],"int8":[2],"int16":[3],"int32":[4],"int64":[5],"uint":[6],"uint8":[7],"uint16":[8],"uint32":[9],"uint64":[10],"float32":[11],"float64":[12],"dur":[1000],"time":["0001-01-01T00:00:00Z"]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArrayMultipleElement(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().
		Strs("string", []string{"foo", "bar"}).
		Errs("err", []error{errors.New("some error"), nil}).
		Bools("bool", []bool{true, false}).
		Ints("int", []int{1, 0}).
		Ints8("int8", []int8{2, 0}).
		Ints16("int16", []int16{3, 0}).
		Ints32("int32", []int32{4, 0}).
		Ints64("int64", []int64{5, 0}).
		Uints("uint", []uint{6, 0}).
		Uints8("uint8", []uint8{7, 0}).
		Uints16("uint16", []uint16{8, 0}).
		Uints32("uint32", []uint32{9, 0}).
		Uints64("uint64", []uint64{10, 0}).
		Floats32("float32", []float32{11, 0}).
		Floats64("float64", []float64{12, 0}).
		Durs("dur", []time.Duration{1 * time.Second, 0}).
		Times("time", []time.Time{time.Time{}, time.Time{}}).
		Msg("")
	if got, want := out.String(), `{"string":["foo","bar"],"err":["some error",null],"bool":[true,false],"int":[1,0],"int8":[2,0],"int16":[3,0],"int32":[4,0],"int64":[5,0],"uint":[6,0],"uint8":[7,0],"uint16":[8,0],"uint32":[9,0],"uint64":[10,0],"float32":[11,0],"float64":[12,0],"dur":[1000,0],"time":["0001-01-01T00:00:00Z","0001-01-01T00:00:00Z"]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsDisabled(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).Level(InfoLevel)
	now := time.Now()
	log.Debug().
		Str("string", "foo").
		Bytes("bytes", []byte("bar")).
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
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestMsgf(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := out.String(), `{"message":"one two 3.4 5 six"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestWithAndFieldsCombined(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).With().Str("f1", "val").Str("f2", "val").Logger()
	log.Log().Str("f3", "val").Msg("")
	if got, want := out.String(), `{"f1":"val","f2":"val","f3":"val"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLevel(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Info().Msg("test")
		if got, want := out.String(), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Log().Msg("test")
		if got, want := out.String(), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Log().Msg("test")
		if got, want := out.String(), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Panic", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(PanicLevel)
		log.Log().Msg("test")
		if got, want := out.String(), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/WithLevel", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.WithLevel(NoLevel).Msg("test")
		if got, want := out.String(), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Info().Msg("test")
		if got, want := out.String(), `{"level":"info","message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSampling(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).Sample(&BasicSampler{N: 2})
	log.Log().Int("i", 1).Msg("")
	log.Log().Int("i", 2).Msg("")
	log.Log().Int("i", 3).Msg("")
	log.Log().Int("i", 4).Msg("")
	if got, want := out.String(), "{\"i\":2}\n{\"i\":4}\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
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
	log.Log().Msg("nolevel-1")
	log.WithLevel(DebugLevel).Msg("5")
	log.WithLevel(InfoLevel).Msg("6")
	log.WithLevel(WarnLevel).Msg("7")
	log.WithLevel(ErrorLevel).Msg("8")
	log.WithLevel(NoLevel).Msg("nolevel-2")

	want := []struct {
		l Level
		p string
	}{
		{DebugLevel, `{"level":"debug","message":"1"}` + "\n"},
		{InfoLevel, `{"level":"info","message":"2"}` + "\n"},
		{WarnLevel, `{"level":"warn","message":"3"}` + "\n"},
		{ErrorLevel, `{"level":"error","message":"4"}` + "\n"},
		{NoLevel, `{"message":"nolevel-1"}` + "\n"},
		{DebugLevel, `{"level":"debug","message":"5"}` + "\n"},
		{InfoLevel, `{"level":"info","message":"6"}` + "\n"},
		{WarnLevel, `{"level":"warn","message":"7"}` + "\n"},
		{ErrorLevel, `{"level":"error","message":"8"}` + "\n"},
		{NoLevel, `{"message":"nolevel-2"}` + "\n"},
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
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
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
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestOutputWithoutTimestamp(t *testing.T) {
	ignoredOut := &bytes.Buffer{}
	out := &bytes.Buffer{}
	log := New(ignoredOut).Output(out).With().Str("foo", "bar").Logger()
	log.Log().Msg("hello world")

	if got, want := out.String(), `{"foo":"bar","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestOutputWithTimestamp(t *testing.T) {
	TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	defer func() {
		TimestampFunc = time.Now
	}()
	ignoredOut := &bytes.Buffer{}
	out := &bytes.Buffer{}
	log := New(ignoredOut).Output(out).With().Timestamp().Str("foo", "bar").Logger()
	log.Log().Msg("hello world")

	if got, want := out.String(), `{"time":"2001-02-03T04:05:06Z","foo":"bar","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}
