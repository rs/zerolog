package zerolog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type loggableError struct {
	error
}

func (l loggableError) MarshalZerologObject(e *Event) {
	if l.error == nil {
		return
	}
	e.Str("l", strings.ToUpper(l.error.Error()))
}

type nonLoggableError struct {
	error
	line int
}

type wrappedError struct {
	error
	msg string
}

func (w wrappedError) Error() string {
	if w.error == nil {
		return w.msg
	}
	return w.error.Error() + ": " + w.msg
}

func TestArrayErrorMarshalFunc(t *testing.T) {
	prefixed := func(s, prefix string) string {
		if s == "null" {
			return ""
		}
		return prefix + s + `,`
	}
	errs := []error{
		nil,
		fmt.Errorf("failure"),
		loggableError{fmt.Errorf("whoops")},
		nonLoggableError{fmt.Errorf("oops"), 402},
	}
	type testCase struct {
		name    string
		marshal func(err error) interface{}
		want    []string
	}
	testCases := []testCase{
		{
			name:    "default",
			marshal: nil,
			want:    []string{`null`, `"failure"`, `{"l":"WHOOPS"}`, `"oops"`},
		},
		{
			name: "string",
			marshal: func(err error) interface{} {
				if err == nil {
					return nil
				}
				return err.Error()
			},
			want: []string{`null`, `"failure"`, `"whoops"`, `"oops"`},
		},
		{
			name: "loggable",
			marshal: func(err error) interface{} {
				if err == nil {
					return nil
				}
				return loggableError{err}
			},
			want: []string{`null`, `{"l":"FAILURE"}`, `{"l":"WHOOPS"}`, `{"l":"OOPS"}`},
		},
		{
			name: "non-loggable",
			marshal: func(err error) interface{} {
				if err == nil {
					return nil
				}
				return nonLoggableError{err, 404}
			},
			want: []string{`null`, `"failure"`, `"whoops"`, `"oops"`},
		},
		{
			name: "nil error",
			marshal: func(err error) interface{} {
				var errNil error = nil
				return errNil
			},
			want: []string{`null`, `null`, `null`, `null`},
		},
		{
			name: "wrapped error",
			marshal: func(err error) interface{} {
				if err == nil {
					return nil
				} else if we, ok := err.(wrappedError); ok {
					return we
				} else {
					return wrappedError{err, "addendum"}
				}
			},
			want: []string{`null`, `"failure: addendum"`, `"whoops: addendum"`, `"oops: addendum"`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalErrorMarshalFunc := ErrorMarshalFunc
			defer func() {
				ErrorMarshalFunc = originalErrorMarshalFunc
			}()

			if tc.marshal != nil {
				ErrorMarshalFunc = tc.marshal
			}

			t.Run("Err", func(t *testing.T) {
				for i, err := range errs {
					want := tc.want[i]
					t.Run("Arr", func(t *testing.T) {
						wants := `[` + want + `]`
						a := Arr().Err(err)
						if got := decodeObjectToStr(a.write([]byte{})); got != wants {
							t.Errorf("%s Array.Err(%v)\ngot:  %s\nwant: %s", tc.name, err, got, wants)
						}
					})
					t.Run("Ctx", func(t *testing.T) {
						wants := `{` + prefixed(want, `"error":`) + `"message":"msg"}` + "\n"
						out := &bytes.Buffer{}
						logger := New(out).With().Err(err).Logger()
						logger.Log().Msg("msg")
						if got := decodeIfBinaryToString(out.Bytes()); got != wants {
							t.Errorf("%s Ctx.Err(%v)\ngot:  %v\nwant: %v", tc.name, err, got, wants)
						}
					})
					t.Run("Event", func(t *testing.T) {
						wants := `{` + prefixed(want, `"error":`) + `"message":"msg"}` + "\n"
						out := &bytes.Buffer{}
						logger := New(out)
						logger.Log().Err(err).Msg("msg")
						if got := decodeIfBinaryToString(out.Bytes()); got != wants {
							t.Errorf("%s Event.Err(%v)\ngot:  %v\nwant: %v", tc.name, err, got, wants)
						}
					})
					t.Run("Fields", func(t *testing.T) {
						wants := `{"err":` + want + `,"message":"msg"}` + "\n"
						out := &bytes.Buffer{}
						logger := New(out)
						logger.Log().Fields(map[string]interface{}{"err": err}).Msg("msg")
						if got := decodeIfBinaryToString(out.Bytes()); got != wants {
							t.Errorf("%s Event.Fields(%v)\ngot:  %v\nwant: %v", tc.name, err, got, wants)
						}
					})
				}
			})

			t.Run("Errs", func(t *testing.T) {
				want := `[` + strings.Join(tc.want, ",") + `]`
				t.Run("Arr", func(t *testing.T) {
					a := Arr().Errs(errs)
					if got := decodeObjectToStr(a.write([]byte{})); got != want {
						t.Errorf("%s Array.Errs()\ngot:  %s\nwant: %s", tc.name, got, want)
					}
				})

				t.Run("Ctx", func(t *testing.T) {
					wants := `{"e":` + want + `,"message":"msg"}` + "\n"
					out := &bytes.Buffer{}
					logger := New(out).With().Errs("e", errs).Logger()
					logger.Log().Msg("msg")
					if got := decodeIfBinaryToString(out.Bytes()); got != wants {
						t.Errorf("%s Ctx.Errs()\ngot:  %v\nwant: %v", tc.name, got, wants)
					}
				})
				t.Run("Event", func(t *testing.T) {
					wants := `{"e":` + want + `,"message":"msg"}` + "\n"
					out := &bytes.Buffer{}
					logger := New(out)
					logger.Log().Errs("e", errs).Msg("msg")
					if got := decodeIfBinaryToString(out.Bytes()); got != wants {
						t.Errorf("%s Ctx.Errs()\ngot:  %v\nwant: %v", tc.name, got, wants)
					}
				})
				t.Run("Fields", func(t *testing.T) {
					wants := `{"e":` + want + `,"message":"msg"}` + "\n"
					out := &bytes.Buffer{}
					logger := New(out)
					logger.Log().Fields(map[string]interface{}{"e": errs}).Msg("msg")
					if got := decodeIfBinaryToString(out.Bytes()); got != wants {
						t.Errorf("%s Ctx.Errs()\ngot:  %v\nwant: %v", tc.name, got, wants)
					}
				})
			})
		})
	}
}
