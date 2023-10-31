package json

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"unicode"
)

var enc = Encoder{}

func TestAppendBytes(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := enc.AppendBytes([]byte{}, []byte(tt.in))
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendBytes(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

func TestAppendHex(t *testing.T) {
	for _, tt := range encodeHexTests {
		b := enc.AppendHex([]byte{}, []byte{tt.in})
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendHex(%x) = %s, want %s", tt.in, got, want)
		}
	}
}

var base64Encodings = []struct {
	name string
	enc  *base64.Encoding
}{
	{"base64.StdEncoding", base64.StdEncoding},
	{"base64.RawStdEncoding", base64.RawStdEncoding},
	{"base64.URLEncoding", base64.URLEncoding},
	{"base64.RawURLEncoding", base64.RawURLEncoding},
}

func TestAppendBase64(t *testing.T) {
	random := make([]byte, 19)
	_, _ = rand.Read(random)
	tests := [][]byte{{}, {'\x00'}, {'\xff'}, random}
	for _, input := range tests {
		for _, tt := range base64Encodings {
			b := enc.AppendBase64(tt.enc, []byte{}, input)
			if got, want := string(b), "\""+tt.enc.EncodeToString(input)+"\""; got != want {
				t.Errorf("appendBase64(%s, %x) = %s, want %s", tt.name, input, got, want)
			}
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

	encStr := string(enc.AppendString([]byte{}, s))
	encBytes := string(enc.AppendBytes([]byte{}, []byte(s)))

	if encStr != encBytes {
		i := 0
		for i < len(encStr) && i < len(encBytes) && encStr[i] == encBytes[i] {
			i++
		}
		encStr = encStr[i:]
		encBytes = encBytes[i:]
		i = 0
		for i < len(encStr) && i < len(encBytes) && encStr[len(encStr)-i-1] == encBytes[len(encBytes)-i-1] {
			i++
		}
		encStr = encStr[:len(encStr)-i]
		encBytes = encBytes[:len(encBytes)-i]

		if len(encStr) > 20 {
			encStr = encStr[:20] + "..."
		}
		if len(encBytes) > 20 {
			encBytes = encBytes[:20] + "..."
		}

		t.Errorf("encodings differ at %#q vs %#q", encStr, encBytes)
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
				_ = enc.AppendBytes(buf, byt)
			}
		})
	}
}
