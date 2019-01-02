package json

import (
	"math"
	"net"
	"reflect"
	"testing"
)

func TestAppendType(t *testing.T) {
	w := map[string]func(interface{}) []byte{
		"AppendInt":     func(v interface{}) []byte { return enc.AppendInt([]byte{}, v.(int)) },
		"AppendInt8":    func(v interface{}) []byte { return enc.AppendInt8([]byte{}, v.(int8)) },
		"AppendInt16":   func(v interface{}) []byte { return enc.AppendInt16([]byte{}, v.(int16)) },
		"AppendInt32":   func(v interface{}) []byte { return enc.AppendInt32([]byte{}, v.(int32)) },
		"AppendInt64":   func(v interface{}) []byte { return enc.AppendInt64([]byte{}, v.(int64)) },
		"AppendUint":    func(v interface{}) []byte { return enc.AppendUint([]byte{}, v.(uint)) },
		"AppendUint8":   func(v interface{}) []byte { return enc.AppendUint8([]byte{}, v.(uint8)) },
		"AppendUint16":  func(v interface{}) []byte { return enc.AppendUint16([]byte{}, v.(uint16)) },
		"AppendUint32":  func(v interface{}) []byte { return enc.AppendUint32([]byte{}, v.(uint32)) },
		"AppendUint64":  func(v interface{}) []byte { return enc.AppendUint64([]byte{}, v.(uint64)) },
		"AppendFloat32": func(v interface{}) []byte { return enc.AppendFloat32([]byte{}, v.(float32)) },
		"AppendFloat64": func(v interface{}) []byte { return enc.AppendFloat64([]byte{}, v.(float64)) },
	}
	tests := []struct {
		name  string
		fn    string
		input interface{}
		want  []byte
	}{
		{"AppendInt8(math.MaxInt8)", "AppendInt8", int8(math.MaxInt8), []byte("127")},
		{"AppendInt16(math.MaxInt16)", "AppendInt16", int16(math.MaxInt16), []byte("32767")},
		{"AppendInt32(math.MaxInt32)", "AppendInt32", int32(math.MaxInt32), []byte("2147483647")},
		{"AppendInt64(math.MaxInt64)", "AppendInt64", int64(math.MaxInt64), []byte("9223372036854775807")},

		{"AppendUint8(math.MaxUint8)", "AppendUint8", uint8(math.MaxUint8), []byte("255")},
		{"AppendUint16(math.MaxUint16)", "AppendUint16", uint16(math.MaxUint16), []byte("65535")},
		{"AppendUint32(math.MaxUint32)", "AppendUint32", uint32(math.MaxUint32), []byte("4294967295")},
		{"AppendUint64(math.MaxUint64)", "AppendUint64", uint64(math.MaxUint64), []byte("18446744073709551615")},

		{"AppendFloat32(-Inf)", "AppendFloat32", float32(math.Inf(-1)), []byte(`"-Inf"`)},
		{"AppendFloat32(+Inf)", "AppendFloat32", float32(math.Inf(1)), []byte(`"+Inf"`)},
		{"AppendFloat32(NaN)", "AppendFloat32", float32(math.NaN()), []byte(`"NaN"`)},
		{"AppendFloat32(0)", "AppendFloat32", float32(0), []byte(`0`)},
		{"AppendFloat32(-1.1)", "AppendFloat32", float32(-1.1), []byte(`-1.1`)},
		{"AppendFloat32(1e20)", "AppendFloat32", float32(1e20), []byte(`100000000000000000000`)},
		{"AppendFloat32(1e21)", "AppendFloat32", float32(1e21), []byte(`1000000000000000000000`)},

		{"AppendFloat64(-Inf)", "AppendFloat64", float64(math.Inf(-1)), []byte(`"-Inf"`)},
		{"AppendFloat64(+Inf)", "AppendFloat64", float64(math.Inf(1)), []byte(`"+Inf"`)},
		{"AppendFloat64(NaN)", "AppendFloat64", float64(math.NaN()), []byte(`"NaN"`)},
		{"AppendFloat64(0)", "AppendFloat64", float64(0), []byte(`0`)},
		{"AppendFloat64(-1.1)", "AppendFloat64", float64(-1.1), []byte(`-1.1`)},
		{"AppendFloat64(1e20)", "AppendFloat64", float64(1e20), []byte(`100000000000000000000`)},
		{"AppendFloat64(1e21)", "AppendFloat64", float64(1e21), []byte(`1000000000000000000000`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := w[tt.fn](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_appendMAC(t *testing.T) {
	MACtests := []struct {
		input string
		want  []byte
	}{
		{"01:23:45:67:89:ab", []byte(`"01:23:45:67:89:ab"`)},
		{"cd:ef:11:22:33:44", []byte(`"cd:ef:11:22:33:44"`)},
	}
	for _, tt := range MACtests {
		t.Run("MAC", func(t *testing.T) {
			ha, _ := net.ParseMAC(tt.input)
			if got := enc.AppendMACAddr([]byte{}, ha); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendMACAddr() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_appendIP(t *testing.T) {
	IPv4tests := []struct {
		input net.IP
		want  []byte
	}{
		{net.IP{0, 0, 0, 0}, []byte(`"0.0.0.0"`)},
		{net.IP{192, 0, 2, 200}, []byte(`"192.0.2.200"`)},
	}

	for _, tt := range IPv4tests {
		t.Run("IPv4", func(t *testing.T) {
			if got := enc.AppendIPAddr([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPAddr() = %s, want %s", got, tt.want)
			}
		})
	}
	IPv6tests := []struct {
		input net.IP
		want  []byte
	}{
		{net.IPv6zero, []byte(`"::"`)},
		{net.IPv6linklocalallnodes, []byte(`"ff02::1"`)},
		{net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}, []byte(`"2001:db8:85a3::8a2e:370:7334"`)},
	}
	for _, tt := range IPv6tests {
		t.Run("IPv6", func(t *testing.T) {
			if got := enc.AppendIPAddr([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPAddr() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_appendIPPrefix(t *testing.T) {
	IPv4Prefixtests := []struct {
		input net.IPNet
		want  []byte
	}{
		{net.IPNet{IP: net.IP{0, 0, 0, 0}, Mask: net.IPv4Mask(0, 0, 0, 0)}, []byte(`"0.0.0.0/0"`)},
		{net.IPNet{IP: net.IP{192, 0, 2, 200}, Mask: net.IPv4Mask(255, 255, 255, 0)}, []byte(`"192.0.2.200/24"`)},
	}
	for _, tt := range IPv4Prefixtests {
		t.Run("IPv4", func(t *testing.T) {
			if got := enc.AppendIPPrefix([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPPrefix() = %s, want %s", got, tt.want)
			}
		})
	}
	IPv6Prefixtests := []struct {
		input net.IPNet
		want  []byte
	}{
		{net.IPNet{IP: net.IPv6zero, Mask: net.CIDRMask(0, 128)}, []byte(`"::/0"`)},
		{net.IPNet{IP: net.IPv6linklocalallnodes, Mask: net.CIDRMask(128, 128)}, []byte(`"ff02::1/128"`)},
		{net.IPNet{IP: net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
			Mask: net.CIDRMask(64, 128)},
			[]byte(`"2001:db8:85a3::8a2e:370:7334/64"`)},
	}
	for _, tt := range IPv6Prefixtests {
		t.Run("IPv6", func(t *testing.T) {
			if got := enc.AppendIPPrefix([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPPrefix() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_appendMac(t *testing.T) {
	MACtests := []struct {
		input net.HardwareAddr
		want  []byte
	}{
		{net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}, []byte(`"12:34:56:78:90:ab"`)},
		{net.HardwareAddr{0x12, 0x34, 0x00, 0x00, 0x90, 0xab}, []byte(`"12:34:00:00:90:ab"`)},
	}

	for _, tt := range MACtests {
		t.Run("MAC", func(t *testing.T) {
			if got := enc.AppendMACAddr([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendMAC() = %s, want %s", got, tt.want)
			}
		})
	}
}
