package diode_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	diodes "code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func ExampleNewWriter() {
	d := diodes.NewManyToOne(1000, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	w := diode.NewWriter(os.Stdout, d, 10*time.Millisecond)
	log := zerolog.New(w)
	log.Print("test")

	w.Close()

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
