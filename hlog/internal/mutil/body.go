package mutil

import "io"

type byteCountReadCloser struct {
	rc   io.ReadCloser
	read int64
}

var _ io.ReadCloser = (*byteCountReadCloser)(nil)

func NewByteCountReadCloser(body io.ReadCloser) *byteCountReadCloser {
	return &byteCountReadCloser{
		rc: body,
	}
}

func (bcrc *byteCountReadCloser) Read(p []byte) (int, error) {
	n, err := bcrc.rc.Read(p)
	bcrc.read += int64(n)
	return n, err
}

func (bcrc *byteCountReadCloser) Close() error {
	return bcrc.rc.Close()
}

func (bcrc *byteCountReadCloser) BytesRead() int64 {
	return bcrc.read
}
