package cbor

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/rs/zerolog/internal"
)

var encodeStringTests = []struct {
	plain  string
	binary string
	json   string //begin and end quotes are implied
}{
	{"", "\x60", ""},
	{"\\", "\x61\x5c", "\\\\"},
	{"\"", "\x61\x22", "\\\""},
	{"\b", "\x61\x08", "\\b"},
	{"\f", "\x61\x0c", "\\f"},
	{"\n", "\x61\x0a", "\\n"},
	{"\r", "\x61\x0d", "\\r"},
	{"\t", "\x61\x09", "\\t"},
	{"Hi\t", "\x63Hi\x09", "Hi\\t"},
	{"\x00", "\x61\x00", "\\u0000"},
	{"\x01", "\x61\x01", "\\u0001"},
	{"\x02", "\x61\x02", "\\u0002"},
	{"\x03", "\x61\x03", "\\u0003"},
	{"\x04", "\x61\x04", "\\u0004"},
	{"*", "\x61*", "*"},
	{"a", "\x61a", "a"},
	{"IETF", "\x64IETF", "IETF"},
	{"abcdefghijklmnopqrstuvwxyzABCD", "\x78\x1eabcdefghijklmnopqrstuvwxyzABCD", "abcdefghijklmnopqrstuvwxyzABCD"},
	{"<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->" +
		"<------------------------------------  This is a 100 character string ----------------------------->",
		"\x79\x01\x2c<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->",
		"<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->" +
			"<------------------------------------  This is a 100 character string ----------------------------->"},
	{"emoji \u2764\ufe0f!", "\x6demoji ❤️!", "emoji \u2764\ufe0f!"},
	{"invalid utf8 \xff", "\x6einvalid utf8 \xff", "invalid utf8 \\ufffd"},
}

var encodeByteTests = []struct {
	plain  []byte
	binary string
}{
	{[]byte{}, "\x40"},
	{[]byte("\\"), "\x41\x5c"},
	{[]byte("\x00"), "\x41\x00"},
	{[]byte("\x01"), "\x41\x01"},
	{[]byte("\x02"), "\x41\x02"},
	{[]byte("\x03"), "\x41\x03"},
	{[]byte("\x04"), "\x41\x04"},
	{[]byte("\f"), "\x41\x0C"},
	{[]byte("\n"), "\x41\x0A"},
	{[]byte("\r"), "\x41\x0D"},
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
		b := enc.AppendString([]byte{}, tt.plain)
		if got, want := string(b), tt.binary; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.plain, got, want)
		}
	}
	//Test a large string > 65535 length

	var buffer bytes.Buffer
	for i := 0; i < 0x00011170; i++ { //70,000 character string
		buffer.WriteString("a")
	}
	inp := buffer.String()
	want := "\x7a\x00\x01\x11\x70" + inp
	b := enc.AppendString([]byte{}, inp)
	if got := string(b); got != want {
		t.Errorf("appendString(%q) = %#q, want %#q", inp, got, want)
	}
}
func TestAppendStrings(t *testing.T) {
	array := []string{}
	for _, tt := range encodeStringTests {
		array = append(array, tt.plain)
	}
	want := make([]byte, 0)
	want = append(want, 0x95) // start array
	for _, tt := range encodeStringTests {
		want = append(want, []byte(tt.binary)...)
	}

	got := enc.AppendStrings([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendStrings(%v)\ngot:  0x%s\nwant: 0x%s",
			array,
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now empty array case
	array = make([]string, 0)
	want = make([]byte, 0)
	want = append(want, 0x80) // start an empty string array
	got = enc.AppendStrings([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendStrings(%v)\ngot:  0x%s\nwant: 0x%s",
			array, hex.EncodeToString(got),
			hex.EncodeToString(want))
	}

	// now large array case
	array = make([]string, 24)
	want = make([]byte, 0)
	want = append(want, 0x98) // start a large array
	want = append(want, 0x18) // of length 24
	for i := 0; i < len(array); i++ {
		array[i] = "test"
		want = append(want, []byte("\x64test")...)
	}
	got = enc.AppendStrings([]byte{}, array)
	if !bytes.Equal(got, want) {
		t.Errorf("AppendStrings(%v)\ngot:  %s\nwant: %s",
			array,
			hex.EncodeToString(got),
			hex.EncodeToString(want))
	}
}

func TestAppendStringer(t *testing.T) {
	oldJSONMarshalFunc := JSONMarshalFunc
	defer func() {
		JSONMarshalFunc = oldJSONMarshalFunc
	}()

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return internal.InterfaceMarshalFunc(v)
	}

	for _, tt := range internal.EncodeStringerTests {
		got := enc.AppendStringer([]byte{}, tt.In)
		want := []byte(tt.Binary)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendStrings(%v)\ngot:  %s\nwant: %s",
				tt.In,
				hex.EncodeToString(got),
				hex.EncodeToString(want))
		}
	}
}

func TestAppendStringers(t *testing.T) {
	for _, tt := range internal.EncodeStringersTests {
		want := make([]byte, 0)
		want = append(want, []byte(tt.Binary)...)

		got := enc.AppendStringers([]byte{}, tt.In)
		if !bytes.Equal(got, want) {
			t.Errorf("AppendStrings(%v)\ngot:  %s\nwant: %s",
				tt,
				hex.EncodeToString(got),
				hex.EncodeToString(want))
		}
	}
}

func TestAppendBytes(t *testing.T) {
	for _, tt := range encodeByteTests {
		b := enc.AppendBytes([]byte{}, tt.plain)
		if got, want := string(b), tt.binary; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.plain, got, want)
		}
	}
	//Test a large string > 65535 length

	inp := []byte{}
	for i := 0; i < 0x00011170; i++ { //70,000 character string
		inp = append(inp, byte('a'))
	}
	want := "\x5a\x00\x01\x11\x70" + string(inp)
	b := enc.AppendBytes([]byte{}, inp)
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
			buf := make([]byte, 0, 120)
			for i := 0; i < b.N; i++ {
				_ = enc.AppendString(buf, str)
			}
		})
	}
}

func TestAppendEmbeddedJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty JSON",
			input: []byte{},
			want:  "\xd9\x01\x06@", // tag 0xd9 + empty byte string
		},
		{
			name:  "small JSON",
			input: []byte(`{"key":"value"}`),
			want:  "\xd9\x01\x06O{\"key\":\"value\"}", // tag 0xd9 + byte string with content
		},
		{
			name:  "large JSON (>23 bytes)",
			input: []byte(`{"key":"this is a very long value that exceeds the 23 byte limit for direct encoding"}`),
			want:  "\xd9\x01\x06XV{\"key\":\"this is a very long value that exceeds the 23 byte limit for direct encoding\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AppendEmbeddedJSON([]byte{}, tt.input)
			if string(got) != tt.want {
				t.Errorf("AppendEmbeddedJSON() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestAppendEmbeddedCBOR(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty CBOR",
			input: []byte{},
			want:  "\xd8?@", // tag 0xd8 + empty byte string
		},
		{
			name:  "small CBOR",
			input: []byte{0x01, 0x02, 0x03},
			want:  "\xd8?C\x01\x02\x03", // tag 0xd8 + byte string with 3 bytes
		},
		{
			name:  "large CBOR (>23 bytes)",
			input: make([]byte, 30),                        // 30 bytes of zeros
			want:  "\xd8?X\x1e" + string(make([]byte, 30)), // tag 0xd8 + byte string with length prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AppendEmbeddedCBOR([]byte{}, tt.input)
			if string(got) != tt.want {
				t.Errorf("AppendEmbeddedCBOR() = %q, want %q", string(got), tt.want)
			}
		})
	}
}
