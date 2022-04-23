//go:build !go1.18
// +build !go1.18

package cbor

import "net"

func ipString(octets []byte) string {
	return net.IP(octets).String()
}

func ipPfxString(octets []byte, bits int) string {
	ip := net.IP(octets)
	var mask net.IPMask
	if len(octets) == 4 {
		mask = net.CIDRMask(bits, 32)
	} else {
		mask = net.CIDRMask(bits, 128)
	}
	ipPfx := net.IPNet{IP: ip, Mask: mask}
	return ipPfx.String()
}
