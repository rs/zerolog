package cbor

// AppendKey adds a key (string) to the binary encoded log message
func AppendKey(dst []byte, key string) []byte {
	if len(dst) < 1 {
		dst = AppendBeginMarker(dst)
	}
	return AppendString(dst, key)
}

// AppendError adds the Error to the log message if error is NOT nil
func AppendError(dst []byte, err error) []byte {
	if err == nil {
		return append(dst, `null`...)
	}
	return AppendString(dst, err.Error())
}

// AppendErrors when given an array of errors,
// adds them to the log message if a specific error is nil, then
// Nil is added, or else the error string is added.
func AppendErrors(dst []byte, errs []error) []byte {
	if len(errs) == 0 {
		return AppendArrayEnd(AppendArrayStart(dst))
	}
	dst = AppendArrayStart(dst)
	if errs[0] != nil {
		dst = AppendString(dst, errs[0].Error())
	} else {
		dst = AppendNull(dst)
	}
	if len(errs) > 1 {
		for _, err := range errs[1:] {
			if err == nil {
				dst = AppendNull(dst)
				continue
			}
			dst = AppendString(dst, err.Error())
		}
	}
	dst = AppendArrayEnd(dst)
	return dst
}
