//go:build go1.18
// +build go1.18

package json

import (
	"net/netip"
	"reflect"
	"testing"
)

func Test_appendNetipAddr(t *testing.T) {
	IPv4tests := []struct {
		input netip.Addr
		want  []byte
	}{
		{netip.IPv4Unspecified(), []byte(`"0.0.0.0"`)},
		{netip.AddrFrom4([4]byte{192, 0, 2, 200}), []byte(`"192.0.2.200"`)},
	}

	for _, tt := range IPv4tests {
		t.Run("IPv4", func(t *testing.T) {
			if got := enc.AppendNetipAddr([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPAddr() = %s, want %s", got, tt.want)
			}
		})
	}
	IPv6tests := []struct {
		input netip.Addr
		want  []byte
	}{
		{netip.IPv6Unspecified(), []byte(`"::"`)},
		{netip.IPv6LinkLocalAllNodes(), []byte(`"ff02::1"`)},
		{netip.AddrFrom16([16]byte{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}), []byte(`"2001:db8:85a3::8a2e:370:7334"`)},
	}
	for _, tt := range IPv6tests {
		t.Run("IPv6", func(t *testing.T) {
			if got := enc.AppendNetipAddr([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPAddr() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_appendNetipPrefix(t *testing.T) {
	IPv4Prefixtests := []struct {
		input netip.Prefix
		want  []byte
	}{
		{netip.PrefixFrom(netip.IPv4Unspecified(), 0), []byte(`"0.0.0.0/0"`)},
		{netip.PrefixFrom(netip.AddrFrom4([4]byte{192, 0, 2, 200}), 24), []byte(`"192.0.2.200/24"`)},
	}
	for _, tt := range IPv4Prefixtests {
		t.Run("IPv4", func(t *testing.T) {
			if got := enc.AppendNetipPrefix([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPPrefix() = %s, want %s", got, tt.want)
			}
		})
	}
	IPv6Prefixtests := []struct {
		input netip.Prefix
		want  []byte
	}{
		{netip.PrefixFrom(netip.IPv6Unspecified(), 0), []byte(`"::/0"`)},
		{netip.PrefixFrom(netip.IPv6LinkLocalAllNodes(), 128), []byte(`"ff02::1/128"`)},
	}
	for _, tt := range IPv6Prefixtests {
		t.Run("IPv6", func(t *testing.T) {
			if got := enc.AppendNetipPrefix([]byte{}, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendIPPrefix() = %s, want %s", got, tt.want)
			}
		})
	}
}
