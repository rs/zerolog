// +build !binary_log

package pkgerrors_test

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"github.com/rs/zerolog/pkgerrors"
)

func Test_StackTrace(t *testing.T) {
	tests := []struct {
		name string
		s    pkgerrors.Stack
		want pkgerrors.StackTrace
	}{
		{
			name: "Nil-Stack",
			s:    pkgerrors.Stack(nil),
			want: pkgerrors.StackTrace([]pkgerrors.Frame{}),
		},
		{
			name: "Empty-Stack",
			s:    pkgerrors.Stack([]uintptr{}),
			want: pkgerrors.StackTrace([]pkgerrors.Frame{}),
		},
		{
			name: "Success",
			s:    pkgerrors.Stack([]uintptr{1, 2, 3}),
			want: pkgerrors.StackTrace([]pkgerrors.Frame{1, 2, 3}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.StackTrace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

type state struct {
	flag        bool
	inputOutput map[string]writeResult
}

type writeResult struct {
	n   int
	err error
}

func (s *state) Write(b []byte) (n int, err error) {
	for input, result := range s.inputOutput {
		pattern := regexp.MustCompile(input)
		if pattern.Match(b) {
			return result.n, result.err
		}
	}

	return -1, nil
}

func (s *state) Width() (wid int, ok bool) {
	return -1, false
}

func (s *state) Precision() (prec int, ok bool) {
	return -1, false
}

func (s *state) Flag(c int) bool {
	return s.flag
}

func TestFrame_Format(t *testing.T) {
	type args struct {
		s    fmt.State
		verb rune
	}

	const (
		verbSourceFile         = 's'
		verbSourceLine         = 'd'
		verbFuncName           = 'n'
		verbSourceFileWithLine = 'v'
	)

	var (
		pc      = make([]uintptr, 10)
		someErr = errors.New("some error")
	)

	runtime.Callers(2, pc)

	tests := []struct {
		name string
		f    pkgerrors.Frame
		args args
	}{
		{
			name: "Source-File-With-Flag-False",
			f:    pkgerrors.Frame(1),
			args: args{
				s:    &state{},
				verb: verbSourceFile,
			},
		},
		{
			name: "Source-File-Invalid-Frame-Number",
			f:    pkgerrors.Frame(0),
			args: args{
				s: &state{
					flag: true,
					inputOutput: map[string]writeResult{
						`^unknown$`: {
							n:   0,
							err: nil,
						},
						`^\s+$`: {
							n:   1,
							err: nil,
						},
					},
				},
				verb: verbSourceFile,
			},
		},
		{
			name: "Source-File-Known-File-Name-and-Path",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					flag: true,
					inputOutput: map[string]writeResult{
						`^testing\.tRunner$`: {
							n:   0,
							err: nil,
						},
						`^\s+$`: {
							n:   1,
							err: nil,
						},
						`^.*testing\.go$`: {
							n:   1,
							err: nil,
						},
					},
				},
				verb: verbSourceFile,
			},
		},
		{
			name: "Source-File-Error-At-Write",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					flag: true,
					inputOutput: map[string]writeResult{
						`^testing\.tRunner$`: {
							n:   0,
							err: someErr,
						},
						`^\s+$`: {
							n:   1,
							err: someErr,
						},
						`^.*testing\.go$`: {
							n:   2,
							err: someErr,
						},
					},
				},
				verb: verbSourceFile,
			},
		},
		{
			name: "Source-Line",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					inputOutput: map[string]writeResult{
						`^[0-9]+$`: {
							n:   1,
							err: nil,
						},
					},
				},
				verb: verbSourceLine,
			},
		},
		{
			name: "Source-Line-Error-At-Write",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					inputOutput: map[string]writeResult{
						`^[0-9]+$`: {
							n:   0,
							err: someErr,
						},
					},
				},
				verb: verbSourceLine,
			},
		},
		{
			name: "Source-Line-Invalid-Frame-Number",
			f:    pkgerrors.Frame(1),
			args: args{
				s: &state{
					inputOutput: map[string]writeResult{
						`^[0-9]+$`: {
							n:   1,
							err: nil,
						},
					},
				},
				verb: verbSourceLine,
			},
		},
		{
			name: "Func-Name",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					inputOutput: map[string]writeResult{
						`^tRunner$`: {
							n:   1,
							err: nil,
						},
					},
				},
				verb: verbFuncName,
			},
		},
		{
			name: "Func-Name-Error-At-Write",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					inputOutput: map[string]writeResult{
						`^tRunner$`: {
							n:   1,
							err: someErr,
						},
					},
				},
				verb: verbFuncName,
			},
		},
		{
			name: "Source-File-and-Line",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					flag: true,
					inputOutput: map[string]writeResult{
						`^testing\.tRunner$`: {
							n:   0,
							err: nil,
						},
						`^\s+$`: {
							n:   1,
							err: nil,
						},
						`^.*testing\.go$`: {
							n:   2,
							err: nil,
						},
						`^:$`: {
							n:   2,
							err: nil,
						},
						`^[0-9]+$`: {
							n:   3,
							err: nil,
						},
					},
				},
				verb: verbSourceFileWithLine,
			},
		},
		{
			name: "Source-File-and-Line-Error-At-Write",
			f:    pkgerrors.Frame(pc[0]),
			args: args{
				s: &state{
					flag: true,
					inputOutput: map[string]writeResult{
						`^testing\.tRunner$`: {
							n:   0,
							err: someErr,
						},
						`^\s+$`: {
							n:   1,
							err: someErr,
						},
						`^.*testing\.go$`: {
							n:   2,
							err: someErr,
						},
						`^:$`: {
							n:   2,
							err: someErr,
						},
						`^[0-9]+$`: {
							n:   3,
							err: someErr,
						},
					},
				},
				verb: verbSourceFileWithLine,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Format(tt.args.s, tt.args.verb)
		})
	}
}
