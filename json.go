package zerolog

import "unicode/utf8"

const hex = "0123456789abcdef"

// appendJSONString encodes the input string to json and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte in the string looking
// for characters that need json or utf8 encoding. If the string
// does not need encoding, then the string is appended in it's
// entirety to the byte slice.
// If we encounter a byte that does need encoding, switch up
// the operation and perform a byte-by-byte read-encode-append.
func appendJSONString(dst []byte, s string) []byte {
	// Start with a double quote.
	dst = append(dst, '"')
	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if s[i] < ' ' || s[i] == '\\' || s[i] == '"' || s[i] > 126 {
			// We encountered a character that needs to be encoded. Let's
			// append the previous simple characters to the byte slice
			// and switch our operation to read and encode the remainder
			// characters byte-by-byte.
			dst = append(dst, s[:i]...)
			for i < len(s) {
				if b := s[i]; b < utf8.RuneSelf {
					switch b {
					case '"', '\\':
						dst = append(dst, '\\', b)
					case '\b':
						dst = append(dst, '\\', 'b')
					case '\f':
						dst = append(dst, '\\', 'f')
					case '\n':
						dst = append(dst, '\\', 'n')
					case '\r':
						dst = append(dst, '\\', 'r')
					case '\t':
						dst = append(dst, '\\', 't')
					default:
						if b >= 0x20 {
							dst = append(dst, b)
						} else {
							dst = append(dst, '\\', 'u', '0', '0',
								hex[b>>4], hex[b&0xF])
						}
					}
					i++
					continue
				}
				r, size := utf8.DecodeRuneInString(s[i:])
				if r == utf8.RuneError && size == 1 {
					dst = append(dst, `\ufffd`...)
					i++
					continue
				}
				dst = append(dst, s[i:i+size]...)
				i += size
			}
			// End with a double quote
			return append(dst, '"')
		}
	}
	// The string has no need for encoding an therefore is directly
	// appended to the byte slice.
	dst = append(dst, s...)
	// End with a double quote
	return append(dst, '"')
}
