//go:build go1.18
// +build go1.18

package zerolog

import (
	"fmt"
)

func AsLogObjectMarshalers[T LogObjectMarshaler](objs []T) []LogObjectMarshaler {
	s := make([]LogObjectMarshaler, len(objs))
	for i, v := range objs {
		s[i] = v
	}
	return s
}

func AsStringers[T fmt.Stringer](objs []T) []fmt.Stringer {
	s := make([]fmt.Stringer, len(objs))
	for i, v := range objs {
		s[i] = v
	}
	return s
}
