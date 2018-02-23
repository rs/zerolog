package diode_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
	"bytes"

	diodes "code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/internal/cbor"
)

func ExampleNewWriter() {
	d := diodes.NewManyToOne(1000, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	buf := bytes.Buffer{}
	w := diode.NewWriter(&buf, d, 10*time.Millisecond)
	log := zerolog.New(w)
	log.Print("test")

	w.Close()
	fmt.Println(cbor.DecodeIfBinaryToString(buf.Bytes()))

	// Output: {"level":"debug","message":"test"}
}

func Benchmark(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	d := diodes.NewManyToOne(100000, nil)
	w := diode.NewWriter(ioutil.Discard, d, 10*time.Millisecond)
	log := zerolog.New(w)
	defer w.Close()

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("test")
		}
	})

}
