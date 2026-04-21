package diode_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/internal/cbor"
)

type signalWriter struct {
	w     io.Writer
	wrote chan struct{}
}

func (w *signalWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	if err == nil {
		select {
		case w.wrote <- struct{}{}:
		default:
		}
	}
	return n, err
}

func TestNewWriter(t *testing.T) {
	var buf bytes.Buffer
	sw := signalWriter{w: &buf, wrote: make(chan struct{}, 1)}
	w := diode.NewWriter(&sw, 1000, 0, func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	})
	log := zerolog.New(w)
	log.Print("test")

	select {
	case <-sw.wrote:
	case <-time.After(2 * time.Second):
		w.Close()
		t.Fatal("timed out waiting for diode writer to flush")
	}

	_ = w.Close()
	want := "{\"level\":\"debug\",\"message\":\"test\"}\n"
	got := cbor.DecodeIfBinaryToString(buf.Bytes())
	if got != want {
		t.Errorf("Diode New Writer Test failed. got:%s, want:%s!", got, want)
	}
}

func TestClose(t *testing.T) {
	buf := bytes.Buffer{}
	w := diode.NewWriter(&buf, 1000, 0, func(missed int) {})
	log := zerolog.New(w)
	log.Print("test")
	w.Close()
}

func TestFatal(t *testing.T) {
	var buf bytes.Buffer
	w := diode.NewWriter(&buf, 1000, 0, func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	})
	log := zerolog.New(w)

	oldExitFunc := zerolog.FatalExitFunc
	zerolog.FatalExitFunc = func() { panic("fatal exit") }
	defer func() { zerolog.FatalExitFunc = oldExitFunc }()

	defer func() {
		if r := recover(); r == nil || r != "fatal exit" {
			t.Fatalf("expected panic %q from log.Fatal(), got %v", "fatal exit", r)
		}
		want := "{\"level\":\"fatal\",\"message\":\"test\"}\n"
		got := cbor.DecodeIfBinaryToString(buf.Bytes())
		if got != want {
			t.Errorf("Diode Fatal Test failed. got:%s, want:%s!", got, want)
		}
	}()

	log.Fatal().Msg("test")
}

type SlowWriter struct{ w io.Writer }

func (rw *SlowWriter) Write(p []byte) (n int, err error) {
	time.Sleep(200 * time.Millisecond)
	return rw.w.Write(p)
}

func TestFatalWithFilteredLevelWriter(t *testing.T) {
	var buf bytes.Buffer
	slowWriter := SlowWriter{w: &buf}
	diodeWriter := diode.NewWriter(&slowWriter, 500, 0, func(missed int) {
		fmt.Printf("Missed %d logs\n", missed)
	})
	leveledDiodeWriter := zerolog.LevelWriterAdapter{
		Writer: &diodeWriter,
	}
	filteredDiodeWriter := zerolog.FilteredLevelWriter{
		Writer: &leveledDiodeWriter,
		Level:  zerolog.InfoLevel,
	}
	logger := zerolog.New(&filteredDiodeWriter)

	oldExitFunc := zerolog.FatalExitFunc
	zerolog.FatalExitFunc = func() { panic("fatal exit") }
	defer func() { zerolog.FatalExitFunc = oldExitFunc }()

	defer func() {
		if r := recover(); r == nil || r != "fatal exit" {
			t.Fatalf("expected panic %q from log.Fatal(), got %v", "fatal exit", r)
		}
		want := "{\"level\":\"fatal\",\"message\":\"test\"}\n"
		got := cbor.DecodeIfBinaryToString(buf.Bytes())
		if got != want {
			t.Errorf("Expected output %q, got: %q", want, got)
		}
	}()

	logger.Fatal().Msg("test")
}

func Benchmark(b *testing.B) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	benchs := map[string]time.Duration{
		"Waiter": 0,
		"Pooler": 10 * time.Millisecond,
	}
	for name, interval := range benchs {
		b.Run(name, func(b *testing.B) {
			w := diode.NewWriter(io.Discard, 100000, interval, nil)
			log := zerolog.New(w)
			defer w.Close()

			b.SetParallelism(1000)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					log.Print("test")
				}
			})
		})
	}
}
