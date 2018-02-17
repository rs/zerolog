// +build enable_binary_log

package zerolog

import (
	"errors"
	stdlog "log"
	"time"
	"bytes"
	"fmt"
	"testing"
	"io/ioutil"
	"reflect"
	"runtime"
)

func ExampleBinaryNew() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Info().Msg("hello world")
	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_With() {
	dst := bytes.Buffer{}
	log := New(&dst).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Msg("hello world")
	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))

	// Output: {"level":"info","foo":"bar","message":"hello world"}
}

func ExampleLogger_Level() {
	dst := bytes.Buffer{}
	log := New(&dst).Level(WarnLevel)

	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"error","message":"kept message"}
}

func ExampleLogger_Sample() {
	dst := bytes.Buffer{}
	log := New(&dst).Sample(&BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","message":"message 2"}
	// {"level":"info","message":"message 4"}
}

type LevelNameHook1 struct{}

func (h LevelNameHook1) Run(e *Event, l Level, msg string) {
	if l != NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *Event, l Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook1
	var messageHook MessageHook = "The message"

	dst := bytes.Buffer{}
	log := New(&dst).Hook(levelNameHook).Hook(messageHook)

	log.Info().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","level_name":"info","the_message":"hello world","message":"hello world"}
}

func ExampleLogger_Print() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Print("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Printf() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Printf("hello %s", "world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Debug() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Info() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Warn() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"warn","foo":"bar","message":"a warning message"}
}

func ExampleLogger_Error() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"error","error":"some error","message":"error doing something"}
}

func ExampleLogger_WithLevel() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.WithLevel(InfoLevel).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_Write() {
	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","message":"hello world"}
}

func ExampleLogger_Log() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz"}
}

func ExampleEvent_Dict() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Dict("dict", Dict().
			Str("bar", "baz").
			Int("n", 1),
		).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Users []User

func (uu Users) MarshalZerologArray(a *Array) {
	for _, u := range uu {
		a.Object(u)
	}
}

func ExampleEvent_Array() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Array("array", Arr().
			Str("baz").
			Int(1),
		).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","array":["baz",1],"message":"hello world"}
}

func ExampleEvent_Array_object() {
	dst := bytes.Buffer{}
	log := New(&dst)

	// Users implements LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log.Log().
		Str("foo", "bar").
		Array("users", u).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

func ExampleEvent_Object() {
	dst := bytes.Buffer{}
	log := New(&dst)

	// User implements LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	log.Log().
		Str("foo", "bar").
		Object("user", u).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
}

func ExampleEvent_Interface() {
	dst := bytes.Buffer{}
	log := New(&dst)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleEvent_Dur() {
	d := time.Duration(10 * time.Second)

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}

func ExampleContext_Dict() {
	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Dict("dict", Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

func ExampleContext_Array() {
	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Array("array", Arr().
			Str("baz").
			Int(1),
		).Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","array":["baz",1],"message":"hello world"}
}

func ExampleContext_Array_object() {
	// Users implements LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Array("users", u).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

func ExampleContext_Object() {
	// User implements LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Object("user", u).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleContext_Dur() {
	d := time.Duration(10 * time.Second)

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}

func ExampleContext_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(DecodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}


func TestArray(t *testing.T) {
	a := Arr().
		Bool(true).
		Int(1).
		Int8(2).
		Int16(3).
		Int32(4).
		Int64(5).
		Uint(6).
		Uint8(7).
		Uint16(8).
		Uint32(9).
		Uint64(10).
		Float32(11).
		Float64(12).
		Str("a").
		Time(time.Time{}).
		Dur(0)
	want := `[true,1,2,3,4,5,6,7,8,9,10,11,12,"a","0001-01-01T00:00:00Z",0]`
	if got := DecodeIfBinaryToString(a.write([]byte{})); got != want {
		t.Errorf("Array.write()\ngot:  %s\nwant: %s", got, want)
	}
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *Event, level Level, msg string) {
	levelName := level.String()
	if level == NoLevel {
		levelName = "nolevel"
	}
	e.Str("level_name", levelName)
}

var levelNameHook LevelNameHook

func TestHook(t *testing.T) {
	t.Run("Message", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Log().Msg("test message")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level_name":"nolevel","message":"test message"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("NoLevel", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Log().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level_name":"nolevel"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Print", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Print("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"debug","level_name":"debug"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Error", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/single/pre", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Hook(levelNameHook).Output(out)
		log.Error().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("Output/single/post", func(t *testing.T) {
		ignored := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := New(ignored).Output(out).Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"error","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/single/pre", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Hook(levelNameHook).With().Str("with", "pre").Logger()
		log.Error().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"pre","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
	t.Run("With/single/post", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).With().Str("with", "post").Logger().Hook(levelNameHook)
		log.Error().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"error","with":"post","level_name":"error"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

type NopHook struct{}

func (h NopHook) Run(e *Event, level Level, msg string) {
}

var nopHook NopHook

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
}

func TestLog(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), "{}\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Str("foo", "bar").Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar"}`+"\n"; got != want {
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
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Str("foo", "bar").Msg("")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar"}`+"\n"; got != want {
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
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestWith(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := New(out).With().
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
		Time("time", time.Time{})
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+3)
	log := ctx.Caller().Logger()
	log.Log().Msg("")
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11,"float64":12,"time":"0001-01-01T00:00:00Z","caller":"`+caller+`"}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"bool":true,"bytes":"bar","dur":1000,"error":"some error","float32":11,"float64":12,"int":1,"int16":3,"int32":4,"int64":5,"int8":2,"nil":null,"string":"foo","time":"0001-01-01T00:00:00Z","uint":6,"uint16":8,"uint32":9,"uint64":10,"uint8":7}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"caller":"`+caller+`","string":"foo","bytes":"bar","error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11,"float64":12,"dur":1000,"time":"0001-01-01T00:00:00Z","diff":10000}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"string":[],"err":[],"bool":[],"int":[],"int8":[],"int16":[],"int32":[],"int64":[],"uint":[],"uint8":[],"uint16":[],"uint32":[],"uint64":[],"float32":[],"float64":[],"dur":[],"time":[]}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"string":["foo"],"err":["some error"],"bool":[true],"int":[1],"int8":[2],"int16":[3],"int32":[4],"int64":[5],"uint":[6],"uint8":[7],"uint16":[8],"uint32":[9],"uint64":[10],"float32":[11],"float64":[12],"dur":[1000],"time":["0001-01-01T00:00:00Z"]}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"string":["foo","bar"],"err":["some error",null],"bool":[true,false],"int":[1,0],"int8":[2,0],"int16":[3,0],"int32":[4,0],"int64":[5,0],"uint":[6,0],"uint8":[7,0],"uint16":[8,0],"uint32":[9,0],"uint64":[10,0],"float32":[11,0],"float64":[12,0],"dur":[1000,0],"time":["0001-01-01T00:00:00Z","0001-01-01T00:00:00Z"]}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestMsgf(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Msgf("one %s %.1f %d %v", "two", 3.4, 5, errors.New("six"))
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"message":"one two 3.4 5 six"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestWithAndFieldsCombined(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out).With().Str("f1", "val").Str("f2", "val").Logger()
	log.Log().Str("f3", "val").Msg("")
	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"f1":"val","f2":"val","f3":"val"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLevel(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Info().Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(Disabled)
		log.Log().Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Log().Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Panic", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(PanicLevel)
		log.Log().Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/WithLevel", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.WithLevel(NoLevel).Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out).Level(InfoLevel)
		log.Info().Msg("test")
		if got, want := DecodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"test"}`+"\n"; got != want {
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
	if got, want := DecodeIfBinaryToString(out.Bytes()), "{\"i\":2}\n{\"i\":4}\n"; got != want {
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

	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
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

	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestOutputWithoutTimestamp(t *testing.T) {
	ignoredOut := &bytes.Buffer{}
	out := &bytes.Buffer{}
	log := New(ignoredOut).Output(out).With().Str("foo", "bar").Logger()
	log.Log().Msg("hello world")

	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","message":"hello world"}`+"\n"; got != want {
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

	if got, want := DecodeIfBinaryToString(out.Bytes()), `{"foo":"bar","time":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}


