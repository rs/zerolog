package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), "{}\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Str("foo", "bar").Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar"}`+"\n"; got != want {
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
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Str("foo", "bar").Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar"}`+"\n"; got != want {
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
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestWith(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := New(out).With().
		Str("string", "foo").
		Bytes("bytes", []byte("bar")).
		Hex("hex", []byte{0x12, 0xef}).
		RawJSON("json", []byte(`{"some":"json"}`)).
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
		Float32("float32", 11.101).
		Float64("float64", 12.30303).
		Time("time", time.Time{})
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+3)
	log := ctx.Caller().Logger()
	log.Log().Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":"foo","bytes":"bar","hex":"12ef","json":{"some":"json"},"error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.101,"float64":12.30303,"time":"0001-01-01T00:00:00Z","caller":"`+caller+`"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	// Validate CallerWithSkipFrameCount.
	out.Reset()
	_, file, line, _ = runtime.Caller(0)
	caller = fmt.Sprintf("%s:%d", file, line+5)
	log = ctx.CallerWithSkipFrameCount(3).Logger()
	func() {
		log.Log().Msg("")
	}()
	// The above line is a little contrived, but the line above should be the line due
	// to the extra frame skip.
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":"foo","bytes":"bar","hex":"12ef","json":{"some":"json"},"error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.101,"float64":12.30303,"time":"0001-01-01T00:00:00Z","caller":"`+caller+`"}`+"\n"; got != want {
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
		"ipv6":    net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
		"dur":     1 * time.Second,
		"time":    time.Time{},
		"obj":     obj{"a", "b", 1},
	}).Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":true,"bytes":"bar","dur":1000,"error":"some error","float32":11,"float64":12,"int":1,"int16":3,"int32":4,"int64":5,"int8":2,"ipv6":"2001:db8:85a3::8a2e:370:7334","nil":null,"obj":{"Pub":"a","Tag":"b","priv":1},"string":"foo","time":"0001-01-01T00:00:00Z","uint":6,"uint16":8,"uint32":9,"uint64":10,"uint8":7}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMapPnt(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Fields(map[string]interface{}{
		"string":  new(string),
		"bool":    new(bool),
		"int":     new(int),
		"int8":    new(int8),
		"int16":   new(int16),
		"int32":   new(int32),
		"int64":   new(int64),
		"uint":    new(uint),
		"uint8":   new(uint8),
		"uint16":  new(uint16),
		"uint32":  new(uint32),
		"uint64":  new(uint64),
		"float32": new(float32),
		"float64": new(float64),
		"dur":     new(time.Duration),
		"time":    new(time.Time),
	}).Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":false,"dur":0,"float32":0,"float64":0,"int":0,"int16":0,"int32":0,"int64":0,"int8":0,"string":"","time":"0001-01-01T00:00:00Z","uint":0,"uint16":0,"uint32":0,"uint64":0,"uint8":0}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMapNilPnt(t *testing.T) {
	var (
		stringPnt  *string
		boolPnt    *bool
		intPnt     *int
		int8Pnt    *int8
		int16Pnt   *int16
		int32Pnt   *int32
		int64Pnt   *int64
		uintPnt    *uint
		uint8Pnt   *uint8
		uint16Pnt  *uint16
		uint32Pnt  *uint32
		uint64Pnt  *uint64
		float32Pnt *float32
		float64Pnt *float64
		durPnt     *time.Duration
		timePnt    *time.Time
	)
	out := &bytes.Buffer{}
	log := New(out)
	fields := map[string]interface{}{
		"string":  stringPnt,
		"bool":    boolPnt,
		"int":     intPnt,
		"int8":    int8Pnt,
		"int16":   int16Pnt,
		"int32":   int32Pnt,
		"int64":   int64Pnt,
		"uint":    uintPnt,
		"uint8":   uint8Pnt,
		"uint16":  uint16Pnt,
		"uint32":  uint32Pnt,
		"uint64":  uint64Pnt,
		"float32": float32Pnt,
		"float64": float64Pnt,
		"dur":     durPnt,
		"time":    timePnt,
	}

	log.Log().Fields(fields).Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":null,"dur":null,"float32":null,"float64":null,"int":null,"int16":null,"int32":null,"int64":null,"int8":null,"string":null,"time":null,"uint":null,"uint16":null,"uint32":null,"uint64":null,"uint8":null}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFields(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	now := time.Now()
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+3)
	log.Log().
		Caller().
		Str("string", "foo").
		Bytes("bytes", []byte("bar")).
		Hex("hex", []byte{0x12, 0xef}).
		RawJSON("json", []byte(`{"some":"json"}`)).
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
		IPAddr("IPv4", net.IP{192, 168, 0, 100}).
		IPAddr("IPv6", net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}).
		MACAddr("Mac", net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}).
		IPPrefix("Prefix", net.IPNet{IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}).
		Float32("float32", 11.1234).
		Float64("float64", 12.321321321).
		Dur("dur", 1*time.Second).
		Time("time", time.Time{}).
		TimeDiff("diff", now, now.Add(-10*time.Second)).
		Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"caller":"`+caller+`","string":"foo","bytes":"bar","hex":"12ef","json":{"some":"json"},"error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"IPv4":"192.168.0.100","IPv6":"2001:db8:85a3::8a2e:370:7334","Mac":"00:14:22:01:23:45","Prefix":"192.168.0.100/24","float32":11.1234,"float64":12.321321321,"dur":1000,"time":"0001-01-01T00:00:00Z","diff":10000}`+"\n"; got != want {
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
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":[],"err":[],"bool":[],"int":[],"int8":[],"int16":[],"int32":[],"int64":[],"uint":[],"uint8":[],"uint16":[],"uint32":[],"uint64":[],"float32":[],"float64":[],"dur":[],"time":[]}`+"\n"; got != want {
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
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":["foo"],"err":["some error"],"bool":[true],"int":[1],"int8":[2],"int16":[3],"int32":[4],"int64":[5],"uint":[6],"uint8":[7],"uint16":[8],"uint32":[9],"uint64":[10],"float32":[11],"float64":[12],"dur":[1000],"time":["0001-01-01T00:00:00Z"]}`+"\n"; got != want {
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
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":["foo","bar"],"err":["some error",null],"bool":[true,false],"int":[1,0],"int8":[2,0],"int16":[3,0],"int32":[4,0],"int64":[5,0],"uint":[6,0],"uint8":[7,0],"uint16":[8,0],"uint32":[9,0],"uint64":[10,0],"float32":[11,0],"float64":[12,0],"dur":[1000,0],"time":["0001-01-01T00:00:00Z","0001-01-01T00:00:00Z"]}`+"\n"; got != want {
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
		Hex("hex", []byte{0x12, 0xef}).
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
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestMsgf(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"one two 3.4 5 six"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestWithAndFieldsCombined(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).With().Str("f1", "val").Str("f2", "val").Logger()
	log.Log().Str("f3", "val").Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"f1":"val","f2":"val","f3":"val"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLevel(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Info().Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Log().Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Log().Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Panic", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(PanicLevel)
		log.Log().Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/WithLevel", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.WithLevel(NoLevel).Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Info().Msg("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestGetLevel(t *testing.T) {
	levels := []Level{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
		NoLevel,
		Disabled,
	}
	for _, level := range levels {
		if got, want := New(nil).Level(level).GetLevel(), level; got != want {
			t.Errorf("GetLevel() = %v, want: %v", got, want)
		}
	}
}

func TestSampling(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).Sample(&BasicSampler{N: 2})
	log.Log().Int("i", 1).Msg("")
	log.Log().Int("i", 2).Msg("")
	log.Log().Int("i", 3).Msg("")
	log.Log().Int("i", 4).Msg("")
	if got, want := decodeIfBinaryToString(out.Bytes()), "{\"i\":1}\n{\"i\":3}\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestDiscard(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Discard().Str("a", "b").Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	// Double call
	log.Log().Discard().Discard().Str("a", "b").Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
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
	p = decodeIfBinaryToBytes(p)
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

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
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

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestOutputWithoutTimestamp(t *testing.T) {
	ignoredOut := &bytes.Buffer{}
	out := &bytes.Buffer{}
	log := New(ignoredOut).Output(out).With().Str("foo", "bar").Logger()
	log.Log().Msg("hello world")

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","message":"hello world"}`+"\n"; got != want {
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

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

type loggableError struct {
	error
}

func (l loggableError) MarshalZerologObject(e *Event) {
	e.Str("message", l.error.Error()+": loggableError")
}

func TestErrorMarshalFunc(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)

	// test default behaviour
	log.Log().Err(errors.New("err")).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	log.Log().Err(loggableError{errors.New("err")}).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":{"message":"err: loggableError"},"message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	// test overriding the ErrorMarshalFunc
	originalErrorMarshalFunc := ErrorMarshalFunc
	defer func() {
		ErrorMarshalFunc = originalErrorMarshalFunc
	}()

	ErrorMarshalFunc = func(err error) interface{} {
		return err.Error() + ": marshaled string"
	}
	log.Log().Err(errors.New("err")).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err: marshaled string","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	out.Reset()
	ErrorMarshalFunc = func(err error) interface{} {
		return errors.New(err.Error() + ": new error")
	}
	log.Log().Err(errors.New("err")).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err: new error","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	out.Reset()
	ErrorMarshalFunc = func(err error) interface{} {
		return loggableError{err}
	}
	log.Log().Err(errors.New("err")).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":{"message":"err: loggableError"},"message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestCallerMarshalFunc(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)

	// test default behaviour this is really brittle due to the line numbers
	// actually mattering for validation
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+2)
	log.Log().Caller().Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"caller":"`+caller+`","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	// test custom behavior. In this case we'll take just the last directory
	origCallerMarshalFunc := CallerMarshalFunc
	defer func() { CallerMarshalFunc = origCallerMarshalFunc }()
	CallerMarshalFunc = func(file string, line int) string {
		parts := strings.Split(file, "/")
		if len(parts) > 1 {
			return strings.Join(parts[len(parts)-2:], "/") + ":" + strconv.Itoa(line)
		} else {
			return file + ":" + strconv.Itoa(line)
		}
	}
	_, file, line, _ = runtime.Caller(0)
	caller = CallerMarshalFunc(file, line+2)
	log.Log().Caller().Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"caller":"`+caller+`","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLevelFieldMarshalFunc(t *testing.T) {
	origLevelFieldMarshalFunc := LevelFieldMarshalFunc
	LevelFieldMarshalFunc = func(l Level) string {
		return strings.ToUpper(l.String())
	}
	defer func() {
		LevelFieldMarshalFunc = origLevelFieldMarshalFunc
	}()
	out := &bytes.Buffer{}
	log := New(out)

	log.Debug().Msg("test")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"DEBUG","message":"test"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	log.Info().Msg("test")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"INFO","message":"test"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	log.Warn().Msg("test")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"WARN","message":"test"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	log.Error().Msg("test")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"ERROR","message":"test"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()
}

type errWriter struct {
	error
}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, w.error
}

func TestErrorHandler(t *testing.T) {
	var got error
	want := errors.New("write error")
	ErrorHandler = func(err error) {
		got = err
	}
	log := New(errWriter{want})
	log.Log().Msg("test")
	if got != want {
		t.Errorf("ErrorHandler err = %#v, want %#v", got, want)
	}
}
