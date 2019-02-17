package zerolog

import (
	"strconv"
	"time"
)
import "sync/atomic"

var (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

	// MessageFieldName is the field name used for the message field.
	MessageFieldName = "message"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName = "error"

	// CallerFieldName is the field name used for caller field.
	CallerFieldName = "caller"

	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 2

	// CallerMarshalFunc allows customization of global caller marshaling
	CallerMarshalFunc = func(file string, line int) string {
		return file+":"+strconv.Itoa(line)
	}

	// ErrorStackFieldName is the field name used for error stacks.
	ErrorStackFieldName = "stack"

	// ErrorStackMarshaler extract the stack from err if any.
	ErrorStackMarshaler func(err error) interface{}

	// ErrorMarshalFunc allows customization of global error marshaling
	ErrorMarshalFunc = func(err error) interface{} {
		return err
	}

	// TimeFieldFormat defines the time format of the Time field type.
	// If set to an empty string, the time is formatted as an UNIX timestamp
	// as integer.
	TimeFieldFormat = time.RFC3339

	// TimestampFunc defines the function called to generate a timestamp.
	TimestampFunc = time.Now

	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Dur method.
	DurationFieldUnit = time.Millisecond

	// DurationFieldInteger renders Dur fields as integer instead of float if
	// set to true.
	DurationFieldInteger = false

	// ErrorHandler is called whenever zerolog fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)
)

var (
	gLevel          = new(uint32)
	disableSampling = new(uint32)
)

// SetGlobalLevel sets the global override for log level. If this
// values is raised, all Loggers will use at least this value.
//
// To globally disable logs, set GlobalLevel to Disabled.
func SetGlobalLevel(l Level) {
	atomic.StoreUint32(gLevel, uint32(l))
}

// GlobalLevel returns the current global log level
func GlobalLevel() Level {
	return Level(atomic.LoadUint32(gLevel))
}

// DisableSampling will disable sampling in all Loggers if true.
func DisableSampling(v bool) {
	var i uint32
	if v {
		i = 1
	}
	atomic.StoreUint32(disableSampling, i)
}

func samplingDisabled() bool {
	return atomic.LoadUint32(disableSampling) == 1
}
