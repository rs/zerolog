package json

import (
	"bytes"
	"testing"
)

func TestAppendKey(t *testing.T) {
	want := make([]byte, 0)
	want = append(want, []byte("{\"key\":")...)

	got := enc.AppendKey([]byte("{"), "key") // test with empty object
	if !bytes.Equal(got, want) {
		t.Errorf("AppendKey(%v)\ngot:  %s\nwant: %s",
			"key",
			string(got),
			string(want))
	}

	want = make([]byte, 0)
	want = append(want, []byte("},\"key\":")...) // test with non-empty object

	got = enc.AppendKey([]byte("}"), "key")
	if !bytes.Equal(got, want) {
		t.Errorf("AppendKey(%v)\ngot:  %s\nwant: %s",
			"key",
			string(got),
			string(want))
	}
}
