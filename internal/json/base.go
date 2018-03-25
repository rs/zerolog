package json

// AppendKey appends a new key to the output JSON.
func AppendKey(dst []byte, key string) []byte {
	if len(dst) > 1 && dst[len(dst)-1] != '{' {
		dst = append(dst, ',')
	}
	dst = AppendString(dst, key)
	return append(dst, ':')
}

// AppendError encodes the error string to json and appends
// the encoded string to the input byte slice.
func AppendError(dst []byte, err error) []byte {
	if err == nil {
		return append(dst, `null`...)
	}
	return AppendString(dst, err.Error())
}

// AppendErrors encodes the error strings to json and
// appends the encoded string list to the input byte slice.
func AppendErrors(dst []byte, errs []error) []byte {
	if len(errs) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	if errs[0] != nil {
		dst = AppendString(dst, errs[0].Error())
	} else {
		dst = append(dst, "null"...)
	}
	if len(errs) > 1 {
		for _, err := range errs[1:] {
			if err == nil {
				dst = append(dst, ",null"...)
				continue
			}
			dst = AppendString(append(dst, ','), err.Error())
		}
	}
	dst = append(dst, ']')
	return dst
}
