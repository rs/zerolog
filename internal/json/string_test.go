package json

import (
	"fmt"
	"net"
	"testing"
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

var encodeHexTests = []struct {
	in  byte
	out string
}{
	{0x00, `"00"`},
	{0x0f, `"0f"`},
	{0x10, `"10"`},
	{0xf0, `"f0"`},
	{0xff, `"ff"`},
}

func TestAppendString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b := enc.AppendString([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendString(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

var encodeStringsTests = []struct {
	in  []string
	out string
}{
	{[]string{}, `[]`},
	{[]string{"A"}, `["A"]`},
	{[]string{"A", "B"}, `["A","B"]`},
}

func TestAppendStrings(t *testing.T) {
	for _, tt := range encodeStringsTests {
		b := enc.AppendStrings([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendStrings(%q) = %#q, want %#q", tt.in, got, want)
		}
	}
}

var encodeStringersTests = []struct {
	in  []fmt.Stringer
	out string
}{
	{[]fmt.Stringer{}, `[]`},
	{[]fmt.Stringer{net.IPv4bcast}, `["255.255.255.255"]`},
	{[]fmt.Stringer{net.IPv4allsys, net.IPv4allrouter}, `["224.0.0.1","224.0.0.2"]`},
}

func TestAppendStringers(t *testing.T) {
	for _, tt := range encodeStringersTests {
		b := enc.AppendStringers([]byte{}, tt.in)
		if got, want := string(b), tt.out; got != want {
			t.Errorf("appendStrings(%q) = %#q, want %#q", tt.in, got, want)
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
