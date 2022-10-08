package goerrors

import (
	"path/filepath"
	"strconv"

	goerrors "github.com/go-errors/errors"
)

var (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
)

// MarshalStack implements go-errors stack trace marshaling.
//
// zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) interface{} {
	sterr, ok := err.(*goerrors.Error)
	if !ok {
		return nil
	}
	st := sterr.StackFrames()
	out := make([]map[string]string, 0, len(st))
	for _, frame := range st {
		out = append(out, map[string]string{
			StackSourceFileName:     filepath.Base(frame.File),
			StackSourceLineName:     strconv.Itoa(frame.LineNumber),
			StackSourceFunctionName: frame.Name,
		})
	}

	return out
}
