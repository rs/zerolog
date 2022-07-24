//go:build !binary_log
// +build !binary_log

package zerolog_test

import (
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func ExampleNew() {
	log := zerolog.New(os.Stdout)

	log.Info().LogLevel().Msg("hello world")
	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_With() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().LogLevel().Msg("hello world")

	// Output: {"foo":"bar","level":"info","message":"hello world"}
}

func ExampleLogger_Level() {
	log := zerolog.New(os.Stdout).Level(zerolog.WarnLevel)

	log.Info().LogLevel().Msg("filtered out message")
	log.Error().LogLevel().Msg("kept message")

	// Output: {"level":"error","message":"kept message"}
}

func ExampleLogger_Sample() {
	log := zerolog.New(os.Stdout).Sample(&zerolog.BasicSampler{N: 2})

	log.Info().LogLevel().Msg("message 1")
	log.Info().LogLevel().Msg("message 2")
	log.Info().LogLevel().Msg("message 3")
	log.Info().LogLevel().Msg("message 4")

	// Output: {"level":"info","message":"message 1"}
	// {"level":"info","message":"message 3"}
}

type LevelNameHook struct{}

func (h LevelNameHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	if l != zerolog.NoLevel {
		e.Str("level_name", l.String())
	} else {
		e.Str("level_name", "NoLevel")
	}
}

type MessageHook string

func (h MessageHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Str("the_message", msg)
}

func ExampleLogger_Hook() {
	var levelNameHook LevelNameHook
	var messageHook MessageHook = "The message"

	log := zerolog.New(os.Stdout).Hook(levelNameHook).Hook(messageHook)

	log.Info().
		LogLevel().
		Msg("hello world")

	// Output: {"level":"info","level_name":"info","the_message":"hello world","message":"hello world"}
}

func ExampleLogger_Print() {
	log := zerolog.New(os.Stdout).LogLevel()

	log.Print("hello world")

	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Printf() {
	log := zerolog.New(os.Stdout).LogLevel()

	log.Printf("hello %s", "world")

	// Output: {"level":"debug","message":"hello world"}
}

func ExampleLogger_Trace() {
	log := zerolog.New(os.Stdout)

	log.Trace().
		Str("foo", "bar").
		Int("n", 123).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","n":123,"level":"trace","message":"hello world"}
}

func ExampleLogger_Debug() {
	log := zerolog.New(os.Stdout)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","n":123,"level":"debug","message":"hello world"}
}

func ExampleLogger_Info() {
	log := zerolog.New(os.Stdout)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","n":123,"level":"info","message":"hello world"}
}

func ExampleLogger_Warn() {
	log := zerolog.New(os.Stdout)

	log.Warn().
		Str("foo", "bar").
		LogLevel().
		Msg("a warning message")

	// Output: {"foo":"bar","level":"warn","message":"a warning message"}
}

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(errors.New("some error")).
		LogLevel().
		Msg("error doing something")

	// Output: {"error":"some error","level":"error","message":"error doing something"}
}

func ExampleLogger_WithLevel() {
	log := zerolog.New(os.Stdout)

	log.WithLevel(zerolog.InfoLevel).
		LogLevel().
		Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_Write() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		LogLevel().
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	// Output: {"foo":"bar","level":"trace","message":"hello world"}
}

func ExampleLogger_Log() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		LogLevel().
		Msg("")

	// Output: {"foo":"bar","bar":"baz","level":""}
}

func ExampleEvent_Dict() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"level":"","message":"hello world"}
}

type User struct {
	Name    string
	Age     int
	Created time.Time
}

func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Str("name", u.Name).
		Int("age", u.Age).
		Time("created", u.Created)
}

type Price struct {
	val  uint64
	prec int
	unit string
}

func (p Price) MarshalZerologObject(e *zerolog.Event) {
	denom := uint64(1)
	for i := 0; i < p.prec; i++ {
		denom *= 10
	}
	result := []byte(p.unit)
	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
	e.Str("price", string(result))
}

type Users []User

func (uu Users) MarshalZerologArray(a *zerolog.Array) {
	for _, u := range uu {
		a.Object(u)
	}
}

func ExampleEvent_Array() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1).
			Dict(zerolog.Dict().
				Str("bar", "baz").
				Int("n", 1),
			),
		).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","array":["baz",1,{"bar":"baz","n":1}],"level":"","message":"hello world"}
}

func ExampleEvent_Array_object() {
	log := zerolog.New(os.Stdout)

	// Users implements zerolog.LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log.Log().
		Str("foo", "bar").
		Array("users", u).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"level":"","message":"hello world"}
}

func ExampleEvent_Object() {
	log := zerolog.New(os.Stdout)

	// User implements zerolog.LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	log.Log().
		Str("foo", "bar").
		Object("user", u).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"level":"","message":"hello world"}
}

func ExampleEvent_EmbedObject() {
	log := zerolog.New(os.Stdout)

	price := Price{val: 6449, prec: 2, unit: "$"}

	log.Log().
		Str("foo", "bar").
		EmbedObject(price).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","price":"$64.49","level":"","message":"hello world"}
}

func ExampleEvent_Interface() {
	log := zerolog.New(os.Stdout)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"level":"","message":"hello world"}
}

func ExampleEvent_Dur() {
	d := 10 * time.Second

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"level":"","message":"hello world"}
}

func ExampleEvent_Durs() {
	d := []time.Duration{
		10 * time.Second,
		20 * time.Second,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"level":"","message":"hello world"}
}

func ExampleEvent_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"level":"","message":"hello world"}
}

func ExampleEvent_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Fields(fields).
		LogLevel().
		Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"level":"","message":"hello world"}
}

func ExampleContext_Dict() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"level":"trace","message":"hello world"}
}

func ExampleContext_Array() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Array("array", zerolog.Arr().
			Str("baz").
			Int(1),
		).Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","array":["baz",1],"level":"trace","message":"hello world"}
}

func ExampleContext_Array_object() {
	// Users implements zerolog.LogArrayMarshaler
	u := Users{
		User{"John", 35, time.Time{}},
		User{"Bob", 55, time.Time{}},
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Array("users", u).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"level":"trace","message":"hello world"}
}

func ExampleContext_Object() {
	// User implements zerolog.LogObjectMarshaler
	u := User{"John", 35, time.Time{}}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Object("user", u).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"level":"trace","message":"hello world"}
}

func ExampleContext_EmbedObject() {

	price := Price{val: 6449, prec: 2, unit: "$"}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		EmbedObject(price).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","price":"$64.49","level":"trace","message":"hello world"}
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"level":"trace","message":"hello world"}
}

func ExampleContext_Dur() {
	d := 10 * time.Second

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"level":"trace","message":"hello world"}
}

func ExampleContext_Durs() {
	d := []time.Duration{
		10 * time.Second,
		20 * time.Second,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"level":"trace","message":"hello world"}
}

func ExampleContext_IPAddr() {
	hostIP := net.IP{192, 168, 0, 100}
	log := zerolog.New(os.Stdout).With().
		IPAddr("HostIP", hostIP).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"HostIP":"192.168.0.100","level":"trace","message":"hello world"}
}

func ExampleContext_IPPrefix() {
	route := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
	log := zerolog.New(os.Stdout).With().
		IPPrefix("Route", route).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"Route":"192.168.0.0/24","level":"trace","message":"hello world"}
}

func ExampleContext_MACAddr() {
	mac := net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
	log := zerolog.New(os.Stdout).With().
		MACAddr("hostMAC", mac).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"hostMAC":"00:14:22:01:23:45","level":"trace","message":"hello world"}
}

func ExampleContext_Fields_map() {
	fields := map[string]interface{}{
		"bar": "baz",
		"n":   1,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Fields(fields).
		Logger().LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"level":"trace","message":"hello world"}
}

func ExampleContext_Fields_slice() {
	fields := []interface{}{
		"bar", "baz",
		"n", 1,
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Fields(fields).
		Logger().
		LogLevel()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","bar":"baz","n":1,"level":"trace","message":"hello world"}
}
