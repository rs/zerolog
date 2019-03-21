package pkgerrors

import (
	"github.com/pkg/errors"
)

var (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return false
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// MarshalStack implements pkg/errors stack trace marshaling.
//
//   zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) interface{} {
	sterr, ok := err.(stackTracer)
	if !ok {
		return nil
	}
	st := sterr.StackTrace()
	s := &state{}
	out := make([]map[string]string, 0, len(st))
	for _, frame := range st {
		out = append(out, map[string]string{
			StackSourceFileName:     frameField(frame, s, 's'),
			StackSourceLineName:     frameField(frame, s, 'd'),
			StackSourceFunctionName: frameField(frame, s, 'n'),
		})
	}
	return out
}

type stackTrace struct {
	Frames []frame `json:"stacktrace"`
}

type frame struct {
	StackSourceFileName string `json:"source"`
	StackSourceLineName string `json:"line"`
	StackSourceFuncName string `json:"func"`
}

// MarshalMultiStack properly implements pkg/errors stack trace marshaling by unwrapping the error stack.
//
//   zerolog.ErrorStackMarshaler = MarshalMultiStack
func MarshalMultiStack(err error) interface{} {
	stackTraces := []stackTrace{}
	currentErr := err
	for currentErr != nil {
		stack, ok := currentErr.(stackTracer)
		if !ok {
			// Unwrap again because errors.Wrap actually adds two
			// layers of wrapping.
			currentErr = unwrapErr(currentErr)
			continue
		}
		st := stack.StackTrace()
		s := &state{}
		stackTrace := stackTrace{}
		for _, f := range st {
			frame := frame{
				StackSourceFileName: frameField(f, s, 's'),
				StackSourceLineName: frameField(f, s, 'd'),
				StackSourceFuncName: frameField(f, s, 'n'),
			}
			stackTrace.Frames = append(stackTrace.Frames, frame)
		}
		stackTraces = append(stackTraces, stackTrace)

		currentErr = unwrapErr(currentErr)
	}
	return stackTraces
}

type causer interface {
	Cause() error
}

func unwrapErr(err error) error {
	cause, ok := err.(causer)
	if !ok {
		return nil
	}
	return cause.Cause()
}
