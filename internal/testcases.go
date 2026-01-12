package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"time"
)

var BooleanTestCases = []struct {
	Val    bool
	Binary string
	Json   string
}{
	{true, "\xf5", "true"},
	{false, "\xf4", "false"},
}

var BooleanArrayTestCases = []struct {
	Val    []bool
	Binary string
	Json   string
}{
	{[]bool{}, "\x9f\xff", "[]"},
	{[]bool{false}, "\x81\xf4", "[false]"},
	{[]bool{true, false, true}, "\x83\xf5\xf4\xf5", "[true,false,true]"},
	{[]bool{true, false, false, true, false, true}, "\x86\xf5\xf4\xf4\xf5\xf4\xf5", "[true,false,false,true,false,true]"},
}

var IntegerTestCases = []struct {
	Val    int
	Binary string
}{
	// Value included in the type.
	{0, "\x00"},
	{1, "\x01"},
	{2, "\x02"},
	{3, "\x03"},
	{8, "\x08"},
	{9, "\x09"},
	{10, "\x0a"},
	{22, "\x16"},
	{23, "\x17"},

	// Value in 1 byte.
	{24, "\x18\x18"},
	{25, "\x18\x19"},
	{26, "\x18\x1a"},
	{127, "\x18\x7f"},
	{254, "\x18\xfe"},
	{255, "\x18\xff"},

	// Value in 2 bytes.
	{256, "\x19\x01\x00"},
	{257, "\x19\x01\x01"},
	{1000, "\x19\x03\xe8"},
	{0xFFFF, "\x19\xff\xff"},

	// Value in 4 bytes.
	{0x10000, "\x1a\x00\x01\x00\x00"},
	{0x7FFFFFFE, "\x1a\x7f\xff\xff\xfe"},
	{1000000, "\x1a\x00\x0f\x42\x40"},

	// Negative number test cases.
	// Value included in the type.
	{-1, "\x20"},
	{-2, "\x21"},
	{-3, "\x22"},
	{-10, "\x29"},
	{-21, "\x34"},
	{-22, "\x35"},
	{-23, "\x36"},
	{-24, "\x37"},

	// Value in 1 byte.
	{-25, "\x38\x18"},
	{-26, "\x38\x19"},
	{-100, "\x38\x63"},
	{-128, "\x38\x7f"},
	{-254, "\x38\xfd"},
	{-255, "\x38\xfe"},
	{-256, "\x38\xff"},

	// Value in 2 bytes.
	{-257, "\x39\x01\x00"},
	{-258, "\x39\x01\x01"},
	{-1000, "\x39\x03\xe7"},

	// Value in 4 bytes.
	{-0x10001, "\x3a\x00\x01\x00\x00"},
	{-0x7FFFFFFE, "\x3a\x7f\xff\xff\xfd"},
	{-1000000, "\x3a\x00\x0f\x42\x3f"},

	//Constants
	{math.MaxInt8, "\x18\x7f"},
	{math.MinInt8, "\x38\x7f"},
	{math.MaxInt16, "\x19\x7f\xff"},
	{math.MinInt16, "\x39\x7f\xff"},
	{math.MaxInt32, "\x1a\x7f\xff\xff\xff"},
	{math.MinInt32, "\x3a\x7f\xff\xff\xff"},
	{math.MaxInt64, "\x1b\x7f\xff\xff\xff\xff\xff\xff\xff"},
	{math.MinInt64, "\x3b\x7f\xff\xff\xff\xff\xff\xff\xff"},
}

type UnsignedIntTestCase struct {
	Val       uint
	Binary    string
	Bigbinary string
}

var AdditionalUnsignedIntegerTestCases = []UnsignedIntTestCase{
	{0x7FFFFFFF, "\x18\xff", "\x1a\x7f\xff\xff\xff"},
	{0x80000000, "\x19\xff\xff", "\x1a\x80\x00\x00\x00"},
	{1000000, "\x1b\x80\x00\x00\x00\x00\x00\x00\x00", "\x1a\x00\x0f\x42\x40"},

	//Constants
	{math.MaxUint8, "\x18\xff", "\x18\xff"},
	{math.MaxUint16, "\x19\xff\xff", "\x19\xff\xff"},
	{math.MaxUint32, "\x1a\xff\xff\xff\xff", "\x1a\xff\xff\xff\xff"},
	{math.MaxUint64, "\x1b\xff\xff\xff\xff\xff\xff\xff\xff", "\x1b\xff\xff\xff\xff\xff\xff\xff\xff"},
}

func unsignedIntegerTestCases() []UnsignedIntTestCase {
	size := len(IntegerTestCases) + len(AdditionalUnsignedIntegerTestCases)
	cases := make([]UnsignedIntTestCase, 0, size)
	cases = append(cases, AdditionalUnsignedIntegerTestCases...)
	for _, itc := range IntegerTestCases {
		if itc.Val < 0 {
			continue
		}
		cases = append(cases, UnsignedIntTestCase{Val: uint(itc.Val), Binary: itc.Binary, Bigbinary: itc.Binary})
	}
	return cases
}

var UnsignedIntegerTestCases = unsignedIntegerTestCases()

var Float32TestCases = []struct {
	Val    float32
	Binary string
}{
	{0.0, "\xfa\x00\x00\x00\x00"},
	{-0.0, "\xfa\x00\x00\x00\x00"},
	{1.0, "\xfa\x3f\x80\x00\x00"},
	{1.5, "\xfa\x3f\xc0\x00\x00"},
	{65504.0, "\xfa\x47\x7f\xe0\x00"},
	{-4.0, "\xfa\xc0\x80\x00\x00"},
	{0.00006103515625, "\xfa\x38\x80\x00\x00"},
	{float32(math.Inf(0)), "\xfa\x7f\x80\x00\x00"},
	{float32(math.Inf(-1)), "\xfa\xff\x80\x00\x00"},
	{float32(math.NaN()), "\xfa\x7f\xc0\x00\x00"},
	{math.SmallestNonzeroFloat32, "\xfa\x00\x00\x00\x01"},
	{math.MaxFloat32, "\xfa\x7f\x7f\xff\xff"},
}

var Float64TestCases = []struct {
	Val    float64
	Binary string
}{
	{0.0, "\xfa\x00\x00\x00\x00"},
	{-0.0, "\xfa\x00\x00\x00\x00"},
	{1.0, "\xfa\x3f\x80\x00\x00"},
	{1.5, "\xfa\x3f\xc0\x00\x00"},
	{65504.0, "\xfa\x47\x7f\xe0\x00"},
	{-4.0, "\xfa\xc0\x80\x00\x00"},
	{0.00006103515625, "\xfa\x38\x80\x00\x00"},
	{math.Inf(0), "\xfa\x7f\x80\x00\x00\x00\x00\x00\x00"},
	{math.Inf(-1), "\xfa\xff\x80\x00\x00\x00\x00\x00\x00"},
	{math.NaN(), "\xfb\x7f\xf8\x00\x00\x00\x00\x00\x00"},
	{math.SmallestNonzeroFloat64, "\xfa\x00\x00\x00\x00\x00\x00\x00\x01"},
	{math.MaxFloat64, "\xfa\x7f\x7f\xff\xff"},
}

var IntegerArrayTestCases = []struct {
	Val    []int
	Binary string
	Json   string
}{
	{[]int{}, "\x9f\xff", "[]"},
	{[]int{32768}, "\x81\x19\x80\x00", "[32768]"},
	{[]int{-1, 0, 200, 20}, "\x84\x20\x00\x18\xc8\x14", "[-1,0,200,20]"},
	{[]int{-200, -10, 200, 400}, "\x84\x38\xc7\x29\x18\xc8\x19\x01\x90", "[-200,-10,200,400]"},
	{[]int{1, 2, 3}, "\x83\x01\x02\x03", "[1,2,3]"},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		"\x98\x19\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x18\x18\x19",
		"[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]"},
}

var IpAddrTestCases = []struct {
	Ipaddr net.IP
	Text   string
	Binary string
}{
	{net.IP{10, 0, 0, 1}, "\"10.0.0.1\"", "\xd9\x01\x04\x44\x0a\x00\x00\x01"},
	{net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x0, 0x0, 0x0, 0x0, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
		"\"2001:db8:85a3::8a2e:370:7334\"",
		"\xd9\x01\x04\x50\x20\x01\x0d\xb8\x85\xa3\x00\x00\x00\x00\x8a\x2e\x03\x70\x73\x34"},
}

var IPAddrArrayTestCases = []struct {
	Val    []net.IP
	Binary string
	Json   string
}{
	{[]net.IP{}, "\x9f\xff", "[]"},
	{[]net.IP{{127, 0, 0, 0}}, "\x81\xd9\x01\x04\x44\x7f\x00\x00\x00", "[127.0.0.0]"},
	{[]net.IP{{0, 0, 0, 0}, {192, 168, 0, 100}}, "\x82\xd9\x01\x04\x44\x00\x00\x00\x00\xd9\x01\x04\x44\xc0\xa8\x00\x64", "[0.0.0.0,192.168.0.100]"},
}

var IPPrefixTestCases = []struct {
	Pfx    net.IPNet
	Text   string // ASCII representation of pfx
	Binary string // CBOR representation of pfx
}{
	{net.IPNet{IP: net.IP{0, 0, 0, 0}, Mask: net.CIDRMask(0, 32)}, "\"0.0.0.0/0\"", "\xd9\x01\x05\xa1\x44\x00\x00\x00\x00\x00"},
	{net.IPNet{IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}, "\"192.168.0.100/24\"",
		"\xd9\x01\x05\xa1\x44\xc0\xa8\x00\x64\x18\x18"},
	{net.IPNet{IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, Mask: net.CIDRMask(128, 128)}, "\"::1/128\"",
		"\xd9\x01\x05\xa1\x50\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x18\x80"},
}

var IPPrefixArrayTestCases = []struct {
	Val    []net.IPNet
	Binary string
	Json   string
}{
	{[]net.IPNet{}, "\x9f\xff", "[]"},
	{[]net.IPNet{{IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(24, 32)}}, "\x81\xd9\x01\x05\xa1\x44\x7f\x00\x00\x00\x18\x18", "[127.0.0.0/24]"},
	{[]net.IPNet{{IP: net.IP{0, 0, 0, 0}, Mask: net.CIDRMask(0, 32)}, {IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}}, "\x82\xd9\x01\x05\xa1\x44\x00\x00\x00\x00\x00\xd9\x01\x05\xa1\x44\xc0\xa8\x00\x64\x18\x18", "[0.0.0.0/0,192.168.0.100/24]"},
}

var MacAddrTestCases = []struct {
	Macaddr net.HardwareAddr
	Text    string // ASCII representation of macaddr
	Binary  string // CBOR representation of macaddr
}{
	{net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}, "\"12:34:56:78:90:ab\"", "\xd9\x01\x04\x46\x12\x34\x56\x78\x90\xab"},
	{net.HardwareAddr{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3}, "\"20:01:0d:b8:85:a3\"", "\xd9\x01\x04\x46\x20\x01\x0d\xb8\x85\xa3"},
}

var EncodeHexTests = []struct {
	In  byte
	Out string
}{
	{0x00, `"00"`},
	{0x0f, `"0f"`},
	{0x10, `"10"`},
	{0xf0, `"f0"`},
	{0xff, `"ff"`},
}

var EncodeStringTests = []struct {
	In  string
	Out string
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

var EncodeStringsTests = []struct {
	In  []string
	Out string
}{
	{nil, `[]`},
	{[]string{}, `[]`},
	{[]string{"A"}, `["A"]`},
	{[]string{"A", "B"}, `["A","B"]`},
}

var EncodeStringerTests = []struct {
	In     fmt.Stringer
	Out    string
	Binary string
}{
	{nil, `null`, "\xf6"},
	{fmt.Stringer(nil), `null`, "\xf6"},
	{net.IPv4bcast, `"255.255.255.255"`, "\x6f\x32\x35\x35\x2e\x32\x35\x35\x2e\x32\x35\x35\x2e\x32\x35\x35"},
}

var EncodeStringersTests = []struct {
	In     []fmt.Stringer
	Out    string
	Binary string
}{
	{nil, `[]`, "\x9f\xff"},
	{[]fmt.Stringer{}, `[]`, "\x9f\xff"},
	{[]fmt.Stringer{net.IPv4bcast}, `["255.255.255.255"]`, "\x9f\x6f255.255.255.255\xff"},
	{[]fmt.Stringer{net.IPv4allsys, net.IPv4allrouter}, `["224.0.0.1","224.0.0.2"]`, "\x9f\x69224.0.0.1\x69224.0.0.2\xff"},
}

var TimeIntegerTestcases = []struct {
	Txt     string
	Binary  string
	RfcStr  string
	UnixInt int
}{
	{"2013-02-03T19:54:00-08:00", "\xc1\x1a\x51\x0f\x30\xd8", "2013-02-04T03:54:00Z", 1359950040},
	{"1950-02-03T19:54:00-08:00", "\xc1\x3a\x25\x71\x93\xa7", "1950-02-04T03:54:00Z", -628200360},
}

var TimeFloatTestcases = []struct {
	RfcStr  string
	Out     string
	UnixInt int
}{
	{"2006-01-02T15:04:05.999999-08:00", "\xc1\xfb\x41\xd0\xee\x6c\x59\x7f\xff\xfc", 1136243045},
	{"1956-01-02T15:04:05.999999-08:00", "\xc1\xfb\xc1\xba\x53\x81\x1a\x00\x00\x11", -441680155},
}

var DurTestcases = []struct {
	Duration   time.Duration
	FloatOut   string
	IntegerOut string
}{
	{1000, "\xfb\x3f\xf0\x00\x00\x00\x00\x00\x00", "\x01"},
	{2000, "\xfb\x40\x00\x00\x00\x00\x00\x00\x00", "\x02"},
	{200000, "\xfb\x40\x69\x00\x00\x00\x00\x00\x00", "\x18\xc8"},
}

// inline copy from globals.go of InterfaceMarshalFunc used in tests to avoid import cycle
func InterfaceMarshalFunc(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 {
		// Remove trailing \n which is added by Encode.
		return b[:len(b)-1], nil
	}
	return b, nil
}
