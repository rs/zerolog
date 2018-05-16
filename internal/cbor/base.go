package cbor

type Encoder struct{}

// AppendKey adds a key (string) to the binary encoded log message
func (e Encoder) AppendKey(dst []byte, key string) []byte {
	if len(dst) < 1 {
		dst = e.AppendBeginMarker(dst)
	}
	return e.AppendString(dst, key)
}

// AppendError adds the Error to the log message if error is NOT nil
func (e Encoder) AppendError(dst []byte, err error) []byte {
	if err == nil {
		return append(dst, `null`...)
	}
	return e.AppendString(dst, err.Error())
}

// AppendErrors when given an array of errors,
// adds them to the log message if a specific error is nil, then
// Nil is added, or else the error string is added.
func (e Encoder) AppendErrors(dst []byte, errs []error) []byte {
	if len(errs) == 0 {
		return e.AppendArrayEnd(e.AppendArrayStart(dst))
	}
	dst = e.AppendArrayStart(dst)
	if errs[0] != nil {
		dst = e.AppendString(dst, errs[0].Error())
	} else {
		dst = e.AppendNil(dst)
	}
	if len(errs) > 1 {
		for _, err := range errs[1:] {
			if err == nil {
				dst = e.AppendNil(dst)
				continue
			}
			dst = e.AppendString(dst, err.Error())
		}
	}
	dst = e.AppendArrayEnd(dst)
	return dst
}
