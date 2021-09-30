package zerolog

import (
	"bytes"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// LevelWriter defines as interface a writer may implement in order
// to receive level information with payload.
type LevelWriter interface {
	io.Writer
	WriteLevel(level Level, p []byte) (n int, err error)
}

type levelWriterAdapter struct {
	io.Writer
}

func (lw levelWriterAdapter) WriteLevel(l Level, p []byte) (n int, err error) {
	return lw.Write(p)
}

type syncWriter struct {
	mu sync.Mutex
	lw LevelWriter
}

// SyncWriter wraps w so that each call to Write is synchronized with a mutex.
// This syncer can be used to wrap the call to writer's Write method if it is
// not thread safe. Note that you do not need this wrapper for os.File Write
// operations on POSIX and Windows systems as they are already thread-safe.
func SyncWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LevelWriter); ok {
		return &syncWriter{lw: lw}
	}
	return &syncWriter{lw: levelWriterAdapter{w}}
}

// Write implements the io.Writer interface.
func (s *syncWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.Write(p)
}

// WriteLevel implements the LevelWriter interface.
func (s *syncWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.WriteLevel(l, p)
}

type multiLevelWriter struct {
	writers []LevelWriter
}

func (t multiLevelWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		if _n, _err := w.Write(p); err == nil {
			n = _n
			if _err != nil {
				err = _err
			} else if _n != len(p) {
				err = io.ErrShortWrite
			}
		}
	}
	return n, err
}

func (t multiLevelWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	for _, w := range t.writers {
		if _n, _err := w.WriteLevel(l, p); err == nil {
			n = _n
			if _err != nil {
				err = _err
			} else if _n != len(p) {
				err = io.ErrShortWrite
			}
		}
	}
	return n, err
}

// MultiLevelWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command. If some writers
// implement LevelWriter, their WriteLevel method will be used instead of Write.
func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	lwriters := make([]LevelWriter, 0, len(writers))
	for _, w := range writers {
		if lw, ok := w.(LevelWriter); ok {
			lwriters = append(lwriters, lw)
		} else {
			lwriters = append(lwriters, levelWriterAdapter{w})
		}
	}
	return multiLevelWriter{lwriters}
}

// TestWriter is a writer that writes to testing.TB.
type TestWriter struct {
	TB testing.TB

	// Frame skips caller frames to capture the original file and line numbers.
	Frame int
}

// NewTestWriter creates a writer that logs to the testing.TB.
func NewTestWriter(t testing.TB) TestWriter {
	return TestWriter{TB: t, Frame: 1}
}

// Write to testing.TB.
func (t TestWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	// Strip trailing newline because t.Log always adds one.
	p = bytes.TrimRight(p, "\n")

	// Try to correct the log file and line number to the caller.
	if t.Frame > 0 {
		_, origFile, origLine, _ := runtime.Caller(0)
		_, frameFile, frameLine, ok := runtime.Caller(t.Frame)
		if ok {
			erase := strings.Repeat("\b", len(path.Base(origFile))+len(strconv.Itoa(origLine))+3)
			t.TB.Logf("%s%s:%d: %s", erase, path.Base(frameFile), frameLine, p)
			return n, err
		}
	}
	t.TB.Logf("%s", p)

	return n, err
}

// ConsoleTestWriter creates an option that correctly sets the file frame depth for testing.TB log.
func ConsoleTestWriter(t testing.TB) func(w *ConsoleWriter) {
	return func(w *ConsoleWriter) {
		w.Out = TestWriter{TB: t, Frame: 7}
	}
}
