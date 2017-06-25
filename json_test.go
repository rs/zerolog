package zerolog

import (
	"testing"
)

func TestAppendJSONString(t *testing.T) {
	encodeStringTests := []struct {
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
		{"ascii", `"ascii"`},
		{"\"a", `"\"a"`},
		{"\x1fa", `"\u001fa"`},
		{"foo\"bar\"baz", `"foo\"bar\"baz"`},
		{"\x1ffoo\x1fbar\x1fbaz", `"\u001ffoo\u001fbar\u001fbaz"`},
		{"emoji \u2764\ufe0f!", `"emoji ❤️!"`},
	}

	for _, tt := range encodeStringTests {
		b := appendJSONString([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendJSONString(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

func BenchmarkAppendJSONString(b *testing.B) {
	tests := map[string]string{
		"NoEncoding":     `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingFirst":  `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingMiddle": `aaaaaaaaaaaaaaaaaaaaaaaaa"aaaaaaaaaaaaaaaaaaaaaaaa`,
		"EncodingLast":   `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`,
		"RuneFirst":      `❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
		"RuneMiddle":     `aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa`,
		"RuneLast":       `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️`,
	}
	for name, str := range tests {
		b.Run(name, func(b *testing.B) {
			buf := make([]byte, 0, 100)
			for i := 0; i < b.N; i++ {
				_ = appendJSONString(buf, str)
			}
		})
	}
}
