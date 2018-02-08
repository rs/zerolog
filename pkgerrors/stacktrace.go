package pkgerrors

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/internal/json"
)

var (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
)

// MarshalStack implements pkg/errors stack trace marshaling.
//
//   zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) []byte {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var st errors.StackTrace
	if err, ok := err.(stackTracer); ok {
		st = err.StackTrace()
	} else {
		return nil
	}
	return appendJSONStack(make([]byte, 0, 500), st)
}

func appendJSONStack(dst []byte, st errors.StackTrace) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	dst = append(dst, '[')
	for i, frame := range st {
		if i > 0 {
			dst = append(dst, ',')
		}

		dst = append(dst, '{')

		fmt.Fprintf(buf, "%s", frame)
		dst = json.AppendString(dst, StackSourceFileName)
		dst = append(dst, ':')
		dst = json.AppendBytes(dst, buf.Bytes())
		dst = append(dst, ',')
		buf.Reset()

		fmt.Fprintf(buf, "%d", frame)
		dst = json.AppendString(dst, StackSourceLineName)
		dst = append(dst, ':')
		dst = json.AppendBytes(dst, buf.Bytes())
		dst = append(dst, ',')
		buf.Reset()

		fmt.Fprintf(buf, "%n", frame)
		dst = json.AppendString(dst, StackSourceFunctionName)
		dst = append(dst, ':')
		dst = json.AppendBytes(dst, buf.Bytes())

		dst = append(dst, '}')
	}
	dst = append(dst, ']')
	return dst
}
