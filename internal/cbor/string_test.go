package cbor

import (
	"bytes"
	"testing"
)

var encodeStringTests = []struct {
	in  string
	out string
}{
	{"", "\x60"},
	{"\\", "\x61\x5c"},
	{"\x00", "\x61\x00"},
	{"\x01", "\x61\x01"},
	{"\x02", "\x61\x02"},
	{"\x03", "\x61\x03"},
	{"\x04", "\x61\x04"},
	{"*", "\x61*"},
	{"a", "\x61a"},
	{"IETF", "\x64IETF"},
	{"abcdefghijklmnopqrstuvwxyzABCD", "\x78\x1eabcdefghijklmnopqrstuvwxyzABCD"},
	{"<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->",
		"\x79\x01\x2c<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->"},
	{"emoji \u2764\ufe0f!", "\x6demoji ❤️!"},
}

var encodeByteTests = []struct {
	in  []byte
	out string
}{
	{[]byte{}, "\x40"},
	{[]byte("\\"), "\x41\x5c"},
	{[]byte("\x00"), "\x41\x00"},
	{[]byte("\x01"), "\x41\x01"},
	{[]byte("\x02"), "\x41\x02"},
	{[]byte("\x03"), "\x41\x03"},
	{[]byte("\x04"), "\x41\x04"},
	{[]byte("*"), "\x41*"},
	{[]byte("a"), "\x41a"},
	{[]byte("IETF"), "\x44IETF"},
	{[]byte("abcdefghijklmnopqrstuvwxyzABCD"), "\x58\x1eabcdefghijklmnopqrstuvwxyzABCD"},
	{[]byte("<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->"),
		"\x59\x01\x2c<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->"},
	{[]byte("emoji \u2764\ufe0f!"), "\x4demoji ❤️!"},
}

func TestAppendString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := AppendString([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
	//Test a large string > 65535 length

	var buffer bytes.Buffer
	for i := 0; i < 0x00011170; i++ { //70,000 character string
		buffer.WriteString("a")
	}
	inp := buffer.String()
	want := "\x7a\x00\x01\x11\x70" + inp
	b := AppendString([]byte{}, inp)
	if got := string(b); got != want {
		t.Errorf("appendString(%q) = %#q, want %#q", inp, got, want)
	}
}

func TestAppendBytes(t *testing.T) {
	for _, tt := range encodeByteTests {
		b := AppendBytes([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
	//Test a large string > 65535 length

	inp := []byte{}
	for i := 0; i < 0x00011170; i++ { //70,000 character string
		inp = append(inp, byte('a'))
	}
	want := "\x5a\x00\x01\x11\x70" + string(inp)
	b := AppendBytes([]byte{}, inp)
	if got := string(b); got != want {
		t.Errorf("appendString(%q) = %#q, want %#q", inp, got, want)
	}
}
func BenchmarkAppendString(b *testing.B) {
	tests := map[string]string{
		"NoEncoding":       `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingFirst":    `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingMiddle":   `aaaaaaaaaaaaaaaaaaaaaaaaa"aaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingLast":     `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`,
		"MultiBytesFirst":  `❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"MultiBytesMiddle": `aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa`,
		"MultiBytesLast":   `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️`,
	}
	for name, str := range tests {
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 110)
			for i := 0; i < b.N; i++ {
				_ = AppendString(buf, str)
			}
		})
	}
}
