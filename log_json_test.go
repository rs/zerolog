//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"encoding/base64"
	"os"
	"testing"
)

func TestJSONBytesMarshalFunc(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(out)
	log.Log().Bytes("bytes", []byte{'a', 'b', 'c', 1, 2, 3, 0xff}).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bytes":"abc\u0001\u0002\u0003\ufffd","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()
	origBytesMarshalFunc := JSONBytesMarshalFunc
	defer func() {
		JSONBytesMarshalFunc = origBytesMarshalFunc
	}()

	JSONBytesMarshalFunc = JSONBytesMarshalBase64(base64.StdEncoding)
	log.Log().Bytes("bytes", []byte{'a', 'b', 'c', 1, 2, 3, 0xff}).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bytes":"YWJjAQID/w==","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	JSONBytesMarshalFunc = JSONBytesMarshalBase64(base64.RawURLEncoding)
	log.Log().Bytes("bytes", []byte{'a', 'b', 'c', 1, 2, 3, 0xff}).Msg("msg")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bytes":"YWJjAQID_w","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()
}

func ExampleJSONBytesMarshalBase64() {
	log := New(os.Stdout)
	JSONBytesMarshalFunc = JSONBytesMarshalBase64(base64.StdEncoding)
	log.Info().Bytes("bytes", []byte{'a', 'b', 'c', 1, 2, 3, 0xff}).Msg("msg")
	// Output: {"level":"info","bytes":"YWJjAQID/w==","message":"msg"}
}
