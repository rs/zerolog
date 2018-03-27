package cbor

// AppendStrings encodes and adds an array of strings to the dst byte array.
func AppendStrings(dst []byte, vals []string) []byte {
	major := majorTypeArray
	l := len(vals)
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	for _, v := range vals {
		dst = AppendString(dst, v)
	}
	return dst
}

// AppendString encodes and adds a string to the dst byte array.
func AppendString(dst []byte, s string) []byte {
	major := majorTypeUtf8String

	l := len(s)
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, majorTypeUtf8String, uint64(l))
	}
	return append(dst, s...)
}

// AppendBytes encodes and adds an array of bytes to the dst byte array.
func AppendBytes(dst, s []byte) []byte {
	major := majorTypeByteString

	l := len(s)
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	return append(dst, s...)
}

// AppendEmbeddedJSON adds a tag and embeds input JSON as such.
func AppendEmbeddedJSON(dst, s []byte) []byte {
	major := majorTypeTags
	minor := additionalTypeEmbeddedJSON
	dst = append(dst, byte(major|minor))

	major = majorTypeByteString

	l := len(s)
	if l <= additionalMax {
		lb := byte(l)
		dst = append(dst, byte(major|lb))
	} else {
		dst = appendCborTypePrefix(dst, major, uint64(l))
	}
	return append(dst, s...)
}
