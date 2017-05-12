package zerolog

import (
	"bytes"
	"unicode/utf8"
)

const hex = "0123456789abcdef"

func writeJSONString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			switch b {
			case '"', '\\':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				if b >= 0x20 {
					buf.WriteByte(b)
				} else {
					buf.WriteString(`\u00`)
					buf.WriteByte(hex[b>>4])
					buf.WriteByte(hex[b&0xF])
				}
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			buf.WriteString(`\ufffd`)
			i++
			continue
		}
		buf.WriteString(s[i : i+size])
		i += size
	}
	buf.WriteByte('"')
}
