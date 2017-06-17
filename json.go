package zerolog

import "unicode/utf8"

const hex = "0123456789abcdef"

func appendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' || s[i] == '\\' || s[i] == '"' || s[i] > 126 {
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
							dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
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
			return append(dst, '"')
		}
	}
	dst = append(dst, s...)
	return append(dst, '"')
}
