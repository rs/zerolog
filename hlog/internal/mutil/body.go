package mutil

import (
	"io"
	"sync/atomic"
)

type byteCountReadCloser struct {
	rc   io.ReadCloser
	read int64
}

var _ io.ReadCloser = (*byteCountReadCloser)(nil)

func NewByteCountReadCloser(rc io.ReadCloser) *byteCountReadCloser {
	return &byteCountReadCloser{
		rc: rc,
	}
}

func (b *byteCountReadCloser) Read(p []byte) (int, error) {
	n, err := b.rc.Read(p)
	atomic.AddInt64(&b.read, int64(n))
	return n, err
}

func (b *byteCountReadCloser) Close() error {
	return b.rc.Close()
}

func (b *byteCountReadCloser) BytesRead() int64 {
	return b.read
}
