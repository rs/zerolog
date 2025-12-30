//go:build binary_log
// +build binary_log

package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	stdlog "log"
	"time"
)

func ExampleLogger_With() {
	dst := bytes.Buffer{}
	log := New(&dst).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Msg("hello world")
	fmt.Println(decodeIfBinaryToString(dst.Bytes()))

	// Output: {"level":"info","foo":"bar","message":"hello world"}
}

func ExampleLogger_Level() {
	dst := bytes.Buffer{}
	log := New(&dst).Level(WarnLevel)

	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"error","message":"kept message"}
}

func ExampleLogger_Sample() {
	dst := bytes.Buffer{}
	log := New(&dst).Sample(&BasicSampler{N: 2})

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","message":"message 1"}
	// {"level":"info","message":"message 3"}
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","level_name":"info","the_message":"hello world","message":"hello world"}
}

func ExampleLogger_Print() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Print("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Printf() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Printf("hello %s", "world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Trace() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Trace().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"trace","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Debug() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"debug","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Info() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"info","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Warn() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"warn","foo":"bar","message":"a warning message"}
}

func ExampleLogger_Error() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"level":"error","error":"some error","message":"error doing something"}
}

func ExampleLogger_WithLevel() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.WithLevel(InfoLevel).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","message":"hello world"}
}

func ExampleLogger_Log() {
	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz"}
}

func ExampleEvent_Dict() {
	dst := bytes.Buffer{}
	log := New(&dst)

	e := log.Log().
		Str("foo", "bar")

	e.Dict("dict", e.CreateDict().
		Str("bar", "baz").
		Int("n", 1),
	).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

// User implements LogObjectMarshaler
func (uu Users) MarshalZerologArray(a *Array) {
	for _, u := range uu {
		a.Object(u)
	}
}

func ExampleEvent_Array() {
	dst := bytes.Buffer{}
	log := New(&dst)

	e := log.Log().
		Str("foo", "bar")

	e.Array("array", e.CreateArray().
		Str("baz").
		Int(1),
	).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
}

func ExampleContext_Objects() {
	// In go, arrays are type invariant so even if you have a variable u of type []User array and User implements
	// the LogObjectMarshaler interface, you cannot pass that to func that takes an []LogObjectMarshaler array in the
	// Objects call. In 1.24+ it allows passing the variadic covariant slice (e.g. u...) but the unit test needs to
	// work in earlier versions so we'll declare the array as []LogObjectMarshaler here.
	u := []LogObjectMarshaler{User{"John", 35, time.Time{}}, User{"Bob", 55, time.Time{}}}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Objects("users", u).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

func ExampleEvent_EmbedObject() {
	price := Price{val: 6449, prec: 2, unit: "$"}

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		EmbedObject(price).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}

func ExampleEvent_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleEvent_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	dst := bytes.Buffer{}
	log := New(&dst)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_Dict() {
	dst := bytes.Buffer{}
	ctx := New(&dst).With().
		Str("foo", "bar")

	logger := ctx.Dict("dict", ctx.CreateDict().
		Str("bar", "baz").
		Int("n", 1),
	).Logger()

	logger.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

func ExampleContext_Array() {
	dst := bytes.Buffer{}
	ctx := New(&dst).With().
		Str("foo", "bar")

	logger := ctx.Array("array", ctx.CreateArray().
		Str("baz").
		Int(1),
	).Logger()

	logger.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func (p Price) MarshalZerologObject(e *Event) {
	denom := uint64(1)
	for i := 0; i < p.prec; i++ {
		denom *= 10
	}
	result := []byte(p.unit)
	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
	e.Str("price", string(result))
}

func ExampleContext_EmbedObject() {
	price := Price{val: 6449, prec: 2, unit: "$"}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		EmbedObject(price).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
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

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}

func ExampleContext_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Fields(fields).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		Fields(fields).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","bar":"baz","n":1,"message":"hello world"}
}

func ExampleContext_IPAddr() {
	ipV4 := net.IP{192, 168, 0, 1}
	ipV6 := net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		IPAddr("v4", ipV4).
		IPAddr("v6", ipV6).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","v4":"192.168.0.1","v6":"2001:db8:85a3::8a2e:370:7334","message":"hello world"}
}
func ExampleContext_IPPrefix() {
	pfxV4 := net.IPNet{IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}
	pfxV6 := net.IPNet{IP: net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x00}, Mask: net.CIDRMask(64, 128)}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		IPPrefix("v4", pfxV4).
		IPPrefix("v6", pfxV6).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","v4":"192.168.0.100/24","v6":"2001:db8:85a3::8a2e:370:7300/64","message":"hello world"}
}
func ExampleContext_MACAddr() {
	mac := net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}

	dst := bytes.Buffer{}
	log := New(&dst).With().
		Str("foo", "bar").
		MACAddr("mac", mac).
		Logger()

	log.Log().Msg("hello world")

	fmt.Println(decodeIfBinaryToString(dst.Bytes()))
	// Output: {"foo":"bar","mac":"12:34:56:78:90:ab","message":"hello world"}
}
