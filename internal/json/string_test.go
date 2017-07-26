package json

import (
	"testing"
	"unicode"
)

var encodeStringTests = []struct {
	in  string
	out string
}{
	{"", `""`},
	{"\\", `"\\"`},
	{"\x00", `"\u0000"`},
	{"\x01", `"\u0001"`},
	{"\x02", `"\u0002"`},
	{"\x03", `"\u0003"`},
	{"\x04", `"\u0004"`},
	{"\x05", `"\u0005"`},
	{"\x06", `"\u0006"`},
	{"\x07", `"\u0007"`},
	{"\x08", `"\b"`},
	{"\x09", `"\t"`},
	{"\x0a", `"\n"`},
	{"\x0b", `"\u000b"`},
	{"\x0c", `"\f"`},
	{"\x0d", `"\r"`},
	{"\x0e", `"\u000e"`},
	{"\x0f", `"\u000f"`},
	{"\x10", `"\u0010"`},
	{"\x11", `"\u0011"`},
	{"\x12", `"\u0012"`},
	{"\x13", `"\u0013"`},
	{"\x14", `"\u0014"`},
	{"\x15", `"\u0015"`},
	{"\x16", `"\u0016"`},
	{"\x17", `"\u0017"`},
	{"\x18", `"\u0018"`},
	{"\x19", `"\u0019"`},
	{"\x1a", `"\u001a"`},
	{"\x1b", `"\u001b"`},
	{"\x1c", `"\u001c"`},
	{"\x1d", `"\u001d"`},
	{"\x1e", `"\u001e"`},
	{"\x1f", `"\u001f"`},
	{"✭", `"✭"`},
	{"foo\xc2\x7fbar", `"foo\ufffd\u007fbar"`}, // invalid sequence
	{"ascii", `"ascii"`},
	{"\"a", `"\"a"`},
	{"\x1fa", `"\u001fa"`},
	{"foo\"bar\"baz", `"foo\"bar\"baz"`},
	{"\x1ffoo\x1fbar\x1fbaz", `"\u001ffoo\u001fbar\u001fbaz"`},
	{"emoji \u2764\ufe0f!", `"emoji ❤️!"`},
}

func TestappendString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := AppendString([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

func TestappendBytes(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := AppendBytes([]byte{}, []byte(tt.in))
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendBytes(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

func TestStringBytes(t *testing.T) {
	t.Parallel()
	// Test that encodeState.stringBytes and encodeState.string use the same encoding.
	var r []rune
	for i := '\u0000'; i <= unicode.MaxRune; i++ {
		r = append(r, i)
	}
	s := string(r) + "\xff\xff\xffhello" // some invalid UTF-8 too

	enc := string(AppendString([]byte{}, s))
	encBytes := string(AppendBytes([]byte{}, []byte(s)))

	if enc != encBytes {
		i := 0
		for i < len(enc) && i < len(encBytes) && enc[i] == encBytes[i] {
			i++
		}
		enc = enc[i:]
		encBytes = encBytes[i:]
		i = 0
		for i < len(enc) && i < len(encBytes) && enc[len(enc)-i-1] == encBytes[len(encBytes)-i-1] {
			i++
		}
		enc = enc[:len(enc)-i]
		encBytes = encBytes[:len(encBytes)-i]

		if len(enc) > 20 {
			enc = enc[:20] + "..."
		}
		if len(encBytes) > 20 {
			encBytes = encBytes[:20] + "..."
		}

		t.Errorf("encodings differ at %#q vs %#q", enc, encBytes)
	}
}

func BenchmarkappendString(b *testing.B) {
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
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = AppendString(buf, str)
			}
		})
	}
}

func BenchmarkappendBytes(b *testing.B) {
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
		byt := []byte(str)
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = AppendBytes(buf, byt)
			}
		})
	}
}
