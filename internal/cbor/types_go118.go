//go:build go1.18
// +build go1.18

package cbor

import "net/netip"

// AppendNetipAddr encodes and inserts an IP Address (IPv4 or IPv6).
func (e Encoder) AppendNetipAddr(dst []byte, addr netip.Addr) []byte {
	dst = append(dst, majorTypeTags|additionalTypeIntUint16)
	dst = append(dst, byte(additionalTypeTagNetworkAddr>>8))
	dst = append(dst, byte(additionalTypeTagNetworkAddr&0xff))
	return e.AppendBytes(dst, addr.AsSlice())
}

// AppendNetipPrefix encodes and inserts an IP Address Prefix (Address + Mask Length).
func (e Encoder) AppendNetipPrefix(dst []byte, pfx netip.Prefix) []byte {
	dst = append(dst, majorTypeTags|additionalTypeIntUint16)
	dst = append(dst, byte(additionalTypeTagNetworkPrefix>>8))
	dst = append(dst, byte(additionalTypeTagNetworkPrefix&0xff))

	// Prefix is a tuple (aka MAP of 1 pair of elements) -
	// first element is prefix, second is mask length.
	dst = append(dst, majorTypeMap|0x1)
	dst = e.AppendBytes(dst, pfx.Addr().AsSlice())
	return e.AppendInt(dst, pfx.Bits())
}
