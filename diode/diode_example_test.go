// +build !binary_log

package diode_test

import (
	"fmt"
	"os"
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
