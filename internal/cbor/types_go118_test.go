//go:build go1.18
// +build go1.18

package cbor

import (
	"encoding/hex"
	"net/netip"
	"testing"
)

var netipAddrTestCases = []struct {
	addr   netip.Addr
	text   string // ASCII representation of ipaddr
	binary string // CBOR representation of ipaddr
}{
	{
		netip.MustParseAddr("10.0.0.2"),
		`"10.0.0.2"`,
		"\xd9\x01\x04\x44\x0a\x00\x00\x02",
	},
	{
		netip.MustParseAddr("2001:db8:85a3::8a2e:370:7334"),
		`"2001:db8:85a3::8a2e:370:7334"`,
		"\xd9\x01\x04\x50\x20\x01\x0d\xb8\x85\xa3\x00\x00\x00\x00\x8a\x2e\x03\x70\x73\x34",
	},
}

func TestAppendNetipAddr(t *testing.T) {
	for _, tc := range netipAddrTestCases {
		s := enc.AppendNetipAddr([]byte{}, tc.addr)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendNetipAddr(%s)=0x%s, want: 0x%s",
				tc.addr, hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}

var netipPrefixTestCases = []struct {
	pfx    netip.Prefix
	text   string
	binary string
}{
	{
		netip.MustParsePrefix("0.0.0.0/0"),
		`"0.0.0.0/0"`,
		"\xd9\x01\x05\xa1\x44\x00\x00\x00\x00\x00",
	},
	{
		netip.MustParsePrefix("192.168.0.100/24"),
		`"192.168.0.100/24"`,
		"\xd9\x01\x05\xa1\x44\xc0\xa8\x00\x64\x18\x18",
	},
}

func TestAppendNetipPrefix(t *testing.T) {
	for _, tc := range netipPrefixTestCases {
		s := enc.AppendNetipPrefix([]byte{}, tc.pfx)
		got := string(s)
		if got != tc.binary {
			t.Errorf("AppendNetipPrefix(%s)=0x%s, want: 0x%s",
				tc.pfx.String(), hex.EncodeToString(s),
				hex.EncodeToString([]byte(tc.binary)))
		}
	}
}
