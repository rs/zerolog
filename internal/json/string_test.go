package json

import (
	"testing"

	"github.com/rs/zerolog/internal"
)

func TestAppendString(t *testing.T) {
	for _, tt := range internal.EncodeStringTests {
		b := enc.AppendString([]byte{}, tt.In)
		if got, want := string(b), tt.Out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.In, got, want)
		}
	}
}

func TestAppendStrings(t *testing.T) {
	for _, tt := range internal.EncodeStringsTests {
		b := enc.AppendStrings([]byte{}, tt.In)
		if got, want := string(b), tt.Out; got != want {
			t.Errorf("appendStrings(%q) = %#q, want %#q", tt.In, got, want)
		}
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
		b := enc.AppendStringer([]byte{}, tt.In)
		if got, want := string(b), tt.Out; got != want {
			t.Errorf("AppendStringer(%q)\ngot:  %#q, want: %#q", tt.In, got, want)
		}
	}
}

func TestAppendStringers(t *testing.T) {
	for _, tt := range internal.EncodeStringersTests {
		b := enc.AppendStringers([]byte{}, tt.In)
		if got, want := string(b), tt.Out; got != want {
			t.Errorf("appendStrings(%q) = %#q, want %#q", tt.In, got, want)
		}
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
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = enc.AppendString(buf, str)
			}
		})
	}
}
