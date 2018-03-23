package json

import (
	"testing"
	"unicode"
)

func TestAppendBytes(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := AppendBytes([]byte{}, []byte(tt.in))
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendBytes(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

func TestAppendHex(t *testing.T) {
	for _, tt := range encodeHexTests {
		b := AppendHex([]byte{}, []byte{tt.in})
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendHex(%x) = %s, want %s", tt.in, got, want)
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

func BenchmarkAppendBytes(b *testing.B) {
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
