//go:build go1.18
// +build go1.18

package cbor

import "net/netip"

func ipString(octets []byte) string {
	addr, ok := netip.AddrFromSlice(octets)
	if !ok {
		return ""
	}
	return addr.String()
}

func ipPfxString(octets []byte, bits int) string {
	addr, ok := netip.AddrFromSlice(octets)
	if !ok {
		return ""
	}
	return netip.PrefixFrom(addr, bits).String()
}
