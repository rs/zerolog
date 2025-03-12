package diode_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/internal/cbor"
)

func TestNewWriter(t *testing.T) {
	buf := bytes.Buffer{}
	w := diode.NewWriter(&buf, 1000, 0, func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	})
	log := zerolog.New(w)
	log.Print("test")

	w.Close()
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
	if os.Getenv("TEST_FATAL") == "1" {
		w := diode.NewWriter(os.Stderr, 1000, 0, func(missed int) {
			fmt.Printf("Dropped %d messages\n", missed)
		})
		defer w.Close()
		log := zerolog.New(w)
		log.Fatal().Msg("test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	slurp, err := io.ReadAll(stderr)
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err == nil {
		t.Error("Expected log.Fatal to exit with non-zero status")
	}

	want := "{\"level\":\"fatal\",\"message\":\"test\"}\n"
	got := cbor.DecodeIfBinaryToString(slurp)
	if got != want {
		t.Errorf("Diode Fatal Test failed. got:%s, want:%s!", got, want)
	}
}

type SlowWriter struct{}

func (rw *SlowWriter) Write(p []byte) (n int, err error) {
	time.Sleep(200 * time.Millisecond)
	fmt.Print(string(p))
	return len(p), nil
}

func TestFatalWithFilteredLevelWriter(t *testing.T) {
	if os.Getenv("TEST_FATAL_SLOW") == "1" {
		slowWriter := SlowWriter{}
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
		logger.Fatal().Msg("test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalWithFilteredLevelWriter")
	cmd.Env = append(os.Environ(), "TEST_FATAL_SLOW=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	slurp, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err == nil {
		t.Error("Expected log.Fatal to exit with non-zero status")
	}

	got := cbor.DecodeIfBinaryToString(slurp)
	want := "{\"level\":\"fatal\",\"message\":\"test\"}\n"
	if got != want {
		t.Errorf("Expected output %q, got: %q", want, got)
	}
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
