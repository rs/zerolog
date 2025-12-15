package cbor

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestAppendKey(t *testing.T) {
	//if len(dst) < 1 {
	//	dst = e.AppendBeginMarker(dst)
	//}
	//return e.AppendString(dst, key)

	want := make([]byte, 0)
	want = append(want, 0xbf) // start string
	want = append(want, 0x63) // length 3
	want = append(want, []byte("key")...)

	got := enc.AppendKey([]byte{}, "key")
	if !bytes.Equal(got, want) {
		t.Errorf("AppendKey(%v)\ngot:  0x%s\nwant: 0x%s",
			"key",
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}
