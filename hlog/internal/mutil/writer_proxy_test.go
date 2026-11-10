package mutil

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"net/textproto"
	"reflect"
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

func TestWriterInformationalResponses(t *testing.T) {
	for _, finalStatus := range []int{0, http.StatusCreated} {
		name := http.StatusText(finalStatus)
		if finalStatus == 0 {
			name = "implicit OK"
		}
		t.Run(name, func(t *testing.T) {
			statuses := make(chan int, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writer := WrapWriter(w)
				writer.WriteHeader(http.StatusEarlyHints)
				writer.WriteHeader(http.StatusProcessing)
				if finalStatus != 0 {
					writer.WriteHeader(finalStatus)
				}
				if _, err := writer.Write([]byte("response body")); err != nil {
					t.Error(err)
				}
				statuses <- writer.Status()
			}))
			defer server.Close()

			var informational []int
			trace := &httptrace.ClientTrace{
				Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
					informational = append(informational, code)
					return nil
				},
			}
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			wantStatus := finalStatus
			if wantStatus == 0 {
				wantStatus = http.StatusOK
			}
			if resp.StatusCode != wantStatus {
				t.Errorf("response status = %d, want %d", resp.StatusCode, wantStatus)
			}
			if got := <-statuses; got != wantStatus {
				t.Errorf("writer status = %d, want %d", got, wantStatus)
			}
			if want := []int{http.StatusEarlyHints, http.StatusProcessing}; !reflect.DeepEqual(informational, want) {
				t.Errorf("informational statuses = %v, want %v", informational, want)
			}
			if string(body) != "response body" {
				t.Errorf("response body = %q", body)
			}
		})
	}
}

func TestWriterSwitchingProtocolsIsFinal(t *testing.T) {
	response := httptest.NewRecorder()
	writer := WrapWriter(response)
	writer.WriteHeader(http.StatusSwitchingProtocols)
	writer.WriteHeader(http.StatusOK)
	if got := writer.Status(); got != http.StatusSwitchingProtocols {
		t.Errorf("writer status = %d, want 101", got)
	}
	if response.Code != http.StatusSwitchingProtocols {
		t.Errorf("response status = %d, want 101", response.Code)
	}
}
