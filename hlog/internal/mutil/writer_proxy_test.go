package mutil

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

type readerFromResponseWriter struct {
	bytes.Buffer
	header http.Header
	status int
}

func (w *readerFromResponseWriter) Header() http.Header {
	return w.header
}

func (w *readerFromResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *readerFromResponseWriter) CloseNotify() <-chan bool {
	return make(chan bool)
}

func (w *readerFromResponseWriter) Flush() {}

func (w *readerFromResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestFancyWriterReadFromWithTeeCountsBytesOnce(t *testing.T) {
	response := &readerFromResponseWriter{header: make(http.Header)}
	writer := WrapWriter(response)
	var tee bytes.Buffer
	writer.Tee(&tee)

	const body = "response body"
	written, err := io.Copy(writer, io.LimitReader(strings.NewReader(body), int64(len(body))))
	if err != nil {
		t.Fatal(err)
	}
	if written != int64(len(body)) {
		t.Fatalf("io.Copy wrote %d bytes, want %d", written, len(body))
	}
	if got := writer.BytesWritten(); got != len(body) {
		t.Errorf("BytesWritten() = %d, want %d", got, len(body))
	}
	if got := response.String(); got != body {
		t.Errorf("response body = %q, want %q", got, body)
	}
	if got := tee.String(); got != body {
		t.Errorf("tee body = %q, want %q", got, body)
	}
}
