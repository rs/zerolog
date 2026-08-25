package mutil

import (
	"io"
	"sync/atomic"
)

type byteCountReadCloser struct {
	rc   io.ReadCloser
	read *int64
}

var _ io.ReadCloser = (*byteCountReadCloser)(nil)
var _ io.WriterTo = (*byteCountReadCloser)(nil)

func NewByteCountReadCloser(rc io.ReadCloser) *byteCountReadCloser {
	read := int64(0)
	return &byteCountReadCloser{
		rc:   rc,
		read: &read,
	}
}

func (b *byteCountReadCloser) Read(p []byte) (int, error) {
	n, err := b.rc.Read(p)
	atomic.AddInt64(b.read, int64(n))
	return n, err
}

func (b *byteCountReadCloser) Close() error {
	return b.rc.Close()
}

func (b *byteCountReadCloser) WriteTo(w io.Writer) (int64, error) {
	n, err := io.Copy(w, b.rc)
	atomic.AddInt64(b.read, n)
	return n, err
}

func (b *byteCountReadCloser) BytesRead() int64 {
	return atomic.LoadInt64(b.read)
}
