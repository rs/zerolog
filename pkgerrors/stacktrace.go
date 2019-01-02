package pkgerrors

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
)

// MarshalStack implements pkg/errors stack trace marshaling.
//
//   zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var st errors.StackTrace
	if err, ok := err.(stackTracer); ok {
		st = err.StackTrace()
	} else {
		return nil
	}
	out := make([]map[string]string, 0, len(st))
	for _, frame := range st {
		out = append(out, map[string]string{
			StackSourceFileName:     fmt.Sprintf("%s", frame),
			StackSourceLineName:     fmt.Sprintf("%d", frame),
			StackSourceFunctionName: fmt.Sprintf("%n", frame),
		})
	}
	return out
}
